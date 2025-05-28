package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"news_service/internal/domain"
	"news_service/internal/tools"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NewsHandler struct {
	service domain.NewsRepository
}

func NewNewsHandler(service domain.NewsRepository) *NewsHandler {
	return &NewsHandler{
		service: service,
	}
}

// validateNews performs simple validation of news fields
func (h *NewsHandler) validateNews(news *domain.News) []string {
	var errors []string

	title := strings.TrimSpace(news.Title)
	if title == "" {
		errors = append(errors, "Title cannot be empty")
	} else if len(title) < 3 {
		errors = append(errors, "Title must contain at least 3 characters")
	} else if len(title) > 200 {
		errors = append(errors, "Title cannot exceed 200 characters")
	}

	content := strings.TrimSpace(news.Content)
	if content == "" {
		errors = append(errors, "Content cannot be empty")
	} else if len(content) < 10 {
		errors = append(errors, "Content must contain at least 10 characters")
	}

	return errors
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

	now := tools.GetCurrentTime()
	news.CreatedAt = now
	news.UpdatedAt = now

	if errors := h.validateNews(&news); len(errors) > 0 {
		c.HTML(http.StatusBadRequest, "news/create.html", gin.H{
			"error":  "Validation Error",
			"errors": errors,
		})
		return
	}

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

	if errors := h.validateNews(existingNews); len(errors) > 0 {
		c.HTML(http.StatusBadRequest, "news/edit.html", gin.H{
			"error":  "Validation Error",
			"errors": errors,
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
