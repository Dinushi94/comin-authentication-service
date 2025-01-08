// cmd/auth/main.go
package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"

    "github.com/yourusername/comin/internal/auth/handler"
    "github.com/yourusername/comin/internal/auth/repository"
    "github.com/yourusername/comin/internal/auth/service"
    "github.com/yourusername/comin/pkg/jwt"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // Database connection
    dbURL := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize dependencies
    jwtService := jwt.NewJWTService(
        os.Getenv("JWT_SECRET_KEY"),
        "comin.auth",
    )

    authRepo := repository.NewAuthRepository(db)
    authService := service.NewAuthService(authRepo, jwtService)
    authHandler := handler.NewAuthHandler(authService)

    // Setup Gin router
    router := gin.Default()

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // Auth routes
    auth := router.Group("/api/v1/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.RefreshToken)
        auth.POST("/logout", authHandler.Logout)
    }

    // Start server
    port := os.Getenv("API_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Starting auth service on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}