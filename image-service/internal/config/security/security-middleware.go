package security

import (
	"errors"
	custom_errors "image-service/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function that checks the authorization token
// in the request header and validates it with the auth service.
// If the token is valid, it sets the username and roles in the context.
// If the token is invalid or missing, it returns a 401 Unauthorized response.
// It allows GET requests to pass through without validation.
func AuthMiddleware(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodGet && !shouldAuthenticate(c, config.Rules) {
			c.Next()
			return
		}

		if config.AuthServiceURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service URL is not set"})
			c.Abort()
			return
		}

		token, err := GetToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		validationResponse, err := Validate(c, config, token)
		if err != nil {
			switch {
			case errors.Is(err, custom_errors.ErrUnauthorized):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			case errors.Is(err, custom_errors.ErrInternalServer):
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Auth service unavailable"})
				c.Abort()
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}
		}

		if !isAdmin(validationResponse.Roles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Set("username", validationResponse.Username)
		c.Set("roles", validationResponse.Roles)

		c.Next()
	}
}
