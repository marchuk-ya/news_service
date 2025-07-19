package handler

import (
	"net/http"

	"news_service/internal/domain"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	newsRepo domain.NewsRepository
}

func NewHealthHandler(newsRepo domain.NewsRepository) *HealthHandler {
	return &HealthHandler{
		newsRepo: newsRepo,
	}
}

func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.HealthCheck)
	router.GET("/health/readiness", h.ReadinessCheck)
	router.GET("/api/v1/health", h.HealthCheckAPI)
	router.GET("/api/v1/health/readiness", h.ReadinessCheckAPI)
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"news_service"`
	Version string `json:"version" example:"1.0.0"`
}

// ReadinessResponse represents the readiness check response
type ReadinessResponse struct {
	Status   string `json:"status" example:"ready"`
	Service  string `json:"service" example:"news_service"`
	Database string `json:"database" example:"connected"`
	Version  string `json:"version" example:"1.0.0"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Basic health check endpoint
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "news_service",
	})
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Readiness check including database connectivity
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} ReadinessResponse
// @Failure 503 {object} ErrorResponse
// @Router /health/readiness [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check database connectivity by performing a simple operation
	_, _, err := h.newsRepo.GetAll(c.Request.Context(), 1, 1)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": "news_service",
			"error":   "database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ready",
		"service":  "news_service",
		"database": "connected",
	})
}

// HealthCheckAPI godoc
// @Summary Health check API
// @Description Basic health check endpoint for API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/v1/health [get]
func (h *HealthHandler) HealthCheckAPI(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Service: "news_service",
		Version: "1.0.0",
	})
}

// ReadinessCheckAPI godoc
// @Summary Readiness check API
// @Description Readiness check including database connectivity for API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} ReadinessResponse
// @Failure 503 {object} ErrorResponse
// @Router /api/v1/health/readiness [get]
func (h *HealthHandler) ReadinessCheckAPI(c *gin.Context) {
	// Check database connectivity by performing a simple operation
	_, _, err := h.newsRepo.GetAll(c.Request.Context(), 1, 1)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "Service unavailable",
			Code:    "SERVICE_UNAVAILABLE",
			Message: "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, ReadinessResponse{
		Status:   "ready",
		Service:  "news_service",
		Database: "connected",
		Version:  "1.0.0",
	})
}
