package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c *gin.Context) {
	var body struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}

	userID := c.GetString("user_id")

	p, err := h.service.Create(c, body.Name, body.Description, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	projects, err := h.service.List(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	project, err := h.service.GetByID(c, id, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	userID := c.GetString("user_id")

	err := h.service.Update(c, id, body.Name, body.Description, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.Delete(c, id, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
