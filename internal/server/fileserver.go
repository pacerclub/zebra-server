package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterFileRoutes registers the file and audio routes
func (s *Server) RegisterFileRoutes(router *gin.Engine) {
	fmt.Println("Registering file and audio routes")
	
	// Files endpoint
	router.GET("/files/:id", func(c *gin.Context) {
		fileIDStr := c.Param("id")
		fmt.Printf("Received file request for ID: %s\n", fileIDStr)
		
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			fmt.Printf("Error parsing file ID: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID format"})
			return
		}

		// Check if file ID is zero UUID
		if fileID == uuid.Nil {
			fmt.Printf("File ID is nil UUID\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID (zero UUID)"})
			return
		}

		fmt.Printf("Looking up file with ID: %s\n", fileID.String())
		file, err := s.sessionRepo.GetFileByID(c, fileID)
		if err != nil {
			fmt.Printf("Error retrieving file: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("File not found: %v", err)})
			return
		}

		if file == nil {
			fmt.Printf("File is nil for ID: %s\n", fileID.String())
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		if file.Data == nil || len(file.Data) == 0 {
			fmt.Printf("File data is empty for ID: %s\n", fileID.String())
			c.JSON(http.StatusNotFound, gin.H{"error": "File data not found"})
			return
		}

		fmt.Printf("Successfully retrieved file: %s, size: %d bytes\n", file.Name, len(file.Data))

		// Set content type and other headers
		contentType := file.Type
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Header("Content-Type", contentType)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", file.Name))
		c.Header("Cache-Control", "public, max-age=31536000")

		// Serve the file data
		c.Data(http.StatusOK, contentType, file.Data)
	})
	
	// Audio endpoint
	router.GET("/audio/:id", func(c *gin.Context) {
		recordIDStr := c.Param("id")
		fmt.Printf("Received audio request for ID: %s\n", recordIDStr)
		
		recordID, err := uuid.Parse(recordIDStr)
		if err != nil {
			fmt.Printf("Error parsing record ID: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID format"})
			return
		}

		// Check if record ID is zero UUID
		if recordID == uuid.Nil {
			fmt.Printf("Record ID is nil UUID\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID (zero UUID)"})
			return
		}

		fmt.Printf("Looking up record with ID: %s\n", recordID.String())
		record, err := s.sessionRepo.GetRecordByID(c, recordID)
		if err != nil {
			fmt.Printf("Error retrieving record: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Record not found: %v", err)})
			return
		}

		if record == nil {
			fmt.Printf("Record is nil for ID: %s\n", recordID.String())
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		if record.AudioData == nil || len(record.AudioData) == 0 {
			fmt.Printf("Audio data is empty for record ID: %s\n", recordID.String())
			c.JSON(http.StatusNotFound, gin.H{"error": "Audio data not found"})
			return
		}

		fmt.Printf("Successfully retrieved audio for record: %s, size: %d bytes\n", recordID.String(), len(record.AudioData))

		// Set content type and other headers
		c.Header("Content-Type", "audio/mpeg")
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=audio_%s.mp3", recordID))
		c.Header("Cache-Control", "public, max-age=31536000")

		// Serve the audio data
		reader := bytes.NewReader(record.AudioData)
		c.DataFromReader(http.StatusOK, int64(len(record.AudioData)), "audio/mpeg", reader, nil)
	})
	
	fmt.Println("File and audio routes registered")
}
