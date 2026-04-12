package task

import (
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
		Title     string `json:"title"`
		ProjectID string `json:"project_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	userID := c.GetString("user_id")

	task, err := h.service.Create(c, body.Title, body.ProjectID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, task)
}

func (h *Handler) Assign(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		UserID string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	userID := c.GetString("user_id")

	err := h.service.Assign(c, id, body.UserID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "assigned"})
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	userID := c.GetString("user_id")

	err := h.service.UpdateStatus(c, id, body.Status, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}

func (h *Handler) List(c *gin.Context) {
	projectID := c.Query("project_id")
	status := c.Query("status")

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	userID := c.GetString("user_id")

	tasks, err := h.service.List(c, projectID, statusPtr, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"tasks": tasks})
}
