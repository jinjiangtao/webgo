package controllers

import (
	"fmt"
	"net/http"
	"time"
	"voting-system/config"
	"voting-system/models"
	"voting-system/utils"

	"github.com/gin-gonic/gin"
)

type VoteRequest struct {
	ActivityID uint   `json:"activity_id" binding:"required"`
	OptionIDs  []uint `json:"option_ids" binding:"required,min=1"`
}

func SubmitVote(c *gin.Context) {
	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	var activity models.Activity
	if err := config.DB.First(&activity, req.ActivityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
		})
		return
	}

	now := time.Now()
	if now.Before(activity.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "活动尚未开始",
		})
		return
	}
	if now.After(activity.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "活动已结束",
		})
		return
	}

	if activity.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "活动已关闭",
		})
		return
	}

	if activity.VoteType == "single" && len(req.OptionIDs) > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "该活动为单选投票，只能选择一个选项",
		})
		return
	}

	if activity.VoteType == "multiple" && len(req.OptionIDs) > activity.MaxChoices {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "最多只能选择 " + fmt.Sprintf("%d", activity.MaxChoices) + " 个选项",
		})
		return
	}

	ip, cookie := utils.GetClientInfo(c)

	if utils.CheckDuplicateVote(req.ActivityID, ip, cookie) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "您已经投过票了，请勿重复投票",
		})
		return
	}

	var validOptions []models.Option
	config.DB.Where("id IN ? AND activity_id = ?", req.OptionIDs, req.ActivityID).Find(&validOptions)
	if len(validOptions) != len(req.OptionIDs) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "存在无效的投票选项",
		})
		return
	}

	tx := config.DB.Begin()

	for _, opt := range validOptions {
		if err := tx.Model(&opt).Update("vote_count", opt.VoteCount+1).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新票数失败",
			})
			return
		}

		record := models.VoteRecord{
			ActivityID: req.ActivityID,
			OptionID:   opt.ID,
			IP:         ip,
			Cookie:     cookie,
			OptionName: opt.Name,
		}
		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "记录投票失败",
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "投票成功",
	})
}
