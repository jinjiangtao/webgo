package controllers

import (
	"net/http"
	"voting-system/config"
	"voting-system/models"
	"voting-system/utils"

	"github.com/gin-gonic/gin"
)

type CreateActivityRequest struct {
	Title       string        `json:"title" binding:"required"`
	Description string        `json:"description"`
	StartTime   utils.DateTime `json:"start_time"`
	EndTime     utils.DateTime `json:"end_time"`
	VoteType    string        `json:"vote_type" binding:"required,oneof=single multiple"`
	MaxChoices  int           `json:"max_choices"`
	Options     []string      `json:"options" binding:"required,min=2"`
}

func CreateActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if req.StartTime.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "开始时间不能为空",
		})
		return
	}
	if req.EndTime.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "结束时间不能为空",
		})
		return
	}
	if !req.EndTime.After(req.StartTime.Time) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "结束时间必须晚于开始时间",
		})
		return
	}

	if req.MaxChoices <= 0 {
		req.MaxChoices = 1
	}
	if req.VoteType == "single" {
		req.MaxChoices = 1
	}

	activity := models.Activity{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime.Time,
		EndTime:     req.EndTime.Time,
		VoteType:    req.VoteType,
		MaxChoices:  req.MaxChoices,
		Status:      1,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&activity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建活动失败",
		})
		return
	}

	for _, optName := range req.Options {
		option := models.Option{
			ActivityID: activity.ID,
			Name:       optName,
		}
		if err := tx.Create(&option).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建选项失败",
			})
			return
		}
	}

	tx.Commit()

	config.DB.Preload("Options").First(&activity, activity.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    activity,
	})
}

func GetActivityList(c *gin.Context) {
	var activities []models.Activity
	config.DB.Preload("Options").Order("created_at desc").Find(&activities)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    activities,
	})
}

func GetActivityDetail(c *gin.Context) {
	id := c.Param("id")
	var activity models.Activity
	if err := config.DB.Preload("Options").First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    activity,
	})
}

func ToggleActivityStatus(c *gin.Context) {
	id := c.Param("id")
	var activity models.Activity
	if err := config.DB.First(&activity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
		})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	activity.Status = req.Status
	config.DB.Save(&activity)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "状态更新成功",
		"data":    activity,
	})
}

func DeleteActivity(c *gin.Context) {
	id := c.Param("id")

	tx := config.DB.Begin()
	if err := tx.Where("activity_id = ?", id).Delete(&models.VoteRecord{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除投票记录失败",
		})
		return
	}

	if err := tx.Where("activity_id = ?", id).Delete(&models.Option{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除选项失败",
		})
		return
	}

	if err := tx.Delete(&models.Activity{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除活动失败",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
