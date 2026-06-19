package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPermissionPreview(c *gin.Context) {
	roleIdStr := c.Query("roleId")
	userIdStr := c.Query("userId")

	var roleIds []uint
	var user *models.User

	if userIdStr != "" {
		userId, _ := strconv.ParseUint(userIdStr, 10, 32)
		var u models.User
		if err := utils.DB.Preload("Roles").First(&u, uint(userId)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
			return
		}
		user = &u
		for _, r := range u.Roles {
			roleIds = append(roleIds, r.ID)
		}
	} else if roleIdStr != "" {
		roleId, _ := strconv.ParseUint(roleIdStr, 10, 32)
		roleIds = append(roleIds, uint(roleId))
	}

	var roles []models.Role
	if len(roleIds) > 0 {
		utils.DB.Where("id IN ?", roleIds).Find(&roles)
	} else {
		utils.DB.Find(&roles)
	}

	var menus []models.Menu
	utils.DB.Preload("Buttons").Where("status = ?", 1).Order("sort ASC, id ASC").Find(&menus)

	var filteredMenuIds []uint
	menuIdSet := make(map[uint]bool)
	for _, rid := range roleIds {
		if utils.IsSuperAdminByRole(rid) {
			for _, m := range menus {
				menuIdSet[m.ID] = true
			}
			break
		}
		mids := utils.GetRoleMenuIds(rid)
		for _, mid := range mids {
			menuIdSet[mid] = true
		}
	}

	for mid := range menuIdSet {
		filteredMenuIds = append(filteredMenuIds, mid)
	}

	filteredMenus := utils.FilterMenusByRole(menus, filteredMenuIds)

	buttonMap := make(map[string][]string)
	buttonCodeSet := make(map[string]bool)
	for _, rid := range roleIds {
		if utils.IsSuperAdminByRole(rid) {
			for _, m := range menus {
				for _, b := range m.Buttons {
					buttonCodeSet[b.Code] = true
					if _, exists := buttonMap[m.Path]; !exists {
						buttonMap[m.Path] = []string{}
					}
				}
			}
			break
		}
		var buttons []models.Button
		utils.DB.Joins("JOIN role_buttons ON role_buttons.button_id = buttons.id").
			Where("role_buttons.role_id = ?", rid).
			Find(&buttons)
		for _, b := range buttons {
			buttonCodeSet[b.Code] = true
		}
	}

	for _, m := range menus {
		var pathBtns []string
		for _, b := range m.Buttons {
			if buttonCodeSet[b.Code] {
				pathBtns = append(pathBtns, b.Code)
			}
		}
		if len(pathBtns) > 0 {
			buttonMap[m.Path] = pathBtns
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"user":    user,
			"roles":   roles,
			"menus":   filteredMenus,
			"buttons": buttonMap,
		},
	})
}

func GetAllPermissionsTree(c *gin.Context) {
	var menus []models.Menu
	if err := utils.DB.Preload("Buttons").Where("status = ?", 1).Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}

	var nodes []models.TreeNode
	for _, menu := range menus {
		nodes = append(nodes, models.TreeNode{
			ID:       menu.ID,
			Label:    menu.Name,
			ParentID: menu.ParentID,
			Type:     "menu",
			Sort:     menu.Sort,
		})
		for _, btn := range menu.Buttons {
			menuId := menu.ID
			nodes = append(nodes, models.TreeNode{
				ID:       btn.ID + 100000,
				Label:    btn.Name,
				ParentID: &menuId,
				Type:     "button",
				Sort:     btn.Sort,
			})
		}
	}

	tree := utils.BuildGenericTree(nodes)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tree})
}

func BatchAssignPermissions(c *gin.Context) {
	var req struct {
		RoleID    uint   `json:"roleId" binding:"required"`
		MenuIds   []uint `json:"menuIds"`
		ButtonIds []uint `json:"buttonIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var role models.Role
	if err := utils.DB.First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}

	tx := utils.DB.Begin()

	tx.Where("role_id = ?", req.RoleID).Delete(&models.RoleMenu{})
	for _, mid := range req.MenuIds {
		tx.Create(&models.RoleMenu{RoleID: req.RoleID, MenuID: mid})
	}

	tx.Where("role_id = ?", req.RoleID).Delete(&models.RoleButton{})
	for _, bid := range req.ButtonIds {
		actualBid := bid
		if bid > 100000 {
			actualBid = bid - 100000
		}
		tx.Create(&models.RoleButton{RoleID: req.RoleID, ButtonID: actualBid})
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "权限配置成功"})
}
