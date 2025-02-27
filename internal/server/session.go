package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	StartTime string      `json:"start_time" binding:"required"`
	Duration  int64       `json:"duration" binding:"required"`
	Records   []RecordDTO `json:"records"`
}

type RecordDTO struct {
	Text      string   `json:"text"`
	GitLink   string   `json:"git_link,omitempty"`
	Files     []File   `json:"files,omitempty"`
	AudioURL  string   `json:"audio_url,omitempty"`
	Timestamp any      `json:"timestamp"` // Can be string or number
}

type UpdateSessionRequest struct {
	Duration int64 `json:"duration" binding:"required"`
}

type File struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
	Type string `json:"type" binding:"required"`
	Size int64  `json:"size" binding:"required"`
}

func (s *Server) handleCreateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get project ID from URL
		projectIDStr := c.Param("id")
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
			return
		}

		// Parse request body
		var req domain.Session
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err)})
			return
		}

		// Debug output
		fmt.Printf("Creating session: %+v\n", &req)
		fmt.Printf("Records: %+v\n", req.Records)

		// Ensure start time is set
		if req.StartTime.IsZero() {
			req.StartTime = time.Now().UTC()
		}

		// Ensure end time is set
		if req.EndTime.IsZero() {
			req.EndTime = time.Now().UTC()
		}

		// Create the session
		err = s.sessionRepo.Create(c, projectID, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create session: %v", err)})
			return
		}

		c.JSON(http.StatusCreated, req)
	}
}

func (s *Server) handleGetSessions() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))

		// Verify project ownership
		project, err := s.projectRepo.GetByID(c, projectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		sessions, err := s.sessionRepo.GetByProjectID(c, projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}

func (s *Server) handleUpdateSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		sessionID, err := uuid.Parse(c.Param("sessionId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))

		// Verify project ownership
		project, err := s.projectRepo.GetByID(c, projectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		var req UpdateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		session, err := s.sessionRepo.GetByID(c, sessionID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}

		if session.ProjectID != projectID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		session.Duration = req.Duration
		session.UpdatedAt = time.Now()

		if err := s.sessionRepo.Update(c, session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
			return
		}

		c.JSON(http.StatusOK, session)
	}
}

func (s *Server) handleDeleteSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		sessionID, err := uuid.Parse(c.Param("sessionId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))

		// Verify project ownership
		project, err := s.projectRepo.GetByID(c, projectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		session, err := s.sessionRepo.GetByID(c, sessionID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}

		if session.ProjectID != projectID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		if err := s.sessionRepo.Delete(c, sessionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
