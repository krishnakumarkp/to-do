package domain

import "time"

// Task represents a domain entity for a to-do item
type Task struct {
	ID          uint      `json:"id"`          // Unique identifier
	Title       string    `json:"title"`       // Title of the task
	Description string    `json:"description"` // Detailed description of the task
	Completed   bool      `json:"completed"`   // Task completion status
	CreatedAt   time.Time `json:"created_at"`  // Timestamp of task creation
}

// TaskRepository is an interface for interacting with task storage
type TaskRepository interface {
	Save(task Task) (uint, error)
	FindByID(id uint) (Task, error)
	FindAll() ([]Task, error)
	Update(task Task) (Task, error)
	Delete(id uint) error
}
