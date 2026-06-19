package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRoleList(c *gin.Context) {
	name := c.Query("name")
	status := c.Query("status")

	db := utils.DB.Model(&models.Role{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var roles []models.Role
	if err := db.Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": roles, "total": len(roles)})
}

func GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := utils.DB.Preload("Menus").Preload("Buttons").First(&role, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": role})
}

func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := utils.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": role})
}

func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := utils.DB.First(&role, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	role.ID = uint(id)
	if err := utils.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": role})
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := utils.DB.First(&role, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	if role.Code == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "超级管理员角色不可删除"})
		return
	}

	var userCount int64
	utils.DB.Table("user_roles").Where("role_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "角色下存在用户，无法删除"})
		return
	}

	utils.DB.Where("role_id = ?", id).Delete(&models.RoleMenu{})
	utils.DB.Where("role_id = ?", id).Delete(&models.RoleButton{})

	if err := utils.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func BindRoleMenus(c *gin.Context) {
	roleId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var req struct {
		MenuIds []uint `json:"menuIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := utils.DB.First(&role, uint(roleId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	utils.DB.Where("role_id = ?", roleId).Delete(&models.RoleMenu{})

	for _, menuId := range req.MenuIds {
		utils.DB.Create(&models.RoleMenu{RoleID: uint(roleId), MenuID: menuId})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "菜单绑定成功"})
}

func BindRoleButtons(c *gin.Context) {
	roleId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var req struct {
		ButtonIds []uint `json:"buttonIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := utils.DB.First(&role, uint(roleId)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	utils.DB.Where("role_id = ?", roleId).Delete(&models.RoleButton{})

	for _, buttonId := range req.ButtonIds {
		utils.DB.Create(&models.RoleButton{RoleID: uint(roleId), ButtonID: buttonId})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "按钮绑定成功"})
}

func GetRoleMenus(c *gin.Context) {
	roleId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	menuIds := utils.GetRoleMenuIds(uint(roleId))
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": menuIds})
}

func GetRoleButtons(c *gin.Context) {
	roleId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	buttonIds := utils.GetRoleButtonIds(uint(roleId))
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": buttonIds})
}
