package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskHandler is a mock implementation of the TaskHandler
type MockTaskHandler struct {
	mock.Mock
}

func (m *MockTaskHandler) CreateTask(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Task created"})
}

func (m *MockTaskHandler) GetAllTasks(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "All tasks"})
}

func (m *MockTaskHandler) GetTaskByID(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Task by ID"})
}

func (m *MockTaskHandler) UpdateTask(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (m *MockTaskHandler) MarkTaskAsDone(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as done"})
}

func (m *MockTaskHandler) DeleteTask(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusNoContent, nil)
}

func TestSetupRouter(t *testing.T) {
	// Create a mock task handler
	mockHandler := new(MockTaskHandler)
	router := SetupRouter(mockHandler)

	// Define test cases
	tests := []struct {
		method       string
		path         string
		expectedCode int
		mockMethod   string
	}{
		{"POST", "/tasks", http.StatusOK, "CreateTask"},
		{"GET", "/tasks", http.StatusOK, "GetAllTasks"},
		{"GET", "/tasks/1", http.StatusOK, "GetTaskByID"},
		{"PUT", "/tasks/1", http.StatusOK, "UpdateTask"},
		{"PATCH", "/tasks/1/done", http.StatusOK, "MarkTaskAsDone"},
		{"DELETE", "/tasks/1", http.StatusNoContent, "DeleteTask"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			// Expect the mock method to be called
			mockHandler.On(tt.mockMethod, mock.Anything).Return().Once()

			// Create an HTTP request and response recorder
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			recorder := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(recorder, req)

			// Assert response code
			assert.Equal(t, tt.expectedCode, recorder.Code)

			// Assert the mock method was called
			mockHandler.AssertCalled(t, tt.mockMethod, mock.Anything)
		})
	}
}
