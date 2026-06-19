package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDepartmentTree(c *gin.Context) {
	var depts []models.Department
	if err := utils.DB.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询部门失败"})
		return
	}

	tree := utils.BuildDepartmentTree(depts)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree})
}

func GetDepartmentList(c *gin.Context) {
	name := c.Query("name")
	status := c.Query("status")

	db := utils.DB.Model(&models.Department{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var depts []models.Department
	if err := db.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询部门失败"})
		return
	}

	tree := utils.BuildDepartmentTree(depts)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree, "total": len(depts)})
}

func GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var dept models.Department
	if err := utils.DB.First(&dept, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": dept})
}

func CreateDepartment(c *gin.Context) {
	var dept models.Department
	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := utils.DB.Create(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建部门失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": dept})
}

func UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var dept models.Department
	if err := utils.DB.First(&dept, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}

	if err := c.ShouldBindJSON(&dept); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if dept.ParentID != nil && *dept.ParentID == dept.ID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不能将自己设为父级部门"})
		return
	}

	if err := utils.DB.Save(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新部门失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": dept})
}

func DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var childCount int64
	utils.DB.Model(&models.Department{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "存在子部门，无法删除"})
		return
	}

	var userCount int64
	utils.DB.Model(&models.User{}).Where("dept_id = ?", id).Count(&userCount)
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "部门下存在用户，无法删除"})
		return
	}

	if err := utils.DB.Delete(&models.Department{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除部门失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func MoveDepartment(c *gin.Context) {
	var req struct {
		ID          uint  `json:"id" binding:"required"`
		ParentID    *uint `json:"parentId"`
		TargetIndex int   `json:"targetIndex"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var dept models.Department
	if err := utils.DB.First(&dept, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "部门不存在"})
		return
	}

	dept.ParentID = req.ParentID
	if err := utils.DB.Save(&dept).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "移动部门失败"})
		return
	}

	var siblings []models.Department
	query := utils.DB.Where("parent_id IS NULL")
	if req.ParentID != nil {
		query = utils.DB.Where("parent_id = ?", *req.ParentID)
	}
	query.Order("sort ASC, id ASC").Find(&siblings)

	for i, s := range siblings {
		s.Sort = i + 1
		if s.ID == dept.ID {
			s.Sort = req.TargetIndex + 1
		}
		utils.DB.Save(&s)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "移动成功"})
}
