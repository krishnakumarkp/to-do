package infrastructure

import (
	"github.com/krishnakumarkp/to-do/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectToDB establishes a connection to the database using the DSN from the configuration
func ConnectToDB() (*gorm.DB, error) {
	// Get the DSN from the global config
	dsn := config.GetDSN()

	// Open a connection to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Return the DB instance
	return db, nil
}
