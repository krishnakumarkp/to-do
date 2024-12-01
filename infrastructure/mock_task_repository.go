package infrastructure

import (
	"errors"
	"sync"

	"github.com/krishnakumarkp/to-do/domain"
)

type MemoryTaskRepository struct {
	tasks  map[uint]domain.Task
	mutex  sync.Mutex
	nextID uint
}

func NewMockTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks:  make(map[uint]domain.Task),
		nextID: 1, // Start IDs from 1
	}
}

func (r *MemoryTaskRepository) Save(task domain.Task) (uint, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if task.ID == 0 {
		task.ID = r.nextID
		r.nextID++
	}
	r.tasks[task.ID] = task
	return task.ID, nil
}

func (r *MemoryTaskRepository) FindByID(id uint) (domain.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return task, errors.New("task not found")
	}
	return task, nil
}

func (r *MemoryTaskRepository) FindAll() ([]domain.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	tasks := make([]domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *MemoryTaskRepository) Update(task domain.Task) (domain.Task, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingTask, exists := r.tasks[task.ID]
	if !exists {
		return existingTask, errors.New("task not found")
	}
	r.tasks[task.ID] = task
	return task, nil
}

func (r *MemoryTaskRepository) Delete(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return errors.New("task not found")
	}
	delete(r.tasks, id)
	return nil
}
