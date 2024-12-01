package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krishnakumarkp/to-do/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskService is a mock implementation of the TaskServiceInterface
type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(title, description string) (domain.Task, error) {
	args := m.Called(title, description)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskService) GetTask(id uint) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskService) GetAllTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskService) UpdateTask(id uint, task domain.Task) (domain.Task, error) {
	args := m.Called(id, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskService) MarkTaskCompleted(id uint) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskService) DeleteTask(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	task := domain.Task{Title: "Test Task", Description: "Test Description"}
	mockService.On("CreateTask", task.Title, task.Description).Return(task, nil)

	router := gin.Default()
	router.POST("/tasks", handler.CreateTask)

	taskJSON, _ := json.Marshal(task)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, task.Title, response["title"])
	mockService.AssertCalled(t, "CreateTask", task.Title, task.Description)
}

func TestGetTaskByID(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	task := domain.Task{ID: 1, Title: "Test Task"}
	mockService.On("GetTask", uint(1)).Return(task, nil)

	router := gin.Default()
	router.GET("/tasks/:id", handler.GetTaskByID)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response domain.Task
	_ = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Equal(t, task.Title, response.Title)
	mockService.AssertCalled(t, "GetTask", uint(1))
}

func TestDeleteTask(t *testing.T) {
	mockService := new(MockTaskService)
	handler := NewTaskHandler(mockService)

	mockService.On("DeleteTask", uint(1)).Return(nil)

	router := gin.Default()
	router.DELETE("/tasks/:id", handler.DeleteTask)

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	mockService.AssertCalled(t, "DeleteTask", uint(1))
}
