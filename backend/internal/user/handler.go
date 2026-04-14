package user

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	log     *slog.Logger
}

func NewHandler(service Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) Register(c *gin.Context) {
	h.log.Info("Register request received")

	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.log.Warn("Register validation failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation failed",
		})
		return
	}

	token, err := h.service.Register(c, body.Name, body.Email, body.Password)
	if err != nil {
		h.log.Error("Register failed", "error", err, "email", body.Email)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.log.Info("User registered successfully", "email", body.Email)
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	h.log.Info("Get all users request", "search", search)

	users, err := h.service.GetAllUsers(c, search)
	if err != nil {
		h.log.Error("Failed to fetch users", "error", err, "search", search)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	h.log.Info("Users fetched successfully", "count", len(users))
	c.JSON(http.StatusOK, users)
}

func (h *Handler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation failed",
		})
		return
	}

	token, err := h.service.Login(c, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
