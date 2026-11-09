package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"llm-service/configs"
	"llm-service/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type handlerCache struct{ store map[string][]byte }

func newHandlerCache() *handlerCache                    { return &handlerCache{store: map[string][]byte{}} }
func (m *handlerCache) Set(key string, value any) error { return m.SetWithTTL(key, value, 0) }
func (m *handlerCache) SetWithTTL(key string, value any, _ time.Duration) error {
	b, err := json.Marshal(value)
	if err == nil {
		m.store[key] = b
	}
	return err
}
func (m *handlerCache) SetNX(key string, value any, _ time.Duration) (bool, error) {
	if _, ok := m.store[key]; ok {
		return false, nil
	}
	return true, m.Set(key, value)
}
func (m *handlerCache) Get(key string, dest any) error {
	b, ok := m.store[key]
	if !ok {
		return fmt.Errorf("%w: missing", cache.ErrNotFound)
	}
	return json.Unmarshal(b, dest)
}
func (m *handlerCache) Del(key string) error { delete(m.store, key); return nil }
func (m *handlerCache) Ping() error          { return nil }

func TestGetSummaryMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := New(&configs.Config{}, newHandlerCache())
	r := gin.New()
	r.GET("/summaries/:language", h.GetSummary)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/summaries/en", nil))
	require.Equal(t, http.StatusNotFound, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"status":"missing"`)
}

func TestGetSummaryRejectsInvalidLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := New(&configs.Config{}, newHandlerCache())
	r := gin.New()
	r.GET("/summaries/:language", h.GetSummary)
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/summaries/fr", nil))
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
