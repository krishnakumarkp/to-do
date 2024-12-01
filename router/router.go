package router

import (
	"github.com/krishnakumarkp/to-do/interfaces/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and returns the Gin router with all the routes
func SetupRouter(taskHandler http.TaskHandlerInterface) *gin.Engine {
	router := gin.Default()

	// Define routes
	router.POST("/tasks", taskHandler.CreateTask)               // Route to create a task
	router.GET("/tasks", taskHandler.GetAllTasks)               // Route to get all tasks
	router.GET("/tasks/:id", taskHandler.GetTaskByID)           // Route to get task by ID
	router.PUT("/tasks/:id", taskHandler.UpdateTask)            // Route to update task by ID
	router.PATCH("/tasks/:id/done", taskHandler.MarkTaskAsDone) // Route to mark task as done
	router.DELETE("/tasks/:id", taskHandler.DeleteTask)         // Route to delete

	return router
}
