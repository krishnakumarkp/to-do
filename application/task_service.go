package application

import (
	"time"

	"github.com/krishnakumarkp/to-do/domain"
)

type TaskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title, description string) (domain.Task, error) {
	task := domain.Task{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	id, err := s.repo.Save(task)
	task.ID = id
	return task, err
}

// GetTask retrieves a task by its ID.
func (s *TaskService) GetTask(id uint) (domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

// GetAllTasks retrieves all tasks.
func (s *TaskService) GetAllTasks() ([]domain.Task, error) {
	return s.repo.FindAll()
}

func (s *TaskService) MarkTaskCompleted(id uint) (domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return task, err
	}
	task.Completed = true
	task, err = s.repo.Update(task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(id uint, task domain.Task) (domain.Task, error) {
	// Fetch the existing task by ID
	existingTask, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Task{}, err // If task doesn't exist, return error
	}

	// Update the task fields with the new data
	existingTask.Title = task.Title
	existingTask.Description = task.Description
	// Optionally, you can update other fields like Completed if necessary

	// Save the updated task
	updatedTask, err := s.repo.Update(existingTask)
	if err != nil {
		return domain.Task{}, err
	}

	return updatedTask, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}
