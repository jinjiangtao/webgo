package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserList(c *gin.Context) {
	username := c.Query("username")
	nickname := c.Query("nickname")
	deptId := c.Query("deptId")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	db := utils.DB.Model(&models.User{})
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	if deptId != "" {
		db = db.Where("dept_id = ?", deptId)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	if err := db.Preload("Roles").Preload("Department").
		Offset(offset).Limit(pageSize).
		Order("id DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": users, "total": total, "page": page, "pageSize": pageSize})
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user models.User
	if err := utils.DB.Preload("Roles").Preload("Department").First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if user.Password == "" {
		user.Password = "123456"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}
	user.Password = string(hashedPassword)

	tx := utils.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建用户失败"})
		return
	}

	if len(user.Roles) > 0 {
		var roleIds []uint
		for _, r := range user.Roles {
			roleIds = append(roleIds, r.ID)
		}
		for _, rid := range roleIds {
			tx.Create(&models.UserRole{UserID: user.ID, RoleID: rid})
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": user})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user models.User
	if err := utils.DB.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	var input struct {
		Username *string       `json:"username"`
		Nickname *string       `json:"nickname"`
		Email    *string       `json:"email"`
		Phone    *string       `json:"phone"`
		Avatar   *string       `json:"avatar"`
		DeptID   *uint         `json:"deptId"`
		Status   *int          `json:"status"`
		Roles    []models.Role `json:"roles"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if input.Username != nil {
		updates["username"] = *input.Username
	}
	if input.Nickname != nil {
		updates["nickname"] = *input.Nickname
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.Avatar != nil {
		updates["avatar"] = *input.Avatar
	}
	if input.DeptID != nil {
		updates["dept_id"] = *input.DeptID
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	tx := utils.DB.Begin()
	if len(updates) > 0 {
		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新用户失败"})
			return
		}
	}

	if input.Roles != nil {
		tx.Where("user_id = ?", id).Delete(&models.UserRole{})
		for _, r := range input.Roles {
			if r.ID > 0 {
				tx.Create(&models.UserRole{UserID: user.ID, RoleID: r.ID})
			}
		}
	}

	tx.Commit()
	utils.DB.First(&user, uint(id))

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": user})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var user models.User
	if err := utils.DB.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	if user.Username == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "超级管理员不可删除"})
		return
	}

	tx := utils.DB.Begin()
	tx.Where("user_id = ?", id).Delete(&models.UserRole{})
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除用户失败"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	if err := utils.DB.Model(&models.User{}).Where("id = ?", id).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "重置密码成功"})
}
