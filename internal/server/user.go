package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// handleGetProfile handles requests to get the user's profile
func (s *Server) handleGetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := uuid.Parse(c.GetString("user_id"))
		
		user, err := s.userRepo.GetByID(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
			return
		}
		
		// Don't return the password hash
		user.Password = ""
		
		c.JSON(http.StatusOK, user)
	}
}

// handleUpdateProfile handles requests to update the user's profile
func (s *Server) handleUpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation will be added later
		c.JSON(http.StatusOK, gin.H{"message": "Update profile endpoint"})
	}
}
