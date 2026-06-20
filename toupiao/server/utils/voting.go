package utils

import (
	"voting-system/config"
	"voting-system/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetClientInfo(c *gin.Context) (string, string) {
	ip := c.ClientIP()
	cookie, err := c.Cookie("voting_id")
	if err != nil {
		cookie = uuid.New().String()
		c.SetCookie("voting_id", cookie, 365*24*3600, "/", "", false, true)
	}
	return ip, cookie
}

func CheckDuplicateVote(activityID uint, ip, cookie string) bool {
	var count int64
	config.DB.Model(&models.VoteRecord{}).Where(
		"activity_id = ? AND (ip = ? OR cookie = ?)",
		activityID, ip, cookie,
	).Count(&count)
	return count > 0
}
