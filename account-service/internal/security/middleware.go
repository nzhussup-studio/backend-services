package security

import (
	"errors"
	"net/http"

	customerrors "account-service/internal/errors"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.JWKSetURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Keycloak auth configuration is incomplete"})
			c.Abort()
			return
		}

		token, err := GetToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		validationResponse, err := Validate(c, config, token)
		if err != nil {
			switch {
			case errors.Is(err, customerrors.ErrUnauthorized):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Keycloak validation failed"})
			}
			c.Abort()
			return
		}

		c.Set("subject", validationResponse.Subject)
		c.Set("username", validationResponse.Username)
		c.Next()
	}
}
