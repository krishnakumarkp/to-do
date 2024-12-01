package application

import (
	"github.com/krishnakumarkp/to-do/domain"

	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(title, description string) (domain.Task, error) {
	args := m.Called(title, description)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *MockTaskService) GetAllTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskService) GetTask(id uint) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
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
