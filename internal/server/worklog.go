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

		// Get work logs from database
		workLogs, err := s.workLogRepo.GetByUserID(c, userID, 50, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch work logs"})
			return
		}

		c.JSON(http.StatusOK, workLogs)
	}
}
