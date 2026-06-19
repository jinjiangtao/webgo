package main

import (
	"log"

	"quanxian/handlers"
	"quanxian/middleware"
	"quanxian/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitDB()
	utils.InitSchema()
	utils.SeedData()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.GET("/logout", middleware.AuthMiddleware(), handlers.Logout)
			auth.GET("/userinfo", middleware.AuthMiddleware(), handlers.GetUserInfo)
		}

		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			departments := authorized.Group("/departments")
			{
				departments.GET("/tree", handlers.GetDepartmentTree)
				departments.GET("", handlers.GetDepartmentList)
				departments.GET("/:id", handlers.GetDepartment)
				departments.POST("", handlers.CreateDepartment)
				departments.PUT("/:id", handlers.UpdateDepartment)
				departments.DELETE("/:id", handlers.DeleteDepartment)
				departments.POST("/move", handlers.MoveDepartment)
			}

			roles := authorized.Group("/roles")
			{
				roles.GET("", handlers.GetRoleList)
				roles.GET("/:id", handlers.GetRole)
				roles.POST("", handlers.CreateRole)
				roles.PUT("/:id", handlers.UpdateRole)
				roles.DELETE("/:id", handlers.DeleteRole)
				roles.GET("/:id/menus", handlers.GetRoleMenus)
				roles.GET("/:id/buttons", handlers.GetRoleButtons)
				roles.POST("/:id/menus", handlers.BindRoleMenus)
				roles.POST("/:id/buttons", handlers.BindRoleButtons)
			}

			menus := authorized.Group("/menus")
			{
				menus.GET("/tree", handlers.GetMenuTree)
				menus.GET("", handlers.GetMenuList)
				menus.GET("/:id", handlers.GetMenu)
				menus.POST("", handlers.CreateMenu)
				menus.PUT("/:id", handlers.UpdateMenu)
				menus.DELETE("/:id", handlers.DeleteMenu)
				menus.POST("/move", handlers.MoveMenu)
			}

			buttons := authorized.Group("/buttons")
			{
				buttons.GET("", handlers.GetButtonList)
				buttons.GET("/:id", handlers.GetButton)
				buttons.POST("", handlers.CreateButton)
				buttons.PUT("/:id", handlers.UpdateButton)
				buttons.DELETE("/:id", handlers.DeleteButton)
			}

			users := authorized.Group("/users")
			{
				users.GET("", handlers.GetUserList)
				users.GET("/:id", handlers.GetUser)
				users.POST("", handlers.CreateUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
				users.POST("/:id/reset-password", handlers.ResetPassword)
			}

			permissions := authorized.Group("/permissions")
			{
				permissions.GET("/tree", handlers.GetAllPermissionsTree)
				permissions.GET("/preview", handlers.GetPermissionPreview)
				permissions.POST("/batch-assign", handlers.BatchAssignPermissions)
			}
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
