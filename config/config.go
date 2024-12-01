package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration values for the application
type Config struct {
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	DBCharset   string
	DBParseTime string
	DBLoc       string
}

// Global variable to hold the loaded config
var AppConfig *Config

// LoadConfig loads the configuration from the .env file
func LoadConfig() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Read configuration values from environment variables
	AppConfig = &Config{
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBName:      os.Getenv("DB_NAME"),
		DBCharset:   os.Getenv("DB_CHARSET"),
		DBParseTime: os.Getenv("DB_PARSE_TIME"),
		DBLoc:       os.Getenv("DB_LOC"),
	}

	// Ensure that all required values are set
	if AppConfig.DBUser == "" || AppConfig.DBPassword == "" || AppConfig.DBHost == "" || AppConfig.DBPort == "" || AppConfig.DBName == "" {
		return fmt.Errorf("missing required environment variables")
	}

	return nil
}

// GetDSN returns the MySQL connection string (DSN) based on the loaded config
func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBHost, AppConfig.DBPort, AppConfig.DBName, AppConfig.DBCharset, AppConfig.DBParseTime, AppConfig.DBLoc,
	)
}
