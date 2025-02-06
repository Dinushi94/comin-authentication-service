// internal/common/middleware/error_middleware.go
package middleware

import (
	"net/http"

	"github.com/Axontik/comin-authentication-service/internal/common/errors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(errors.AppError); ok {
				status := getHTTPStatus(appErr.Code)
				c.JSON(status, appErr)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
	}
}

func getHTTPStatus(code errors.ErrorCode) int {
	switch code {
	case errors.ErrInvalidInput, errors.ErrValidation:
		return http.StatusBadRequest
	case errors.ErrUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrForbidden:
		return http.StatusForbidden
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrDuplicateEntity:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
