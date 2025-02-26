package server

import (
	"net/http"
	"time"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) handleCreateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateProjectRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))
		project := &domain.Project{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.projectRepo.Create(c, project); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
			return
		}

		c.JSON(http.StatusCreated, project)
	}
}

func (s *Server) handleGetProjects() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := uuid.Parse(c.GetString("user_id"))

		projects, err := s.projectRepo.GetByUserID(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
			return
		}

		c.JSON(http.StatusOK, projects)
	}
}

func (s *Server) handleGetProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		userIDStr := c.GetString("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user authentication"})
			return
		}

		project, err := s.projectRepo.GetByID(c, projectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Project already includes sessions due to preloading
		c.JSON(http.StatusOK, project)
	}
}

func (s *Server) handleUpdateProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		var req UpdateProjectRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))

		project, err := s.projectRepo.GetByID(c, projectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		project.Name = req.Name
		project.Description = req.Description
		project.UpdatedAt = time.Now()

		if err := s.projectRepo.Update(c, project); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
			return
		}

		c.JSON(http.StatusOK, project)
	}
}

func (s *Server) handleDeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}

		userID, _ := uuid.Parse(c.GetString("user_id"))

		project, err := s.projectRepo.GetByID(c, projectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}

		if project.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		if err := s.projectRepo.Delete(c, projectID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
