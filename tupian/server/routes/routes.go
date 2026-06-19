package routes

import (
	"image-editor-server/handlers"
	"image-editor-server/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.Static("/uploads", "./"+storage.UploadDir)

	api := r.Group("/api")
	{
		images := api.Group("/images")
		{
			images.POST("/upload", handlers.UploadImages)
			images.GET("", handlers.GetImages)
			images.GET("/:uuid", handlers.GetImage)
			images.DELETE("/:uuid", handlers.DeleteImage)
		}

		editLogs := api.Group("/edit-logs")
		{
			editLogs.POST("", handlers.CreateEditLog)
			editLogs.GET("/:imageUuid", handlers.GetEditLogs)
		}
	}

	return r
}
