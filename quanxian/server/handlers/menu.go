package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMenuTree(c *gin.Context) {
	var menus []models.Menu
	if err := utils.DB.Preload("Buttons").Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询菜单失败"})
		return
	}

	tree := utils.BuildMenuTree(menus)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree})
}

func GetMenuList(c *gin.Context) {
	name := c.Query("name")
	status := c.Query("status")

	db := utils.DB.Model(&models.Menu{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var menus []models.Menu
	if err := db.Preload("Buttons").Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询菜单失败"})
		return
	}

	tree := utils.BuildMenuTree(menus)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree, "total": len(menus)})
}

func GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var menu models.Menu
	if err := utils.DB.Preload("Buttons").First(&menu, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": menu})
}

func CreateMenu(c *gin.Context) {
	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if err := utils.DB.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建成功", "data": menu})
}

func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var menu models.Menu
	if err := utils.DB.First(&menu, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if menu.ParentID != nil && *menu.ParentID == menu.ID {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不能将自己设为父级菜单"})
		return
	}

	menu.ID = uint(id)
	if err := utils.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功", "data": menu})
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var childCount int64
	utils.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "存在子菜单，无法删除"})
		return
	}

	utils.DB.Where("menu_id = ?", id).Delete(&models.Button{})
	utils.DB.Where("menu_id = ?", id).Delete(&models.RoleMenu{})

	if err := utils.DB.Delete(&models.Menu{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除菜单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func MoveMenu(c *gin.Context) {
	var req struct {
		ID          uint  `json:"id" binding:"required"`
		ParentID    *uint `json:"parentId"`
		TargetIndex int   `json:"targetIndex"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var menu models.Menu
	if err := utils.DB.First(&menu, req.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "菜单不存在"})
		return
	}

	menu.ParentID = req.ParentID
	if err := utils.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "移动菜单失败"})
		return
	}

	var siblings []models.Menu
	query := utils.DB.Where("parent_id IS NULL")
	if req.ParentID != nil {
		query = utils.DB.Where("parent_id = ?", *req.ParentID)
	}
	query.Order("sort ASC, id ASC").Find(&siblings)

	for i, s := range siblings {
		s.Sort = i + 1
		if s.ID == menu.ID {
			s.Sort = req.TargetIndex + 1
		}
		utils.DB.Save(&s)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "移动成功"})
}
