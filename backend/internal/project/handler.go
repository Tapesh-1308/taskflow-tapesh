package project

import (
	"net/http"

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
	h.log.Info("Create project request received")

	var body struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Warn("Create project validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}

	userID := c.GetString("user_id")

	p, err := h.service.Create(c, body.Name, body.Description, userID)
	if err != nil {
		h.log.Error("Create project failed", "error", err, "name", body.Name, "userID", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Project created successfully", "id", p.ID, "name", body.Name)
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	h.log.Info("List projects request", "userID", userID)

	projects, err := h.service.List(c, userID)
	if err != nil {
		h.log.Error("List projects failed", "error", err, "userID", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Projects listed successfully", "count", len(projects), "userID", userID)
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	h.log.Info("Get project by ID request", "id", id, "userID", userID)

	project, err := h.service.GetByID(c, id, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			h.log.Warn("Get project forbidden", "id", id, "userID", userID)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		h.log.Error("Get project failed", "error", err, "id", id, "userID", userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	h.log.Info("Project retrieved successfully", "id", id, "taskCount", len(project.Tasks))
	c.JSON(http.StatusOK, project)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	h.log.Info("Update project request", "id", id)

	var body struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Warn("Update project validation failed", "error", err, "id", id)
		c.JSON(400, gin.H{"error": "invalid"})
		return
	}

	userID := c.GetString("user_id")

	err := h.service.Update(c, id, body.Name, body.Description, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			h.log.Warn("Update project forbidden", "id", id, "userID", userID)
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		h.log.Error("Update project failed", "error", err, "id", id, "userID", userID)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Project updated successfully", "id", id)
	c.JSON(200, gin.H{"message": "updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	h.log.Info("Delete project request", "id", id, "userID", userID)

	err := h.service.Delete(c, id, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			h.log.Warn("Delete project forbidden", "id", id, "userID", userID)
			c.JSON(403, gin.H{"error": "forbidden"})
			return
		}
		h.log.Error("Delete project failed", "error", err, "id", id, "userID", userID)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Project deleted successfully", "id", id)
	c.JSON(200, gin.H{"message": "deleted"})
}
