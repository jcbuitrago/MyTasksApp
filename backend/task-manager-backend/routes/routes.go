package routes

import (
	"github.com/gin-gonic/gin"
	"task-manager-backend/controllers"
)

func SetupRoutes(router *gin.Engine) {
	// User routes
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/", controllers.GetUsers)
		userGroup.GET("/:id", controllers.GetUserByID)
		userGroup.PUT("/:id", controllers.UpdateUser)
		userGroup.DELETE("/:id", controllers.DeleteUser)
	}

	// Category routes
	categoryGroup := router.Group("/api/categories")
	{
		categoryGroup.POST("/", controllers.CreateCategory)
		categoryGroup.GET("/", controllers.GetCategories)
		categoryGroup.GET("/:id", controllers.GetCategoryByID)
		categoryGroup.PUT("/:id", controllers.UpdateCategory)
		categoryGroup.DELETE("/:id", controllers.DeleteCategory)
	}

	// Task routes
	taskGroup := router.Group("/api/tasks")
	{
		taskGroup.POST("/", controllers.CreateTask)
		taskGroup.GET("/", controllers.GetTasks)
		taskGroup.GET("/:id", controllers.GetTaskByID)
		taskGroup.PUT("/:id", controllers.UpdateTask)
		taskGroup.DELETE("/:id", controllers.DeleteTask)
	}
}