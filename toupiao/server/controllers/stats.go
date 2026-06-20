package controllers

import (
	"net/http"
	"sort"
	"voting-system/config"
	"voting-system/models"

	"github.com/gin-gonic/gin"
)

func GetVoteStats(c *gin.Context) {
	activityID := c.Param("id")

	var activity models.Activity
	if err := config.DB.Preload("Options").First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
		})
		return
	}

	options := make([]models.Option, len(activity.Options))
	copy(options, activity.Options)

	sort.Slice(options, func(i, j int) bool {
		return options[i].VoteCount > options[j].VoteCount
	})

	totalVotes := 0
	for _, opt := range options {
		totalVotes += opt.VoteCount
	}

	ranking := make([]gin.H, len(options))
	for i, opt := range options {
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(opt.VoteCount) / float64(totalVotes) * 100
		}
		ranking[i] = gin.H{
			"id":         opt.ID,
			"name":       opt.Name,
			"image_url":  opt.ImageURL,
			"vote_count": opt.VoteCount,
			"percentage": percentage,
			"rank":       i + 1,
		}
	}

	chartData := make([]gin.H, len(options))
	for i, opt := range options {
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(opt.VoteCount) / float64(totalVotes) * 100
		}
		chartData[i] = gin.H{
			"name":       opt.Name,
			"value":      opt.VoteCount,
			"percentage": percentage,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"activity":    activity,
			"total_votes": totalVotes,
			"ranking":     ranking,
			"chart_data":  chartData,
		},
	})
}

func GetVoteRecords(c *gin.Context) {
	activityID := c.Param("id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	var activity models.Activity
	if err := config.DB.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
		})
		return
	}

	var records []models.VoteRecord
	var total int64

	query := config.DB.Model(&models.VoteRecord{}).Where("activity_id = ?", activityID)
	query.Count(&total)

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	query.Order("created_at desc").Offset(offset).Limit(parseInt(pageSize)).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"records": records,
			"total":   total,
			"page":    parseInt(page),
			"size":    parseInt(pageSize),
		},
	})
}

func GetDashboardStats(c *gin.Context) {
	var activityCount int64
	var totalVotes int64
	var todayVotes int64

	config.DB.Model(&models.Activity{}).Count(&activityCount)
	config.DB.Model(&models.VoteRecord{}).Count(&totalVotes)

	var activities []models.Activity
	config.DB.Preload("Options").Order("created_at desc").Limit(5).Find(&activities)

	recentActivities := make([]gin.H, len(activities))
	for i, act := range activities {
		total := 0
		for _, opt := range act.Options {
			total += opt.VoteCount
		}
		recentActivities[i] = gin.H{
			"id":          act.ID,
			"title":       act.Title,
			"status":      act.Status,
			"vote_type":   act.VoteType,
			"total_votes": total,
			"start_time":  act.StartTime,
			"end_time":    act.EndTime,
			"created_at":  act.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"activity_count":    activityCount,
			"total_votes":       totalVotes,
			"today_votes":       todayVotes,
			"recent_activities": recentActivities,
		},
	})
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	if result <= 0 {
		result = 1
	}
	return result
}
