package main

import (
	"dashboard/internal/controller"
	"dashboard/internal/middleware"
	"dashboard/internal/repository"
	"dashboard/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	seedService := service.NewDataSeedService()
	if err := seedService.SeedAllData(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}
	log.Println("Data initialization completed")

	r := gin.Default()

	r.Use(middleware.CORS())

	dataController := controller.NewDataController()

	api := r.Group("/api")
	{
		api.GET("/health", dataController.Health)
		api.GET("/dimensions", dataController.GetDimensions)

		data := api.Group("/data")
		{
			data.POST("/aggregate", dataController.Aggregate)
			data.POST("/drill", dataController.Drill)
			data.POST("/trace", dataController.Trace)
		}

		snapshots := api.Group("/snapshots")
		{
			snapshots.POST("", dataController.SaveSnapshot)
			snapshots.GET("", dataController.GetSnapshots)
			snapshots.GET("/:id", dataController.GetSnapshot)
		}

		export := api.Group("/export")
		{
			export.POST("/excel", dataController.ExportExcel)
			export.POST("/aggregate-excel", dataController.ExportAggregateExcel)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
