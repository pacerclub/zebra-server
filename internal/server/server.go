package server

import (
	"fmt"
	"time"

	"github.com/ZigaoWang/zebra-server/internal/config"
	"github.com/ZigaoWang/zebra-server/internal/repository"
	"github.com/ZigaoWang/zebra-server/internal/repository/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router      *gin.Engine
	cfg         *config.Config
	sessionRepo repository.SessionRepository
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
	workLogRepo repository.WorkLogRepository
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	// Set gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router
	router := gin.Default()

	// Setup CORS
	server := &Server{
		router: router,
		cfg:    cfg,
	}
	server.SetupCORS()

	// Initialize repositories
	sessionRepo := postgres.NewSessionRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	userRepo := postgres.NewUserRepository(db)
	workLogRepo := postgres.NewWorkLogRepository(db)

	// Create server instance
	server.sessionRepo = sessionRepo
	server.projectRepo = projectRepo
	server.userRepo = userRepo
	server.workLogRepo = workLogRepo

	// Setup routes
	server.setupRoutes()

	return server
}

func (s *Server) SetupCORS() {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.router.Use(cors.New(corsConfig))
}

func (s *Server) setupRoutes() {
	// Debug output
	fmt.Println("Setting up routes...")

	// Public routes
	s.router.POST("/api/v1/users/register", s.handleRegister())
	s.router.POST("/api/v1/users/login", s.handleLogin())

	// Direct file and audio access routes (outside of API group)
	s.router.GET("/files/:id", s.handleGetFile())
	s.router.GET("/audio/:id", s.handleGetAudio())

	// Protected API v1 group
	v1 := s.router.Group("/api/v1")
	v1.Use(s.authMiddleware())
	{
		// Projects
		projects := v1.Group("/projects")
		{
			projects.GET("", s.handleGetProjects())
			projects.POST("", s.handleCreateProject())
			projects.GET("/:id", s.handleGetProject())
			projects.PUT("/:id", s.handleUpdateProject())
			projects.DELETE("/:id", s.handleDeleteProject())

			// Sessions for a project
			sessions := projects.Group("/:id/sessions")
			{
				sessions.GET("", s.handleGetSessions())
				sessions.POST("", s.handleCreateSession())
				sessions.PUT("/:sessionId", s.handleUpdateSession())
				sessions.DELETE("/:sessionId", s.handleDeleteSession())
			}
		}

		// Work logs
		logs := v1.Group("/logs")
		{
			logs.GET("", s.handleGetWorkLogs())
			logs.POST("", s.handleCreateWorkLog())
			logs.GET("/:id", s.handleGetWorkLog())
			logs.PUT("/:id", s.handleUpdateWorkLog())
			logs.DELETE("/:id", s.handleDeleteWorkLog())
		}
	}

	// Debug output - print all routes
	routes := s.router.Routes()
	fmt.Println("=== ALL REGISTERED ROUTES ===")
	for _, route := range routes {
		fmt.Printf("Method: %s, Path: %s\n", route.Method, route.Path)
	}
	fmt.Println("============================")
}

func (s *Server) Run() error {
	// Start the main server
	return s.router.Run(fmt.Sprintf(":%s", s.cfg.Server.Port))
}
