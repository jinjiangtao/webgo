package middleware

import (
	"net/http"
	"voting-system/config"
	"voting-system/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AntiRepeatVote() gin.HandlerFunc {
	return func(c *gin.Context) {
		activityID := c.GetUint("activity_id")
		if activityID == 0 {
			c.Next()
			return
		}

		ip := c.ClientIP()
		cookie, err := c.Cookie("voting_id")
		if err != nil {
			cookie = uuid.New().String()
			c.SetCookie("voting_id", cookie, 365*24*3600, "/", "", false, true)
		}

		c.Set("client_ip", ip)
		c.Set("client_cookie", cookie)

		var count int64
		config.DB.Model(&models.VoteRecord{}).Where(
			"activity_id = ? AND (ip = ? OR cookie = ?)",
			activityID, ip, cookie,
		).Count(&count)

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "您已经投过票了，请勿重复投票",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckActivityStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		activityID := c.GetUint("activity_id")
		if activityID == 0 {
			id := c.Param("id")
			var activity models.Activity
			if err := config.DB.First(&activity, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"code":    404,
					"message": "活动不存在",
				})
				c.Abort()
				return
			}
			c.Set("activity_id", activity.ID)
			c.Set("activity", activity)
		}
		c.Next()
	}
}
