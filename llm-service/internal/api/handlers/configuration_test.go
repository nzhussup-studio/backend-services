package handlers

import (
	"bytes"
	"encoding/json"
	"llm-service/configs"
	"llm-service/internal/summarizer"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type noopCache struct{}

func (noopCache) Set(string, any) error { return nil }
func (noopCache) Get(string, any) error { return nil }
func (noopCache) Del(string) error      { return nil }
func (noopCache) Ping() error           { return nil }
func str(v string) *string              { return &v }
func b(v bool) *bool                    { return &v }
func newHandlerForConfig() *Handler     { return New(&configs.Config{}, noopCache{}) }
func setGinTestMode()                   { gin.SetMode(gin.TestMode) }
func mustUnmarshalMap(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}

func restoreConfig(t *testing.T, cfg summarizer.RuntimeConfig) {
	t.Helper()
	err := summarizer.UpdateRuntimeConfig(summarizer.RuntimeConfigUpdate{
		Model:                    str(cfg.Model),
		SystemPromptEN:           str(cfg.SystemPromptEN),
		SystemPromptDE:           str(cfg.SystemPromptDE),
		SystemPromptKZ:           str(cfg.SystemPromptKZ),
		EnableParallelGeneration: b(cfg.EnableParallelGeneration),
	})
	if err != nil {
		t.Fatalf("failed to restore runtime config: %v", err)
	}
}

func TestConfigurationHandlers(t *testing.T) {
	setGinTestMode()
	original := summarizer.GetRuntimeConfig()
	defer restoreConfig(t, original)

	t.Run("get configuration", func(t *testing.T) {
		h := newHandlerForConfig()
		r := gin.New()
		r.GET("/configuration", h.GetConfiguration)

		req := httptest.NewRequest(http.MethodGet, "/configuration", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		body := mustUnmarshalMap(t, rec.Body.Bytes())
		if _, ok := body["model"]; !ok {
			t.Fatalf("expected model in response")
		}
	})

	t.Run("put configuration success", func(t *testing.T) {
		h := newHandlerForConfig()
		r := gin.New()
		r.PUT("/configuration", h.PutConfiguration)

		payload := map[string]any{
			"model":                      "openai/gpt-4.1-mini",
			"system_prompt_en":           "Custom EN",
			"enable_parallel_generation": false,
		}
		raw, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/configuration", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
		}

		cfg := summarizer.GetRuntimeConfig()
		if cfg.Model != "openai/gpt-4.1-mini" || cfg.SystemPromptEN != "Custom EN" || cfg.EnableParallelGeneration {
			t.Fatalf("runtime config not updated as expected: %#v", cfg)
		}
	})

	t.Run("put configuration invalid payload", func(t *testing.T) {
		h := newHandlerForConfig()
		r := gin.New()
		r.PUT("/configuration", h.PutConfiguration)

		req := httptest.NewRequest(http.MethodPut, "/configuration", bytes.NewReader([]byte(`{"model":"   "}`)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})
}
