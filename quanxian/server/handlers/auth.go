package handlers

import (
	"net/http"
	"quanxian/models"
	"quanxian/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var user models.User
	if err := utils.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "账号已被禁用"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	utils.DB.Preload("Roles").Preload("Department").First(&user, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func GetUserInfo(c *gin.Context) {
	userId := c.GetUint("userId")

	var user models.User
	if err := utils.DB.Preload("Roles").Preload("Department").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	menus, err := utils.GetUserMenus(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取菜单失败"})
		return
	}

	buttons := utils.GetUserAllButtonCodes(userId)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"user":    user,
			"menus":   menus,
			"buttons": buttons,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "退出成功"})
}
