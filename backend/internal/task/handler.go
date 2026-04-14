package task

import (
	"fmt"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	log     *slog.Logger
}

func NewHandler(s Service, log *slog.Logger) *Handler {
	return &Handler{service: s, log: log}
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

func (h *Handler) ListByProject(c *gin.Context) {
	projectID := c.Param("id")

	status := c.Query("status")
	assignee := c.Query("assignee")

	var statusPtr *string
	var assigneePtr *string

	if status != "" {
		statusPtr = &status
	}
	if assignee != "" {
		assigneePtr = &assignee
	}

	userID := c.GetString("user_id")

	tasks, err := h.service.ListByProject(c, projectID, statusPtr, assigneePtr, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"tasks": tasks})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var body UpdateTaskInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	userID := c.GetString("user_id")

	fmt.Println(body)
	err := h.service.Update(c, id, body, userID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	userID := c.GetString("user_id")

	err := h.service.Delete(c, id, userID)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "deleted"})
}
