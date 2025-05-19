package handler

import (
	"net/http"
	"strconv"

	"news_service/internal/domain"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	service domain.NewsService
}

func NewNewsHandler(service domain.NewsService) *NewsHandler {
	return &NewsHandler{
		service: service,
	}
}

func (h *NewsHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/", h.ListNews)
	router.GET("/news/create", h.ShowCreateForm)
	router.POST("/news", h.CreateNews)
	router.GET("/news/:id", h.GetNews)
	router.GET("/news/:id/edit", h.ShowEditForm)
	router.PUT("/news/:id", h.UpdateNews)
	router.DELETE("/news/:id", h.DeleteNews)
	router.GET("/news/search", h.SearchNews)
}

func (h *NewsHandler) ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	news, total, err := h.service.GetAllNews(page, limit)
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
	if err := c.ShouldBind(&news); err != nil {
		c.HTML(http.StatusBadRequest, "news/create.html", gin.H{
			"error": "Invalid input",
		})
		return
	}

	if err := h.service.CreateNews(&news); err != nil {
		c.HTML(http.StatusInternalServerError, "news/create.html", gin.H{
			"error": "Failed to create news",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	id := c.Param("id")
	news, err := h.service.GetNewsByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "News not found",
		})
		return
	}

	c.HTML(http.StatusOK, "news/view.html", gin.H{
		"News": news,
	})
}

func (h *NewsHandler) ShowEditForm(c *gin.Context) {
	id := c.Param("id")
	news, err := h.service.GetNewsByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "News not found",
		})
		return
	}

	c.HTML(http.StatusOK, "news/edit.html", gin.H{
		"News": news,
	})
}

func (h *NewsHandler) UpdateNews(c *gin.Context) {
	id := c.Param("id")
	var news domain.News
	if err := c.ShouldBind(&news); err != nil {
		c.HTML(http.StatusBadRequest, "news/edit.html", gin.H{
			"error": "Invalid input",
		})
		return
	}

	existingNews, err := h.service.GetNewsByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "News not found",
		})
		return
	}

	news.ID = existingNews.ID
	news.CreatedAt = existingNews.CreatedAt

	if err := h.service.UpdateNews(&news); err != nil {
		c.HTML(http.StatusInternalServerError, "news/edit.html", gin.H{
			"error": "Failed to update news",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/news/"+id)
}

func (h *NewsHandler) DeleteNews(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteNews(id); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to delete news",
		})
		return
	}

	c.HTML(http.StatusOK, "news/empty.html", nil)
}

func (h *NewsHandler) SearchNews(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	news, total, err := h.service.SearchNews(query, page, limit)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Failed to search news",
		})
		return
	}

	c.HTML(http.StatusOK, "news/list.html", gin.H{
		"News":  news,
		"Total": total,
		"Page":  page,
		"Limit": limit,
		"Query": query,
	})
}
