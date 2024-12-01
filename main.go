package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krishnakumarkp/to-do/application"
	"github.com/krishnakumarkp/to-do/config"
	"github.com/krishnakumarkp/to-do/domain"
	"github.com/krishnakumarkp/to-do/infrastructure"
	httpHandler "github.com/krishnakumarkp/to-do/interfaces/http"
	"github.com/krishnakumarkp/to-do/router"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the database
	db, err := infrastructure.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the Task model
	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repository, service, and handler
	repo := infrastructure.NewMySQLTaskRepository(db)
	//repo := infrastructure.NewMockTaskRepository()
	service := application.NewTaskService(repo)
	taskHandler := httpHandler.NewTaskHandler(service)

	// Set up the router using the router package
	router := router.SetupRouter(taskHandler)

	// Create the HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a goroutine so it doesn't block
	go func() {
		log.Println("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown: listen for SIGINT and SIGTERM signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Wait for an interrupt signal
	<-c
	log.Println("Received shutdown signal. Shutting down gracefully...")

	// Create a context with a timeout for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Perform any cleanup here (like closing DB connections, etc.)
	log.Println("Server stopped gracefully.")
}
