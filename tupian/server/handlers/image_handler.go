package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"image-editor-server/database"
	"image-editor-server/models"
	"image-editor-server/storage"

	"github.com/gin-gonic/gin"
)

type UploadResponse struct {
	Success []*models.Image `json:"success"`
	Failed  []string        `json:"failed"`
}

func UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		Fail(c, http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}

	if len(files) == 0 {
		Fail(c, http.StatusBadRequest, "No files uploaded")
		return
	}

	response := &UploadResponse{
		Success: make([]*models.Image, 0),
		Failed:  make([]string, 0),
	}

	for _, file := range files {
		savedFile, err := storage.SaveFile(file)
		if err != nil {
			response.Failed = append(response.Failed, file.Filename)
			continue
		}

		image := &models.Image{
			UUID:         savedFile.UUID,
			Filename:     savedFile.Filename,
			OriginalName: savedFile.OriginalName,
			FilePath:     savedFile.FilePath,
			FileSize:     savedFile.FileSize,
			ContentType:  savedFile.ContentType,
		}

		if err := database.DB.Create(image).Error; err != nil {
			storage.DeleteFile(savedFile.FilePath)
			response.Failed = append(response.Failed, file.Filename)
			continue
		}

		response.Success = append(response.Success, image)
	}

	Success(c, response)
}

func GetImages(c *gin.Context) {
	var images []models.Image
	if err := database.DB.Order("created_at DESC").Find(&images).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "Failed to fetch images: "+err.Error())
		return
	}
	Success(c, images)
}

func GetImage(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		Fail(c, http.StatusBadRequest, "Image UUID is required")
		return
	}

	var image models.Image
	if err := database.DB.Where("uuid = ?", uuid).First(&image).Error; err != nil {
		FailWithStatus(c, http.StatusNotFound, http.StatusNotFound, "Image not found")
		return
	}

	accept := c.GetHeader("Accept")
	if accept == "application/json" || c.Query("format") == "json" {
		Success(c, image)
		return
	}

	filePath := filepath.Join(storage.UploadDir, image.Filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		FailWithStatus(c, http.StatusNotFound, http.StatusNotFound, "File not found on disk")
		return
	}

	c.File(filePath)
}

func DeleteImage(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		Fail(c, http.StatusBadRequest, "Image UUID is required")
		return
	}

	var image models.Image
	if err := database.DB.Where("uuid = ?", uuid).First(&image).Error; err != nil {
		FailWithStatus(c, http.StatusNotFound, http.StatusNotFound, "Image not found")
		return
	}

	storage.DeleteFile(image.FilePath)

	if err := database.DB.Delete(&image).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "Failed to delete image record: "+err.Error())
		return
	}

	Success(c, gin.H{"uuid": uuid})
}
