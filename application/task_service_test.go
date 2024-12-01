package application

import (
	"errors"
	"testing"
	"time"

	"github.com/krishnakumarkp/to-do/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository is a mock implementation of the TaskRepository interface
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Save(task domain.Task) (uint, error) {
	args := m.Called(task)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockTaskRepository) FindByID(id uint) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindAll() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task domain.Task) (domain.Task, error) {
	args := m.Called(task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	task := domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Completed:   false,
		CreatedAt:   time.Now(),
	}
	mockRepo.On("Save", mock.Anything).Return(uint(1), nil)

	result, err := service.CreateTask(task.Title, task.Description)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, task.Title, result.Title)
	mockRepo.AssertCalled(t, "Save", mock.Anything)
}

func TestGetTask_Success(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	task := domain.Task{ID: 1, Title: "Test Task", Description: "Test Description"}
	mockRepo.On("FindByID", uint(1)).Return(task, nil)

	result, err := service.GetTask(1)

	assert.NoError(t, err)
	assert.Equal(t, task, result)
	mockRepo.AssertCalled(t, "FindByID", uint(1))
}

func TestGetTask_NotFound(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(domain.Task{}, errors.New("not found"))

	_, err := service.GetTask(1)

	assert.Error(t, err)
	mockRepo.AssertCalled(t, "FindByID", uint(1))
}

func TestGetAllTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	tasks := []domain.Task{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}
	mockRepo.On("FindAll").Return(tasks, nil)

	result, err := service.GetAllTasks()

	assert.NoError(t, err)
	assert.Equal(t, tasks, result)
	mockRepo.AssertCalled(t, "FindAll")
}

func TestMarkTaskCompleted(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	task := domain.Task{ID: 1, Title: "Test Task", Completed: false}
	updatedTask := task
	updatedTask.Completed = true

	mockRepo.On("FindByID", uint(1)).Return(task, nil)
	mockRepo.On("Update", updatedTask).Return(updatedTask, nil)

	result, err := service.MarkTaskCompleted(1)

	assert.NoError(t, err)
	assert.True(t, result.Completed)
	mockRepo.AssertCalled(t, "FindByID", uint(1))
	mockRepo.AssertCalled(t, "Update", updatedTask)
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := NewTaskService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteTask(1)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", uint(1))
}
