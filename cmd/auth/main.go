// cmd/auth/main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// _ "github.com/Axontik/comin-authentication-service/docs"
	"github.com/Axontik/comin-authentication-service/internal/auth/handler"
	authmiddleware "github.com/Axontik/comin-authentication-service/internal/auth/middleware"
	"github.com/Axontik/comin-authentication-service/internal/auth/repository"
	"github.com/Axontik/comin-authentication-service/internal/auth/service"
	commonmiddleware "github.com/Axontik/comin-authentication-service/internal/common/middleware"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	// "github.com/swaggo/gin-swagger"
)

// @title           ComIn Authentication Service API
// @version         1.0
// @description     Authentication service for ComIn platform
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@comin.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type Services struct {
	jwtService  *jwt.JWTService
	authService *service.AuthService
	orgService  service.OrganizationService
	authHandler *handler.AuthHandler
	orgHandler  *handler.OrganizationHandler
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to init database:", err)
	}

	services := initServices(db)
	router := setupRouter(services)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() (*gorm.DB, error) {
	config := gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	return gorm.Open(postgres.Open("postgresql://comin_owner:Ye5rfjcIB7FX@ep-flat-shadow-a8onelva.eastus2.azure.neon.tech/comin?sslmode=require"), &config)
}

func initServices(db *gorm.DB) *Services {
	jwtService := jwt.NewJWTService(
		os.Getenv("JWT_SECRET_KEY"),
		"comin.auth",
	)

	authRepo := repository.NewAuthRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)

	authService := service.NewAuthService(authRepo, jwtService)
	orgService := service.NewOrganizationService(orgRepo)

	authHandler := handler.NewAuthHandler(authService)
	orgHandler := handler.NewOrganizationHandler(orgService)

	return &Services{
		jwtService:  jwtService,
		authService: authService,
		orgService:  *orgService,
		authHandler: authHandler,
		orgHandler:  orgHandler,
	}
}

func setupRouter(s *Services) *gin.Engine {
	router := gin.Default()

	router.Use(commonmiddleware.ErrorHandler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/refresh", s.authHandler.RefreshToken)
			auth.POST("/logout", s.authHandler.Logout)
			auth.GET("/validate", s.authHandler.ValidateToken)
		}

		protected := api.Group("/")
		protected.Use(authmiddleware.AuthMiddleware(s.jwtService))
		{
			org := protected.Group("/organizations")
			{
				org.POST("/", s.orgHandler.Create)
				org.GET("/profile", s.orgHandler.GetProfile)
				org.PUT("/profile", s.orgHandler.UpdateProfile)
			}
		}
	}

	return router
}
