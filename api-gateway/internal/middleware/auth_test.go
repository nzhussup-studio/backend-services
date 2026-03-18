package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareWithMockClient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := NewTestClient(func(req *http.Request) *http.Response {
		// Check the request payload for the token
		var vr ValidationRequest
		err := json.NewDecoder(req.Body).Decode(&vr)
		assert.NoError(t, err)

		if vr.Token == "valid-token" {
			body := `{"username":"user1","roles":["admin","user"]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader("invalid token")),
			Header:     make(http.Header),
		}
	})

	router := gin.New()
	router.Use(AuthMiddlewareWithClient("http://mock-auth", mockClient))
	router.GET("/test", func(c *gin.Context) {
		username := c.Request.Header.Get("X-User-Username")
		roles := c.Request.Header.Get("X-User-Roles")
		c.JSON(http.StatusOK, gin.H{"username": username, "roles": roles})
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUser   string
		expectedRoles  string
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			expectedUser:   "",
			expectedRoles:  "",
		},
		{
			name:           "Invalid Authorization header format",
			authHeader:     "Token abc",
			expectedStatus: http.StatusOK,
			expectedUser:   "",
			expectedRoles:  "",
		},
		{
			name:           "Valid token",
			authHeader:     "Bearer valid-token",
			expectedStatus: http.StatusOK,
			expectedUser:   "user1",
			expectedRoles:  "admin,user",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var resp map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, resp["username"])
				assert.Equal(t, tt.expectedRoles, resp["roles"])
			} else {
				var resp map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				_, exists := resp["error"]
				assert.True(t, exists)
			}
		})
	}
}
