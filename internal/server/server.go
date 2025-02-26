package server

import (
	"fmt"

	"github.com/ZigaoWang/zebra-server/internal/config"
	"github.com/ZigaoWang/zebra-server/internal/middleware"
	"github.com/ZigaoWang/zebra-server/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	config      *config.Config
	router      *gin.Engine
	userRepo    *postgres.UserRepository
	workLogRepo *postgres.WorkLogRepository
	projectRepo *postgres.ProjectRepository
	sessionRepo *postgres.SessionRepository
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	server := &Server{
		config:      cfg,
		router:      gin.Default(),
		userRepo:    postgres.NewUserRepository(db),
		workLogRepo: postgres.NewWorkLogRepository(db),
		projectRepo: postgres.NewProjectRepository(db),
		sessionRepo: postgres.NewSessionRepository(db),
	}

	// Set gin mode
	gin.SetMode(cfg.Server.Mode)

	// Setup middleware
	server.router.Use(middleware.CORS())

	// Setup routes
	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// Public routes
	s.router.POST("/api/v1/auth/register", s.handleRegister())
	s.router.POST("/api/v1/auth/login", s.handleLogin())

	// Protected routes
	authorized := s.router.Group("/api/v1")
	authorized.Use(s.authMiddleware())
	{
		// Work logs
		authorized.POST("/logs", s.handleCreateWorkLog())
		authorized.GET("/logs", s.handleGetWorkLogs())
		authorized.GET("/logs/:id", s.handleGetWorkLog())
		authorized.PUT("/logs/:id", s.handleUpdateWorkLog())
		authorized.DELETE("/logs/:id", s.handleDeleteWorkLog())

		// Log entries
		authorized.POST("/logs/:id/entries", s.handleCreateLogEntry())
		authorized.GET("/logs/:id/entries", s.handleGetLogEntries())
		authorized.PUT("/logs/:id/entries/:entryId", s.handleUpdateLogEntry())
		authorized.DELETE("/logs/:id/entries/:entryId", s.handleDeleteLogEntry())

		// Projects
		authorized.POST("/projects", s.handleCreateProject())
		authorized.GET("/projects", s.handleGetProjects())
		authorized.GET("/projects/:id", s.handleGetProject())
		authorized.PUT("/projects/:id", s.handleUpdateProject())
		authorized.DELETE("/projects/:id", s.handleDeleteProject())

		// Sessions
		authorized.POST("/projects/:id/sessions", s.handleCreateSession())
		authorized.GET("/projects/:id/sessions", s.handleGetSessions())
		authorized.PUT("/projects/:id/sessions/:sessionId", s.handleUpdateSession())
		authorized.DELETE("/projects/:id/sessions/:sessionId", s.handleDeleteSession())

		// User
		authorized.GET("/user/profile", s.handleGetProfile())
		authorized.PUT("/user/profile", s.handleUpdateProfile())
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.Server.Port))
}

// Handler function declarations are in their respective files:
// - auth.go: handleRegister, handleLogin, authMiddleware
// - worklog.go: handleCreateWorkLog, handleGetWorkLogs, etc.
// - project.go: handleCreateProject, handleGetProjects, etc.

func (s *Server) handleGetWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleUpdateWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleDeleteWorkLog() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleCreateLogEntry() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleGetLogEntries() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleUpdateLogEntry() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleDeleteLogEntry() gin.HandlerFunc {
	return func(c *gin.Context) {}
}



func (s *Server) handleGetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Server) handleUpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
