package infrastructure

import (
	"errors"

	"github.com/krishnakumarkp/to-do/domain"

	"gorm.io/gorm"
)

type MySQLTaskRepository struct {
	db *gorm.DB
}

func NewMySQLTaskRepository(db *gorm.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

func (r *MySQLTaskRepository) Save(task domain.Task) (uint, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return 0, result.Error
	}
	return task.ID, nil
}

func (r *MySQLTaskRepository) FindByID(id uint) (domain.Task, error) {
	var task domain.Task
	result := r.db.First(&task, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return task, errors.New("task not found")
	}
	return task, result.Error
}

func (r *MySQLTaskRepository) FindAll() ([]domain.Task, error) {
	var tasks []domain.Task
	result := r.db.Find(&tasks)
	return tasks, result.Error
}

// UpdateTask updates a task in the database
func (r *MySQLTaskRepository) Update(task domain.Task) (domain.Task, error) {
	// Use GORM's Save method to update the task in the database
	if err := r.db.Save(&task).Error; err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (r *MySQLTaskRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Task{}, id)
	if result.RowsAffected == 0 {
		return errors.New("task not found")
	}
	return result.Error
}
