package handlers

import (
	"daka/server/database"
	"daka/server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := database.DB.Where("is_active = ? OR is_active = ?", true, false).Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")
	result := make([]models.TaskWithStatus, len(tasks))

	for i, task := range tasks {
		var record models.Record
		database.DB.Where("task_id = ? AND record_date = ?", task.ID, today).First(&record)
		result[i] = models.TaskWithStatus{
			Task: task,
		}
		if record.ID > 0 {
			result[i].TodayRecord = &record
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func CreateTask(c *gin.Context) {
	var req models.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		Name:             req.Name,
		Description:      req.Description,
		CycleType:        req.CycleType,
		CountdownSeconds: req.CountdownSeconds,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Color:            req.Color,
		IsActive:         req.IsActive,
	}

	if task.Color == "" {
		task.Color = "#6366f1"
	}
	if req.IsActive == false && req.Name != "" {
		task.IsActive = req.IsActive
	} else if req.Name != "" {
		task.IsActive = true
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var req models.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.CycleType != "" {
		task.CycleType = req.CycleType
	}
	if req.CountdownSeconds > 0 {
		task.CountdownSeconds = req.CountdownSeconds
	}
	if req.StartTime != "" {
		task.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		task.EndTime = req.EndTime
	}
	if req.Color != "" {
		task.Color = req.Color
	}
	if req.IsActive != nil {
		task.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	database.DB.Where("task_id = ?", id).Delete(&models.Record{})

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
