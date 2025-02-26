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
	StartTime time.Time `json:"start_time" binding:"required"`
	Duration  int64     `json:"duration" binding:"required"`
	Records   []struct {
		Text      string    `json:"text"`
		GitLink   string    `json:"git_link,omitempty"`
		Files     []File    `json:"files,omitempty"`
		AudioURL  string    `json:"audio_url,omitempty"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"records"`
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

		var req CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.StartTime.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
			return
		}

		session := &domain.Session{
			ProjectID: projectID,
			StartTime: req.StartTime,
			Duration:  req.Duration,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create records
		for _, r := range req.Records {
			record := &domain.Record{
				Text:      r.Text,
				GitLink:   r.GitLink,
				AudioURL:  r.AudioURL,
				Timestamp: r.Timestamp,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Create files
			for _, f := range r.Files {
				file := &domain.File{
					Name:      f.Name,
					URL:       f.URL,
					Type:      f.Type,
					Size:      f.Size,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				record.Files = append(record.Files, *file)
			}

			session.Records = append(session.Records, *record)
		}

		// Log the session being created
		fmt.Printf("Creating session: %+v\n", session)
		fmt.Printf("Records: %+v\n", session.Records)

		if err := s.sessionRepo.Create(c, session); err != nil {
			fmt.Printf("Error saving session: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(http.StatusCreated, session)
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
