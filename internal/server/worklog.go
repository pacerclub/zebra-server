package server

import (
	"net/http"
	"time"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateWorkLogRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (s *Server) handleCreateWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateWorkLogRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))
		workLog := &domain.WorkLog{
			ID:          uuid.New(),
			UserID:      userID,
			Title:       req.Title,
			Description: req.Description,
			StartTime:   time.Now(),
			Status:      "active",
			Tags:        req.Tags,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Save to database
		if err := s.workLogRepo.Create(c, workLog); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work log"})
			return
		}

		c.JSON(http.StatusCreated, workLog)
	}
}

func (s *Server) handleGetWorkLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := uuid.Parse(c.GetString("user_id"))
		workLogs, err := s.workLogRepo.GetByUserID(c, userID, 10, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, workLogs)
	}
}

func (s *Server) handleGetWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Get work log endpoint"})
	}
}

func (s *Server) handleUpdateWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Update work log endpoint"})
	}
}

func (s *Server) handleDeleteWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Delete work log endpoint"})
	}
}

func (s *Server) handleCreateLogEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Create log entry endpoint"})
	}
}

func (s *Server) handleGetLogEntries() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Get log entries endpoint"})
	}
}

func (s *Server) handleUpdateLogEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Update log entry endpoint"})
	}
}

func (s *Server) handleDeleteLogEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Delete log entry endpoint"})
	}
}
