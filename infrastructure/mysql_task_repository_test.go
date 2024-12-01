package infrastructure

import (
	"errors"
	"testing"
	"time"

	"github.com/krishnakumarkp/to-do/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSave(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()

	// Create GORM DB from sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm: %v", err)
	}

	// Mock repository
	repo := NewMySQLTaskRepository(gormDB)

	// Define the task to save
	task := domain.Task{
		Title:       "Test Task",
		Description: "Test description",
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	// Set up expectations
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `tasks`").WithArgs(task.Title, task.Description, task.Completed, task.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Execute the function
	id, err := repo.Save(task)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Assert the returned ID
	if id != 1 {
		t.Errorf("expected ID to be 1, got %d", id)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestFindByID(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()

	// Create GORM DB with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable GORM logs for cleaner test output
	})
	if err != nil {
		t.Fatalf("failed to initialize gorm: %v", err)
	}

	// Mock repository
	repo := NewMySQLTaskRepository(gormDB)

	// Define test cases
	tests := []struct {
		name        string
		taskID      uint
		mockSetup   func()
		expectedErr error
		expectedRes domain.Task
	}{
		{
			name:   "Task Found",
			taskID: 1,
			mockSetup: func() {
				// Mock the SQL query
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "Test Task", "Test description")
				mock.ExpectQuery("^SELECT \\* FROM `tasks` WHERE `tasks`.`id` = \\? ORDER BY `tasks`.`id` LIMIT \\?$").
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			expectedErr: nil,
			expectedRes: domain.Task{ID: 1, Title: "Test Task", Description: "Test description"},
		},
		{
			name:   "Task Not Found",
			taskID: 2,
			mockSetup: func() {
				// Mock the SQL query to return no rows
				mock.ExpectQuery("^SELECT \\* FROM `tasks` WHERE `tasks`.`id` = \\? ORDER BY `tasks`.`id` LIMIT \\?$").
					WithArgs(2, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedErr: errors.New("task not found"),
			expectedRes: domain.Task{},
		},
		{
			name:   "Database Error",
			taskID: 3,
			mockSetup: func() {
				// Mock the SQL query to simulate a database error
				mock.ExpectQuery("^SELECT \\* FROM `tasks` WHERE `tasks`.`id` = \\? ORDER BY `tasks`.`id` LIMIT \\?$").
					WithArgs(3, 1).
					WillReturnError(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
			expectedRes: domain.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock behavior
			tt.mockSetup()

			// Call the method
			result, err := repo.FindByID(tt.taskID)

			// Assert the results
			if err != nil && tt.expectedErr == nil || err == nil && tt.expectedErr != nil || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if result != tt.expectedRes {
				t.Errorf("expected result: %+v, got: %+v", tt.expectedRes, result)
			}

			// Verify that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	defer db.Close()

	// Create GORM DB from sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable GORM logs for test clarity
	})
	if err != nil {
		t.Fatalf("failed to initialize gorm: %v", err)
	}

	// Mock repository
	repo := NewMySQLTaskRepository(gormDB)

	// Define test cases
	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr error
		expectedRes []domain.Task
	}{
		{
			name: "Tasks Found",
			mockSetup: func() {
				// Mock SQL query to return multiple rows
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "Test Task 1", "Test description 1").
					AddRow(2, "Test Task 2", "Test description 2")
				mock.ExpectQuery("^SELECT \\* FROM `tasks`").
					WillReturnRows(rows)
			},
			expectedErr: nil,
			expectedRes: []domain.Task{
				{ID: 1, Title: "Test Task 1", Description: "Test description 1"},
				{ID: 2, Title: "Test Task 2", Description: "Test description 2"},
			},
		},
		{
			name: "No Tasks Found",
			mockSetup: func() {
				// Mock SQL query to return no rows
				rows := sqlmock.NewRows([]string{"id", "title", "description"})
				mock.ExpectQuery("^SELECT \\* FROM `tasks`").
					WillReturnRows(rows)
			},
			expectedErr: nil,
			expectedRes: []domain.Task{},
		},
		{
			name: "Database Error",
			mockSetup: func() {
				// Mock SQL query to simulate a database error
				mock.ExpectQuery("^SELECT \\* FROM `tasks`").
					WillReturnError(gorm.ErrInvalidTransaction)
			},
			expectedErr: gorm.ErrInvalidTransaction,
			expectedRes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock behavior
			tt.mockSetup()

			// Call the method
			result, err := repo.FindAll()

			// Assert the results
			if err != nil && tt.expectedErr == nil || err == nil && tt.expectedErr != nil || (err != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if len(result) != len(tt.expectedRes) {
				t.Errorf("expected result length: %d, got: %d", len(tt.expectedRes), len(result))
			}
			for i := range result {
				if result[i] != tt.expectedRes[i] {
					t.Errorf("expected task at index %d: %+v, got: %+v", i, tt.expectedRes[i], result[i])
				}
			}

			// Verify that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

// func TestUpdate(t *testing.T) {
// 	// Initialize sqlmock
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("failed to open mock db: %v", err)
// 	}
// 	defer db.Close()

// 	// Create GORM DB from sqlmock
// 	gormDB, err := gorm.Open(mysql.New(mysql.Config{
// 		Conn:                      db,
// 		SkipInitializeWithVersion: true,
// 	}), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Silent), // Disable GORM logs for cleaner test output
// 	})
// 	if err != nil {
// 		t.Fatalf("failed to initialize gorm: %v", err)
// 	}

// 	// Mock repository
// 	repo := NewMySQLTaskRepository(gormDB)

// 	// Define test cases
// 	tests := []struct {
// 		name        string
// 		inputTask   domain.Task
// 		mockSetup   func()
// 		expectedErr error
// 		expectedRes domain.Task
// 	}{
// 		{
// 			name: "Successful Update",
// 			inputTask: domain.Task{
// 				ID:    1,
// 				Title: "Updated Task Title",
// 			},
// 			mockSetup: func() {
// 				mock.ExpectBegin()
// 				mock.ExpectExec("^UPDATE `tasks` SET `title`= \\?,,`description`= \\?,`completed`= \\?,`created_at`= \\? WHERE `id` = \\?$").
// 					WithArgs("Updated Task Title", "Updated Task description", false, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
// 					WillReturnResult(sqlmock.NewResult(1, 1))
// 				mock.ExpectCommit()
// 			},
// 			expectedErr: nil,
// 			expectedRes: domain.Task{
// 				ID:    1,
// 				Title: "Updated Task Name",
// 			},
// 		},
// 		{
// 			name: "Task Not Found",
// 			inputTask: domain.Task{
// 				ID:   2,
// 				Name: "Nonexistent Task",
// 			},
// 			mockSetup: func() {
// 				mock.ExpectBegin()
// 				mock.ExpectExec("^UPDATE `tasks` SET `name`=\\?,`updated_at`=\\? WHERE `id` = \\?$").
// 					WithArgs("Nonexistent Task", sqlmock.AnyArg(), 2).
// 					WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
// 				mock.ExpectCommit()
// 			},
// 			expectedErr: nil,
// 			expectedRes: domain.Task{
// 				ID:   2,
// 				Name: "Nonexistent Task",
// 			},
// 		},
// 		{
// 			name: "Database Error",
// 			inputTask: domain.Task{
// 				ID:   3,
// 				Name: "Task with DB Error",
// 			},
// 			mockSetup: func() {
// 				mock.ExpectBegin()
// 				mock.ExpectExec("^UPDATE `tasks` SET `name`=\\?,`updated_at`=\\? WHERE `id` = \\?$").
// 					WithArgs("Task with DB Error", sqlmock.AnyArg(), 3).
// 					WillReturnError(gorm.ErrInvalidTransaction) // Simulate DB error
// 				mock.ExpectRollback()
// 			},
// 			expectedErr: gorm.ErrInvalidTransaction,
// 			expectedRes: domain.Task{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Set up the mock behavior
// 			tt.mockSetup()

// 			// Call the method
// 			result, err := repo.Update(tt.inputTask)

// 			// Assert the results
// 			if err != nil && tt.expectedErr == nil || err == nil && tt.expectedErr != nil || (err != nil && err.Error() != tt.expectedErr.Error()) {
// 				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
// 			}
// 			if result != tt.expectedRes {
// 				t.Errorf("expected result: %+v, got: %+v", tt.expectedRes, result)
// 			}

// 			// Verify that all expectations were met
// 			if err := mock.ExpectationsWereMet(); err != nil {
// 				t.Errorf("unmet expectations: %v", err)
// 			}
// 		})
// 	}
// }
