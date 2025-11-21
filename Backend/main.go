package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phanindra08/snap-attend/internal/features/controller"
	"github.com/phanindra08/snap-attend/internal/shared/config"
	"github.com/phanindra08/snap-attend/internal/shared/database"
	"github.com/phanindra08/snap-attend/internal/shared/middleware"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

func main() {
	config.GetConfig()
	database.GetDB()
	defer database.CloseDB()

	err := database.HealthCheck()
	if err != nil {
		return
	}

	// Check the environment variable SNAP_ATTEND_ENV for figuring out the gin mode to be used
	if os.Getenv(utils.ENVIRONMENT_VARIABLE) == utils.DEV_ENVIRONMENT {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middlewares
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	controller := controller.GetControllers()
	// Registering routes
	registerRoutes(router, controller)

	// Create HTTP server with graceful shutdown
	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.GetConfig().GetServerPort()),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Starting the server in a goroutine
	go func() {
		log.Printf("Server starting on port %d.", config.GetConfig().GetServerPort())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attempting shutdown gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited properly")
}

func registerRoutes(router *gin.Engine, controller *controller.Controller) {
	api := router.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(context *gin.Context) {
			if err := database.HealthCheck(); err != nil {
				context.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
				return
			}
			context.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", controller.User.Signup)
			auth.POST("/login", controller.User.Login)
		}

		user := api.Group("/user")
		user.Use(middleware.JWTAuthMiddleware())
		{
			user.PUT("/update-profile", controller.User.UpdateProfile)
			user.POST("/logout", controller.User.Logout)
		}
		//staff := auth.Group("/staff")
		//staff.Use(middleware.RequireRoles(models.Admin, models.Professor))
		//staff.GET("/reports", controllers.StaffReports)
	}
}
