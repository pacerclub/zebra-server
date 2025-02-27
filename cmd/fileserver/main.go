package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ZigaoWang/zebra-server/internal/config"
	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/ZigaoWang/zebra-server/internal/repository/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Initialize repository
	sessionRepo := postgres.NewSessionRepository(db)

	// Create router
	router := gin.Default()

	// Setup CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Handle OPTIONS requests for files and audio endpoints
	router.OPTIONS("/files/:id", func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}
		c.Status(http.StatusNoContent)
	})

	router.OPTIONS("/audio/:id", func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}
		c.Status(http.StatusNoContent)
	})

	// Files endpoint
	router.GET("/files/:id", func(c *gin.Context) {
		fileIDStr := c.Param("id")
		log.Printf("Received file request for ID: %s", fileIDStr)
		
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			log.Printf("Error parsing file ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID format"})
			return
		}

		file, err := sessionRepo.GetFileByID(c, fileID)
		if err != nil {
			log.Printf("Error retrieving file: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("File not found: %v", err)})
			return
		}

		if file == nil || file.Data == nil || len(file.Data) == 0 {
			log.Printf("File data is empty for ID: %s", fileID.String())
			c.JSON(http.StatusNotFound, gin.H{"error": "File data not found"})
			return
		}

		log.Printf("Successfully retrieved file: %s, size: %d bytes", file.Name, len(file.Data))

		// Set content type and other headers
		contentType := file.Type
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Set CORS headers - use specific origin instead of wildcard
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Content-Type", contentType)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.Name))
		c.Data(http.StatusOK, contentType, file.Data)
	})
	
	// Audio endpoint
	router.GET("/audio/:id", func(c *gin.Context) {
		recordIDStr := c.Param("id")
		log.Printf("Received audio request for ID: %s", recordIDStr)
		
		recordID, err := uuid.Parse(recordIDStr)
		if err != nil {
			log.Printf("Error parsing record ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID format"})
			return
		}

		record, err := sessionRepo.GetRecordByID(c, recordID)
		if err != nil {
			log.Printf("Error retrieving record: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Record not found: %v", err)})
			return
		}

		if record == nil || record.AudioData == nil || len(record.AudioData) == 0 {
			log.Printf("Audio data is empty for record ID: %s", recordID.String())
			c.JSON(http.StatusNotFound, gin.H{"error": "Audio data not found"})
			return
		}

		log.Printf("Successfully retrieved audio for record: %s, size: %d bytes", recordID.String(), len(record.AudioData))

		// Set content type and other headers
		// Set CORS headers - use specific origin instead of wildcard
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		c.Header("Content-Type", "audio/mpeg")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=audio_%s.mp3", recordID))
		c.Data(http.StatusOK, "audio/mpeg", record.AudioData)
	})

	// Debug endpoint to list all routes
	router.GET("/debug/routes", func(c *gin.Context) {
		routes := router.Routes()
		var routeList []string
		for _, route := range routes {
			routeList = append(routeList, fmt.Sprintf("%s %s", route.Method, route.Path))
		}
		c.JSON(http.StatusOK, gin.H{"routes": routeList})
	})

	// Start the server
	port := "8080"
	log.Printf("File server starting on port %s...", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start file server: %v", err)
	}
}
