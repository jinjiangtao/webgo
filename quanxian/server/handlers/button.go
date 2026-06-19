package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetButtonList(c *gin.Context) {
	menuId := c.Query("menuId")
	name := c.Query("name")

	db := utils.DB.Model(&models.Button{})
	if menuId != "" {
		db = db.Where("menu_id = ?", menuId)
	}
	if name != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	var buttons []models.Button
	if err := db.Order("sort ASC, id ASC").Find(&buttons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询按钮失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": buttons, "total": len(buttons)})
}

func GetButton(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var button models.Button
	if err := utils.DB.First(&button, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "按钮不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": button})
}

func CreateButton(c *gin.Context) {
	var button models.Button
	if err := c.ShouldBindJSON(&button); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := utils.DB.Create(&button).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建按钮失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": button})
}

func UpdateButton(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var button models.Button
	if err := utils.DB.First(&button, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "按钮不存在"})
		return
	}

	if err := c.ShouldBindJSON(&button); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	button.ID = uint(id)
	if err := utils.DB.Save(&button).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新按钮失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": button})
}

func DeleteButton(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	utils.DB.Where("button_id = ?", id).Delete(&models.RoleButton{})

	if err := utils.DB.Delete(&models.Button{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除按钮失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
