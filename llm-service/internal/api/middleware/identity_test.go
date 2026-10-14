package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireGatewayIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		userID string
		want   int
	}{
		{name: "missing identity", want: http.StatusUnauthorized},
		{name: "trusted identity", userID: "user-1", want: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/configuration", RequireGatewayIdentity(), func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })
			req := httptest.NewRequest(http.MethodGet, "/configuration", nil)
			req.Header.Set("X-User-ID", tt.userID)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, response.Code)
			}
		})
	}
}
