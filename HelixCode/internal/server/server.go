package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"dev.helix.code/internal/config"
	"dev.helix.code/internal/database"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	db     *database.Database
	server *http.Server
	router *gin.Engine
}

// New creates a new HTTP server
func New(cfg *config.Config, db *database.Database) *Server {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(SecurityMiddleware())

	server := &Server{
		config: cfg,
		db:     db,
		router: router,
	}

	// Setup routes
	server.setupRoutes()

	// Create HTTP server
	server.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("🚀 Starting HelixCode server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// setupRoutes sets up all HTTP routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.notImplemented)
			auth.POST("/login", s.notImplemented)
			auth.POST("/logout", s.notImplemented)
			auth.POST("/refresh", s.notImplemented)
		}

		// User routes
		users := api.Group("/users")
		users.Use(s.authMiddleware())
		{
			users.GET("/me", s.notImplemented)
			users.PUT("/me", s.notImplemented)
			users.DELETE("/me", s.notImplemented)
		}

		// Worker routes
		workers := api.Group("/workers")
		workers.Use(s.authMiddleware())
		{
			workers.GET("", s.notImplemented)
			workers.POST("", s.notImplemented)
			workers.GET("/:id", s.notImplemented)
			workers.PUT("/:id", s.notImplemented)
			workers.DELETE("/:id", s.notImplemented)
			workers.POST("/:id/heartbeat", s.notImplemented)
			workers.GET("/:id/metrics", s.notImplemented)
		}

		// Task routes
		tasks := api.Group("/tasks")
		tasks.Use(s.authMiddleware())
		{
			tasks.GET("", s.notImplemented)
			tasks.POST("", s.notImplemented)
			tasks.GET("/:id", s.notImplemented)
			tasks.PUT("/:id", s.notImplemented)
			tasks.DELETE("/:id", s.notImplemented)
			tasks.POST("/:id/assign", s.notImplemented)
			tasks.POST("/:id/start", s.notImplemented)
			tasks.POST("/:id/complete", s.notImplemented)
			tasks.POST("/:id/fail", s.notImplemented)
			tasks.POST("/:id/checkpoint", s.notImplemented)
			tasks.GET("/:id/checkpoints", s.notImplemented)
			tasks.POST("/:id/retry", s.notImplemented)
		}

		// Project routes
		projects := api.Group("/projects")
		projects.Use(s.authMiddleware())
		{
			projects.GET("", s.notImplemented)
			projects.POST("", s.notImplemented)
			projects.GET("/:id", s.notImplemented)
			projects.PUT("/:id", s.notImplemented)
			projects.DELETE("/:id", s.notImplemented)
			projects.GET("/:id/sessions", s.notImplemented)
		}

		// Session routes
		sessions := api.Group("/sessions")
		sessions.Use(s.authMiddleware())
		{
			sessions.GET("", s.notImplemented)
			sessions.POST("", s.notImplemented)
			sessions.GET("/:id", s.notImplemented)
			sessions.PUT("/:id", s.notImplemented)
			sessions.DELETE("/:id", s.notImplemented)
		}

		// System routes
		system := api.Group("/system")
		system.Use(s.authMiddleware())
		{
			system.GET("/stats", s.notImplemented)
			system.GET("/status", s.notImplemented)
		}
	}

	// WebSocket routes
	s.router.GET("/ws", s.notImplemented)

	// Static file serving for web interface
	s.router.Static("/static", "./web/frontend/static")
	s.router.StaticFile("/", "./web/frontend/index.html")
	s.router.StaticFile("/favicon.ico", "./assets/icons/icon-32x32.png")
}

// Handler methods

func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	if err := s.db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC(),
	})
}

func (s *Server) notImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"status":  "error",
		"message": "Not implemented yet",
	})
}

// Middleware

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement authentication middleware
		// For now, just continue
		c.Next()
	}
}

// CORSMiddleware provides CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecurityMiddleware provides security headers
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}