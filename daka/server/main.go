package main

import (
	"daka/server/database"
	"daka/server/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.POST("", handlers.CreateTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}

		records := api.Group("/records")
		{
			records.GET("", handlers.GetRecords)
			records.GET("/:id", handlers.GetRecord)
			records.DELETE("/:id", handlers.DeleteRecord)
		}

		api.POST("/check-in", handlers.CheckIn)
		api.POST("/mark-absent", handlers.MarkAbsent)

		stats := api.Group("/stats")
		{
			stats.GET("/calendar", handlers.GetCalendarData)
			stats.GET("/monthly", handlers.GetMonthlyStats)
			stats.GET("/task/:task_id", handlers.GetTaskStats)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
