package handler

import (
	"log"
	"net/http"
	"strconv"

	"news_service/internal/domain"
	"news_service/internal/tools"
	"news_service/internal/validation"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsHandler struct {
	service   domain.NewsService
	validator *validation.Validator
}

func NewNewsHandler(service domain.NewsService) *NewsHandler {
	return &NewsHandler{
		service:   service,
		validator: validation.NewValidator(),
	}
}

func (h *NewsHandler) RegisterRoutes(router *gin.Engine) {
	log.Println("Start Registering Service Routes")

	router.GET("/", h.ListNews)
	log.Println("GET /")

	router.GET("/news/create", h.ShowCreateForm)
	log.Println("GET /news/create")

	router.POST("/news", h.CreateNews)
	log.Println("POST /news")

	router.GET("/news/:id", h.GetNews)
	log.Println("GET /news/:id")

	router.GET("/news/:id/edit", h.ShowEditForm)
	log.Println("GET /news/:id/edit")

	router.POST("/news/:id", h.HandleNewsUpdate)
	log.Println("POST /news/:id (with _method=PUT)")

	router.POST("/news/:id/delete", h.DeleteNews)
	log.Println("POST /news/:id/delete")

	router.GET("/news/search", h.SearchNews)
	log.Println("GET /news/search")

	log.Println("All Routes Registered Successfully")
}

// HandleNewsUpdate handles both PUT and POST requests for updating news
func (h *NewsHandler) HandleNewsUpdate(c *gin.Context) {
	method := c.PostForm("_method")
	if method == "PUT" {
		h.UpdateNews(c)
		return
	}
	c.HTML(http.StatusMethodNotAllowed, "error.html", gin.H{
		"error": "Method not allowed",
	})
}

func (h *NewsHandler) ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	news, total, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to fetch news",
		})
		return
	}

	c.HTML(http.StatusOK, "news/list.html", gin.H{
		"News":  news,
		"Total": total,
		"Page":  page,
		"Limit": limit,
	})
}

func (h *NewsHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "news/create.html", nil)
}

func (h *NewsHandler) CreateNews(c *gin.Context) {
	var news domain.News
	news.Title = c.PostForm("title")
	news.Content = c.PostForm("content")

	// Validate input
	result := h.validator.ValidateNews(news.Title, news.Content)
	if !result.IsValid {
		c.HTML(http.StatusBadRequest, "news/create.html", gin.H{
			"error":  "Validation Error",
			"errors": result.Errors,
		})
		return
	}

	now := tools.GetCurrentTime()
	news.CreatedAt = now
	news.UpdatedAt = now

	if err := h.service.Create(c.Request.Context(), &news); err != nil {
		log.Printf("Error creating news: %v", err)
		c.HTML(http.StatusInternalServerError, "news/create.html", gin.H{
			"error": "Error creating news",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

// validateObjectID validates the news ID and returns the ObjectID if valid
func (h *NewsHandler) validateObjectID(c *gin.Context) (primitive.ObjectID, bool) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid news ID",
		})
		return primitive.ObjectID{}, false
	}
	return objectID, true
}

// renderError renders an error page with the given error message
func (h *NewsHandler) renderError(c *gin.Context, status int, message string) {
	c.HTML(status, "error.html", gin.H{
		"error": message,
	})
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	objectID, valid := h.validateObjectID(c)
	if !valid {
		return
	}

	news, err := h.service.GetByID(c.Request.Context(), objectID.Hex())
	if err != nil {
		h.renderError(c, http.StatusNotFound, "News not found")
		return
	}

	c.HTML(http.StatusOK, "news/view.html", gin.H{
		"News": news,
	})
}

func (h *NewsHandler) ShowEditForm(c *gin.Context) {
	objectID, valid := h.validateObjectID(c)
	if !valid {
		return
	}

	news, err := h.service.GetByID(c.Request.Context(), objectID.Hex())
	if err != nil {
		h.renderError(c, http.StatusNotFound, "News not found")
		return
	}

	c.HTML(http.StatusOK, "news/edit.html", gin.H{
		"News": news,
	})
}

func (h *NewsHandler) UpdateNews(c *gin.Context) {
	objectID, valid := h.validateObjectID(c)
	if !valid {
		return
	}

	existingNews, err := h.service.GetByID(c.Request.Context(), objectID.Hex())
	if err != nil {
		h.renderError(c, http.StatusNotFound, "News not found")
		return
	}

	existingNews.Title = c.PostForm("title")
	existingNews.Content = c.PostForm("content")
	existingNews.UpdatedAt = tools.GetCurrentTime()

	// Validate input
	result := h.validator.ValidateNews(existingNews.Title, existingNews.Content)
	if !result.IsValid {
		c.HTML(http.StatusBadRequest, "news/edit.html", gin.H{
			"error":  "Validation Error",
			"errors": result.Errors,
			"News":   existingNews,
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), existingNews); err != nil {
		log.Printf("Error updating news: %v", err)
		c.HTML(http.StatusInternalServerError, "news/edit.html", gin.H{
			"error": "Error updating news",
			"News":  existingNews,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *NewsHandler) DeleteNews(c *gin.Context) {
	objectID, valid := h.validateObjectID(c)
	if !valid {
		return
	}

	if err := h.service.Delete(c.Request.Context(), objectID.Hex()); err != nil {
		log.Printf("Error deleting news: %v", err)
		h.renderError(c, http.StatusInternalServerError, "Failed to delete news")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *NewsHandler) SearchNews(c *gin.Context) {
	query := c.Query("q")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Validate search query before calling service
	result := h.validator.ValidateSearchQuery(query)
	if !result.IsValid {
		// Get current news list to show the page with error
		news, total, _ := h.service.GetAll(c.Request.Context(), 1, 10)
		c.HTML(http.StatusBadRequest, "news/list.html", gin.H{
			"error":  "Search query cannot be empty",
			"errors": result.Errors,
			"News":   news,
			"Total":  total,
			"Page":   1,
			"Limit":  10,
		})
		return
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	news, total, err := h.service.SearchNews(c.Request.Context(), query, pageNum, limitNum)
	if err != nil {
		log.Printf("Error searching news: %v", err)

		// Check if it's a validation error
		if domain.IsInvalidInput(err) {
			// Get current news list to show the page with error
			news, total, _ := h.service.GetAll(c.Request.Context(), 1, 10)
			c.HTML(http.StatusBadRequest, "news/list.html", gin.H{
				"error": "Invalid search query",
				"News":  news,
				"Total": total,
				"Page":  1,
				"Limit": 10,
			})
			return
		}

		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to search news",
		})
		return
	}

	c.HTML(http.StatusOK, "news/list.html", gin.H{
		"News":  news,
		"Total": total,
		"Page":  pageNum,
		"Limit": limitNum,
		"Query": query,
	})
}
