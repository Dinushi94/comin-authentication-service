// internal/auth/handler/auth_handler.go
package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/service"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user and organization
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration details"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} domain.AuthResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
		return
	}

	// Remove Bearer prefix
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token) 
	log.Printf("Token after removing Bearer prefix: %s", token)

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":         claims.UserID,
		"organization_id": claims.OrganizationID,
		"email":           claims.Email,
		"role":            claims.Role,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear client-side tokens
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Middleware for protected routes
func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set claims to context
		c.Set("userID", claims.UserID)
		c.Set("organizationID", claims.OrganizationID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
