package handlers

import (
	"errors"
	"llm-service/configs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type healthCacheMock struct {
	pingErr error
}

func (m *healthCacheMock) Set(string, any) error { return nil }
func (m *healthCacheMock) Get(string, any) error { return nil }
func (m *healthCacheMock) Del(string) error      { return nil }
func (m *healthCacheMock) Ping() error           { return m.pingErr }

func TestGetHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
	}{
		{name: "healthy", pingErr: nil, wantStatus: http.StatusOK},
		{name: "redis error", pingErr: errors.New("redis down"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(&configs.Config{}, &healthCacheMock{pingErr: tt.pingErr})

			r := gin.New()
			r.GET("/health", h.GetHealth)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
