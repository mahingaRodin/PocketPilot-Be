package main

import (
	"log"
	"pocketpilot-api/internal/config"
	"pocketpilot-api/internal/handlers"
	"pocketpilot-api/internal/middleware"
	"pocketpilot-api/internal/repository"
	"pocketpilot-api/internal/services"
	"pocketpilot-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Initialize database
    db := database.Connect(cfg.DatabaseURL)
    defer db.Close()
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(db.DB)
    
    // Initialize services
    authService := services.NewAuthService(userRepo, cfg.JWTSecret)
    
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    
    // Initialize Gin router
    router := gin.Default()
    
    // Middleware
    router.Use(middleware.CORS())
    router.Use(middleware.RateLimit())
    
    // Routes
    setupRoutes(router, authHandler, cfg.JWTSecret)
    
    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(router.Run(":" + cfg.Port))
}

func setupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, jwtSecret string) {
    // Public routes
    router.POST("/api/auth/register", authHandler.Register)
    router.POST("/api/auth/login", authHandler.Login)
    
    // Protected routes
    auth := router.Group("/api")
    auth.Use(middleware.AuthMiddleware(jwtSecret))
    {
        auth.GET("/auth/profile", authHandler.GetProfile)
    }
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running",
        })
    })
}