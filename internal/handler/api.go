package handler

import (
	"net/http"
	"strconv"

	"news_service/internal/domain"

	"github.com/gin-gonic/gin"
)

// NewsAPIHandler handles REST API requests for news
type NewsAPIHandler struct {
	service domain.NewsService
}

// NewNewsAPIHandler creates a new API handler
func NewNewsAPIHandler(service domain.NewsService) *NewsAPIHandler {
	return &NewsAPIHandler{
		service: service,
	}
}

// RegisterAPIRoutes registers REST API routes
func (h *NewsAPIHandler) RegisterAPIRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/news", h.ListNews)
		api.POST("/news", h.CreateNews)
		api.GET("/news/:id", h.GetNews)
		api.PUT("/news/:id", h.UpdateNews)
		api.DELETE("/news/:id", h.DeleteNews)
		api.GET("/news/search", h.SearchNews)
	}
}

// CreateNewsRequest represents the request body for creating news
type CreateNewsRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200" example:"Breaking News"`
	Content string `json:"content" binding:"required,min=10" example:"This is a breaking news article content."`
}

// UpdateNewsRequest represents the request body for updating news
type UpdateNewsRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=200" example:"Updated Breaking News"`
	Content string `json:"content" binding:"required,min=10" example:"This is an updated breaking news article content."`
}

// NewsResponse represents the response for news operations
type NewsResponse struct {
	ID        string `json:"id" example:"507f1f77bcf86cd799439011"`
	Title     string `json:"title" example:"Breaking News"`
	Content   string `json:"content" example:"This is a breaking news article content."`
	CreatedAt string `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

// ListNewsResponse represents the response for listing news
type ListNewsResponse struct {
	News       []*NewsResponse `json:"news"`
	Total      int64           `json:"total" example:"100"`
	Page       int             `json:"page" example:"1"`
	Limit      int             `json:"limit" example:"10"`
	TotalPages int             `json:"total_pages" example:"10"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"News not found"`
	Code    string `json:"code" example:"NOT_FOUND"`
	Message string `json:"message" example:"The requested news article was not found"`
}

// ListNews godoc
// @Summary List all news articles
// @Description Get a paginated list of all news articles
// @Tags news
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Number of items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Success 200 {object} ListNewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news [get]
func (h *NewsAPIHandler) ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	news, total, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch news",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Convert domain news to response format
	newsResponse := make([]*NewsResponse, len(news))
	for i, n := range news {
		newsResponse[i] = &NewsResponse{
			ID:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: n.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, ListNewsResponse{
		News:       newsResponse,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// CreateNews godoc
// @Summary Create a new news article
// @Description Create a new news article with title and content
// @Tags news
// @Accept json
// @Produce json
// @Param news body CreateNewsRequest true "News article data"
// @Success 201 {object} NewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news [post]
func (h *NewsAPIHandler) CreateNews(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	news := &domain.News{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.service.Create(c.Request.Context(), news); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create news",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	response := &NewsResponse{
		ID:        news.ID,
		Title:     news.Title,
		Content:   news.Content,
		CreatedAt: news.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: news.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetNews godoc
// @Summary Get a news article by ID
// @Description Get a specific news article by its ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID" example("507f1f77bcf86cd799439011")
// @Success 200 {object} NewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news/{id} [get]
func (h *NewsAPIHandler) GetNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid news ID",
			Code:    "VALIDATION_ERROR",
			Message: "News ID cannot be empty",
		})
		return
	}

	news, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if err == domain.ErrNotFound {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		} else if err == domain.ErrInvalidInput {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to get news",
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	response := &NewsResponse{
		ID:        news.ID,
		Title:     news.Title,
		Content:   news.Content,
		CreatedAt: news.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: news.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateNews godoc
// @Summary Update a news article
// @Description Update an existing news article by ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID" example("507f1f77bcf86cd799439011")
// @Param news body UpdateNewsRequest true "Updated news article data"
// @Success 200 {object} NewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news/{id} [put]
func (h *NewsAPIHandler) UpdateNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid news ID",
			Code:    "VALIDATION_ERROR",
			Message: "News ID cannot be empty",
		})
		return
	}

	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Get existing news to preserve creation time
	existingNews, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if err == domain.ErrNotFound {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to get existing news",
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	// Update fields
	existingNews.Title = req.Title
	existingNews.Content = req.Content

	if err := h.service.Update(c.Request.Context(), existingNews); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update news",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	response := &NewsResponse{
		ID:        existingNews.ID,
		Title:     existingNews.Title,
		Content:   existingNews.Content,
		CreatedAt: existingNews.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: existingNews.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteNews godoc
// @Summary Delete a news article
// @Description Delete a news article by ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "News ID" example("507f1f77bcf86cd799439011")
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news/{id} [delete]
func (h *NewsAPIHandler) DeleteNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid news ID",
			Code:    "VALIDATION_ERROR",
			Message: "News ID cannot be empty",
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"

		if err == domain.ErrNotFound {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		} else if err == domain.ErrInvalidInput {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}

		c.JSON(status, ErrorResponse{
			Error:   "Failed to delete news",
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchNewsResponse represents the response for searching news
type SearchNewsResponse struct {
	News       []*NewsResponse `json:"news"`
	Total      int64           `json:"total" example:"5"`
	Query      string          `json:"query" example:"breaking"`
	Page       int             `json:"page" example:"1"`
	Limit      int             `json:"limit" example:"10"`
	TotalPages int             `json:"total_pages" example:"1"`
}

// SearchNews godoc
// @Summary Search news articles
// @Description Search news articles by title and content
// @Tags news
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Number of items per page (default: 10, max: 100)" minimum(1) maximum(100)
// @Success 200 {object} SearchNewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news/search [get]
func (h *NewsAPIHandler) SearchNews(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Search query is required",
			Code:    "VALIDATION_ERROR",
			Message: "Query parameter 'q' cannot be empty",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	news, total, err := h.service.SearchNews(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to search news",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Convert domain news to response format
	newsResponse := make([]*NewsResponse, len(news))
	for i, n := range news {
		newsResponse[i] = &NewsResponse{
			ID:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: n.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, SearchNewsResponse{
		News:       newsResponse,
		Total:      total,
		Query:      query,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
