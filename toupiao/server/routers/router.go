package routers

import (
	"voting-system/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/dashboard", controllers.GetDashboardStats)

		activity := api.Group("/activity")
		{
			activity.POST("", controllers.CreateActivity)
			activity.GET("", controllers.GetActivityList)
			activity.GET("/:id", controllers.GetActivityDetail)
			activity.PUT("/:id/status", controllers.ToggleActivityStatus)
			activity.DELETE("/:id", controllers.DeleteActivity)
		}

		vote := api.Group("/vote")
		{
			vote.POST("", controllers.SubmitVote)
		}

		stats := api.Group("/stats")
		{
			stats.GET("/:id", controllers.GetVoteStats)
			stats.GET("/:id/records", controllers.GetVoteRecords)
		}
	}

	return r
}
