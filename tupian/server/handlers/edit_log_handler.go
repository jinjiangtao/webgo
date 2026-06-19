package handlers

import (
	"net/http"

	"image-editor-server/database"
	"image-editor-server/models"

	"github.com/gin-gonic/gin"
)

type CreateEditLogRequest struct {
	ImageUUID  string                 `json:"image_uuid" binding:"required"`
	ActionType string                 `json:"action_type" binding:"required"`
	ActionData map[string]interface{} `json:"action_data"`
}

func CreateEditLog(c *gin.Context) {
	var req CreateEditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	var image models.Image
	if err := database.DB.Where("uuid = ?", req.ImageUUID).First(&image).Error; err != nil {
		Fail(c, http.StatusBadRequest, "Image not found: "+req.ImageUUID)
		return
	}

	editLog := &models.EditLog{
		ImageUUID:  req.ImageUUID,
		ActionType: req.ActionType,
		ActionData: req.ActionData,
	}

	if err := database.DB.Create(editLog).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "Failed to create edit log: "+err.Error())
		return
	}

	Success(c, editLog)
}

func GetEditLogs(c *gin.Context) {
	imageUUID := c.Param("imageUuid")
	if imageUUID == "" {
		Fail(c, http.StatusBadRequest, "Image UUID is required")
		return
	}

	var logs []models.EditLog
	if err := database.DB.Where("image_uuid = ?", imageUUID).Order("created_at ASC").Find(&logs).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "Failed to fetch edit logs: "+err.Error())
		return
	}

	Success(c, logs)
}
