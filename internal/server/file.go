package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// handleGetFile handles GET requests for file content
func (s *Server) handleGetFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("id")
		if fileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
			return
		}

		// Parse UUID
		id, err := uuid.Parse(fileID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID format"})
			return
		}

		// Get file from database
		file, err := s.sessionRepo.GetFileByID(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve file: %v", err)})
			return
		}

		if file == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		// Set content type based on file type
		contentType := file.Type
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
		c.Data(http.StatusOK, contentType, file.Data)
	}
}

// handleGetAudio handles GET requests for audio content
func (s *Server) handleGetAudio() gin.HandlerFunc {
	return func(c *gin.Context) {
		audioID := c.Param("id")
		if audioID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Audio ID is required"})
			return
		}

		// Parse UUID
		id, err := uuid.Parse(audioID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audio ID format"})
			return
		}

		// Get record from database
		record, err := s.sessionRepo.GetRecordByID(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve audio: %v", err)})
			return
		}

		if record == nil || record.AudioData == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Audio not found"})
			return
		}

		// Set content type for audio
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audio-%s.webm", audioID))
		c.Data(http.StatusOK, "audio/webm", record.AudioData)
	}
}
