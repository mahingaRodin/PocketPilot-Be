package main

import (
	// "expense-tracker/internal/handlers"
	// "expense-tracker/internal/middleware"
	"log"
	"pocketpilot-api/internal/config"
	"pocketpilot-api/internal/middleware"
	"pocketpilot-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
    // load config
    cfg := config.Load()
    
    // db init
    db := database.Connect(cfg.DatabaseURL)
    defer db.Close()
    
    // gin router setup
    router := gin.Default()
    
    // midware setup
    router.Use(middleware.CORS())
    router.Use(middleware.RateLimit())
    
    // Routes
    setupRoutes(router, db)
    
    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(router.Run(":" + cfg.Port))
}

func setupRoutes(router *gin.Engine, db *database.DB) {
    router.GET("/health", func(c *gin.Context) {
        err := db.Ping()
        if err != nil {
            c.JSON(500, gin.H{
                "status": "error",
                "message": "Database connection failed",
            })
            return
        }
        
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running",
        })
    })
    
//routes to be defined here
};