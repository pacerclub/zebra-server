package main

import (
	"log"
	"os"

	"github.com/ZigaoWang/zebra-server/internal/config"
	"github.com/ZigaoWang/zebra-server/internal/database"
	"github.com/ZigaoWang/zebra-server/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Printf("Database initialization error: %v\n", err)
		os.Exit(1)
	}

	// Create and start server
	srv := server.NewServer(cfg, db)
	if err := srv.Run(); err != nil {
		log.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
