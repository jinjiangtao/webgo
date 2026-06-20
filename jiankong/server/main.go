package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dbPath := filepath.Join(".", "data.db")
	if envPath := os.Getenv("DB_PATH"); envPath != "" {
		dbPath = envPath
	}
	initDB(dbPath)
	defer db.Close()

	startSimulator()

	if os.Getenv("GIN_MODE") != "release" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	h := newHandler()
	api := r.Group("/api")
	{
		api.GET("/stats", h.GetStats)
		api.GET("/devices", h.GetDevices)
		api.GET("/devices/:id", h.GetDevice)
		api.POST("/devices", h.CreateDevice)
		api.PUT("/devices/:id", h.UpdateDevice)
		api.DELETE("/devices/:id", h.DeleteDevice)
		api.POST("/devices/:id/data", h.UpdateDeviceData)
		api.GET("/device-types", h.GetDeviceTypes)
		api.POST("/device-types", h.CreateDeviceType)
		api.DELETE("/device-types/:id", h.DeleteDeviceType)
		api.GET("/alarms", h.GetAlarms)
		api.POST("/alarms/:id/acknowledge", h.AcknowledgeAlarm)
		api.GET("/logs", h.GetLogs)
		api.DELETE("/logs", h.ClearLogs)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
