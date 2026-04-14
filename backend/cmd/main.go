package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskflow/internal/config"
	"taskflow/internal/db"
	"taskflow/internal/logger"
	"taskflow/internal/middleware"
	"taskflow/internal/project"
	"taskflow/internal/task"
	"taskflow/internal/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Logger init
	log := logger.GetLogger()

	// DB init
	database := db.NewDB(cfg.DBUrl)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Info("Server running on port", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	userRepo := user.NewRepository(database, log)
	userService := user.NewService(userRepo, cfg, log)
	userHandler := user.NewHandler(userService, log)

	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)

	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(cfg))

	authRoutes.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userName, _ := c.Get("name")

		c.JSON(200, gin.H{
			"user_id": userID,
			"name":    userName,
		})
	})

	authRoutes.GET("/users", userHandler.GetAllUsers)

	projRepo := project.NewRepository(database, log)
	projService := project.NewService(projRepo, log)
	projHandler := project.NewHandler(projService, log)

	authRoutes.POST("/projects", projHandler.Create)
	authRoutes.GET("/projects", projHandler.List)
	authRoutes.GET("/projects/:id", projHandler.GetByID)
	authRoutes.PATCH("/projects/:id", projHandler.Update)
	authRoutes.DELETE("/projects/:id", projHandler.Delete)

	taskRepo := task.NewRepository(database)
	taskService := task.NewService(taskRepo, log)
	taskHandler := task.NewHandler(taskService, log)

	authRoutes.GET("/projects/:id/tasks", taskHandler.ListByProject)
	authRoutes.POST("/projects/:id/tasks", taskHandler.Create)

	authRoutes.PATCH("/tasks/:id", taskHandler.Update)
	authRoutes.DELETE("/tasks/:id", taskHandler.Delete)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	database.Close()

	log.Info("Server exited properly")
}
