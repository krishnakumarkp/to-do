package http

import "github.com/gin-gonic/gin"

// TaskHandlerInterface defines the contract for task handler operations.
type TaskHandlerInterface interface {
	CreateTask(c *gin.Context)
	GetAllTasks(c *gin.Context)
	GetTaskByID(c *gin.Context)
	UpdateTask(c *gin.Context)
	MarkTaskAsDone(c *gin.Context)
	DeleteTask(c *gin.Context)
}
