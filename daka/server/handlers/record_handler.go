package handlers

import (
	"daka/server/database"
	"daka/server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func determineStatus(task models.Task, isMakeup bool, now time.Time) string {
	if isMakeup {
		return "makeup"
	}

	if task.StartTime != "" && task.EndTime != "" {
		currentTime := now.Format("15:04")
		if currentTime >= task.StartTime && currentTime <= task.EndTime {
			return "on_time"
		}
		if currentTime > task.EndTime {
			return "late"
		}
	}

	return "checked_in"
}

func CheckIn(c *gin.Context) {
	var req models.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := database.DB.First(&task, req.TaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	now := time.Now()
	today := now.Format("2006-01-02")

	var existingRecord models.Record
	err := database.DB.Where("task_id = ? AND record_date = ?", req.TaskID, today).First(&existingRecord).Error

	if err == nil && existingRecord.ID > 0 && !req.IsMakeup {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already checked in today"})
		return
	}

	status := determineStatus(task, req.IsMakeup, now)

	if existingRecord.ID > 0 {
		existingRecord.Status = status
		existingRecord.IsMakeup = req.IsMakeup
		existingRecord.Note = req.Note
		existingRecord.Duration = req.Duration
		existingRecord.CheckInTime = &now
		database.DB.Save(&existingRecord)
		c.JSON(http.StatusOK, gin.H{"data": existingRecord})
		return
	}

	record := models.Record{
		TaskID:      req.TaskID,
		RecordDate:  today,
		Status:      status,
		IsMakeup:    req.IsMakeup,
		Note:        req.Note,
		Duration:    req.Duration,
		CheckInTime: &now,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": record})
}

func GetRecords(c *gin.Context) {
	taskID := c.Query("task_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := database.DB.Preload("Task")

	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}
	if startDate != "" {
		query = query.Where("record_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("record_date <= ?", endDate)
	}

	var records []models.Record
	if err := query.Order("record_date desc, id desc").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}

func GetRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.Record
	if err := database.DB.Preload("Task").First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": record})
}

func DeleteRecord(c *gin.Context) {
	id := c.Param("id")
	var record models.Record
	if err := database.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if err := database.DB.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

func MarkAbsent(c *gin.Context) {
	type MarkAbsentRequest struct {
		TaskID     uint   `json:"task_id" binding:"required"`
		RecordDate string `json:"record_date" binding:"required"`
	}

	var req MarkAbsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingRecord models.Record
	err := database.DB.Where("task_id = ? AND record_date = ?", req.TaskID, req.RecordDate).First(&existingRecord).Error

	if err == nil && existingRecord.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record already exists for this date"})
		return
	}

	record := models.Record{
		TaskID:     req.TaskID,
		RecordDate: req.RecordDate,
		Status:     "absent",
		IsMakeup:   false,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": record})
}
