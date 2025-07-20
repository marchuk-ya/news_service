package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"news_service/internal/config"
	"news_service/internal/handler"
	"news_service/internal/logger"
	"news_service/internal/middleware"
	"news_service/internal/observability"
	"news_service/internal/repository/mongodb"
	"news_service/internal/service"
)

// @title           News Service API
// @version         1.0
// @description     A RESTful API for managing news articles
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API key for authentication

// Server represents the HTTP server with graceful shutdown capabilities
type Server struct {
	httpServer  *http.Server
	logger      *logger.Logger
	tracer      *observability.Tracer
	mongoClient *mongo.Client
	shutdownWg  sync.WaitGroup
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, logger *logger.Logger, tracer *observability.Tracer, mongoClient *mongo.Client, router *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
		logger:      logger,
		tracer:      tracer,
		mongoClient: mongoClient,
	}
}

// Start starts the server in a goroutine
func (s *Server) Start() {
	s.shutdownWg.Add(1)
	go func() {
		defer s.shutdownWg.Done()
		s.logger.Info("Starting HTTP server", "port", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start server", "error", err)
			log.Fatal(err)
		}
	}()
}

// Shutdown gracefully shuts down the server and all resources
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Starting graceful shutdown...")

	// Create a channel to signal when shutdown is complete
	shutdownComplete := make(chan struct{})

	// Start shutdown process in a goroutine
	go func() {
		defer close(shutdownComplete)

		// Shutdown HTTP server
		s.logger.Info("Shutting down HTTP server...")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP server shutdown error", "error", err)
		} else {
			s.logger.Info("HTTP server shutdown completed")
		}

		// Shutdown tracer
		s.logger.Info("Shutting down tracer...")
		if err := s.tracer.Shutdown(ctx); err != nil {
			s.logger.Error("Tracer shutdown error", "error", err)
		} else {
			s.logger.Info("Tracer shutdown completed")
		}

		// Shutdown MongoDB connection
		s.logger.Info("Shutting down MongoDB connection...")
		if err := s.mongoClient.Disconnect(ctx); err != nil {
			s.logger.Error("MongoDB disconnect error", "error", err)
		} else {
			s.logger.Info("MongoDB disconnect completed")
		}
	}()

	// Wait for shutdown to complete or context to timeout
	select {
	case <-shutdownComplete:
		s.logger.Info("Graceful shutdown completed successfully")
		return nil
	case <-ctx.Done():
		s.logger.Error("Shutdown timeout exceeded", "error", ctx.Err())
		return ctx.Err()
	}
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.App.LogLevel)
	logger.Info("Starting news service", "environment", cfg.App.Environment, "version", "1.0.0")

	// Initialize tracer
	tracer, err := observability.NewTracer("news-service", "1.0.0")
	if err != nil {
		logger.Error("Failed to initialize tracer", "error", err)
		log.Fatal(err)
	}

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to MongoDB with connection pooling
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.Timeout)
	defer cancel()

	// Configure MongoDB client with connection pooling
	clientOptions := options.Client().ApplyURI(cfg.Database.URI)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetMaxConnIdleTime(30 * time.Minute)
	clientOptions.SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		log.Fatal(err)
	}

	// Test database connection
	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Failed to ping MongoDB", "error", err)
		log.Fatal(err)
	}

	logger.Info("Successfully connected to MongoDB", "uri", cfg.Database.URI, "database", cfg.Database.Database)

	// Initialize dependencies
	newsRepo := mongodb.NewNewsRepository(client, cfg.Database.Database)
	newsUseCase := service.NewNewsUseCase(newsRepo)
	newsHandler := handler.NewNewsHandler(newsUseCase)
	newsAPIHandler := handler.NewNewsAPIHandler(newsUseCase)
	healthHandler := handler.NewHealthHandler(newsRepo)

	// Setup Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorHandlingMiddleware(logger))
	router.Use(middleware.SecurityMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.TimeoutMiddleware(30 * time.Second))

	// Setup template functions
	funcMap := template.FuncMap{
		"subtract": func(a, b int) int { return a - b },
		"add":      func(a, b int) int { return a + b },
		"multiply": func(a, b int) int { return a * b },
	}
	router.SetFuncMap(funcMap)
	router.LoadHTMLGlob("web/templates/**/*")

	// Serve static files
	router.Static("/static", "./web/static")

	// Register routes
	newsHandler.RegisterRoutes(router)
	newsAPIHandler.RegisterAPIRoutes(router)
	healthHandler.RegisterRoutes(router)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API documentation redirect
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Create and start server
	server := NewServer(cfg, logger, tracer, client, router)
	server.Start()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Received shutdown signal, starting graceful shutdown...")

	// Create context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Server shutdown completed successfully")
}
