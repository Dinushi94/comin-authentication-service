// internal/auth/middleware/auth_middleware.go
package middleware

import (
	"net/http"
	"strings"

	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		claims, err := jwtService.ValidateToken(tokenParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("organizationID", claims.OrganizationID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
			return
		}

		for _, role := range roles {
			if role == userRole.(string) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
	}
}

func RequireOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID, exists := c.Get("organizationID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
			return
		}

		if orgID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid organization"})
			return
		}

		c.Next()
	}
}
