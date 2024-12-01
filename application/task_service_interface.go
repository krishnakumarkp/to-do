package application

import "github.com/krishnakumarkp/to-do/domain"

// TaskService defines the methods for managing tasks.

type TaskServiceInterface interface {
	CreateTask(title, description string) (domain.Task, error)
	GetAllTasks() ([]domain.Task, error)
	GetTask(id uint) (domain.Task, error)
	UpdateTask(id uint, task domain.Task) (domain.Task, error)
	MarkTaskCompleted(id uint) (domain.Task, error)
	DeleteTask(id uint) error
}
