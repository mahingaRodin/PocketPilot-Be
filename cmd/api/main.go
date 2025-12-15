// @title           PocketPilot API
// @version         1.0
// @description     Backend API for PocketPilot expense tracking app.

// @contact.name    pocketpilot
// @contact.email   mahingarodin@gmail.com

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host            localhost:909
// @BasePath        /api

// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
    "log"
    "pocketpilot/internal/config"
    "pocketpilot/internal/handlers"
    "pocketpilot/internal/middleware"
    "pocketpilot/internal/repository"
    "pocketpilot/internal/services"
    "pocketpilot/pkg/database"

    "github.com/gin-gonic/gin"

    _ "pocketpilot/docs" // <-- IMPORTANT: Side-effect import for generated Swagger docs (replace 'pocketpilot' with your module name from go.mod)

    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
)

func main() {
    // loading config
    cfg := config.Load()
    
    // db init
    db := database.Connect(cfg.DatabaseURL)
    defer db.Close()

    // repo init
    userRepo := repository.NewUserRepository(db.DB)
    // expenseRepo := repository.NewExpenseRepository(db.DB)
    
    // service init
    authService := services.NewAuthService(userRepo, cfg.JWTSecret)
    // expenseService := services.NewExpenseService(expenseRepo, userRepo)
    
    // handlers init
    authHandler := handlers.NewAuthHandler(authService)
    // expenseHandler := handlers.NewExpenseHandler(expenseService)
    
    // gin router
    router := gin.Default()
    
    // middleware
    router.Use(middleware.CORS())
    router.Use(middleware.RateLimit())

    // Add Swagger UI endpoint (before routes for easy access)
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // routes (NOW passing expenseHandler)
    setupRoutes(router, authHandler, cfg.JWTSecret)
    
    // start server
    log.Printf("Server starting on port %s", cfg.Port)
    log.Fatal(router.Run(":" + cfg.Port))
}

func setupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, jwtSecret string) {
    // Public auth routes
    router.POST("/api/auth/register", authHandler.Register)
    router.POST("/api/auth/login", authHandler.Login)

    // Protected group
    auth := router.Group("/api")
    auth.Use(middleware.AuthMiddleware(jwtSecret))

    // Auth profile
    auth.GET("/auth/profile", authHandler.GetProfile)

    // // Expense routes
    // expenses := auth.Group("/expenses")
    // {
    //     expenses.POST("/", expenseHandler.CreateExpense)
    //     expenses.GET("/", expenseHandler.GetExpenses)
    //     expenses.GET("/:id", expenseHandler.GetExpense)
    //     expenses.PUT("/:id", expenseHandler.UpdateExpense)
    //     expenses.DELETE("/:id", expenseHandler.DeleteExpense)
    //     expenses.GET("/team/:teamId", expenseHandler.GetTeamExpenses)
    // }

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "message": "Server is running",
        })
    })
}