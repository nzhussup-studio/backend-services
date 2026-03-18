package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestRateLimiter_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		rate       rate.Limit
		burst      int
		requests   int
		delay      time.Duration
		wantStatus int
	}{
		"allows under limit": {
			rate:       2,
			burst:      2,
			requests:   2,
			delay:      0,
			wantStatus: http.StatusOK,
		},
		"blocks over limit": {
			rate:       1,
			burst:      1,
			requests:   3,
			delay:      0,
			wantStatus: http.StatusTooManyRequests,
		},
		"resets after delay": {
			rate:       1,
			burst:      1,
			requests:   2,
			delay:      time.Second,
			wantStatus: http.StatusOK,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rl := NewRateLimiter(tt.rate, tt.burst, 1*time.Minute)

			router := gin.New()
			router.Use(rl.Middleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "ok"})
			})

			var resp *httptest.ResponseRecorder

			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.RemoteAddr = "127.0.0.1:1234"

				resp = httptest.NewRecorder()
				router.ServeHTTP(resp, req)

				if tt.delay > 0 {
					time.Sleep(tt.delay)
				}
			}

			if resp.Code != tt.wantStatus {
				t.Errorf("Final response code = %d; want %d", resp.Code, tt.wantStatus)
			}
		})
	}
}
