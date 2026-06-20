package main

import (
	"jiansuo/database"
	"jiansuo/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	log.Println("Database initialized successfully")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
	}))

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Jiansuo search service is running",
		})
	})

	api := r.Group("/api")
	{
		search := api.Group("/search")
		{
			search.GET("", handlers.Search)
			search.GET("/suggest", handlers.Suggest)
			search.GET("/hot", handlers.HotKeywords)
			search.GET("/stats", handlers.GetStatistics)
		}

		keywords := api.Group("/keywords")
		{
			keywords.GET("", handlers.ListKeywords)
			keywords.GET("/:id", handlers.KeywordDetail)
			keywords.POST("", handlers.CreateKeyword)
			keywords.PUT("/:id", handlers.UpdateKeyword)
			keywords.DELETE("/:id", handlers.DeleteKeyword)
			keywords.PATCH("/:id/status", handlers.SetKeywordStatus)
			keywords.POST("/:id/view", handlers.IncrementView)
			keywords.POST("/batch", handlers.BatchKeywords)
			keywords.POST("/import", handlers.ImportCSV)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", handlers.ListCategories)
			categories.POST("", handlers.CreateCategory)
			categories.PUT("/:id", handlers.UpdateCategory)
			categories.DELETE("/:id", handlers.DeleteCategory)
			categories.PATCH("/:id/status", handlers.SetCategoryStatus)
		}

		logs := api.Group("/logs")
		{
			logs.GET("", handlers.ListSearchLogs)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
