package handlers

import (
	"encoding/json"
	"fmt"
	"llm-service/configs"
	"llm-service/internal/model"
	"llm-service/internal/summarizer"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type summarizerCacheMock struct {
	store map[string]string
}

func newSummarizerCacheMock() *summarizerCacheMock {
	return &summarizerCacheMock{store: make(map[string]string)}
}

func (m *summarizerCacheMock) Set(key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.store[key] = string(b)
	return nil
}

func (m *summarizerCacheMock) Get(key string, dest any) error {
	raw, ok := m.store[key]
	if !ok {
		return fmt.Errorf("not found")
	}
	return json.Unmarshal([]byte(raw), dest)
}

func (m *summarizerCacheMock) Del(key string) error {
	delete(m.store, key)
	return nil
}

func (m *summarizerCacheMock) Ping() error { return nil }

func TestGetSummarizer_Success(t *testing.T) {
	original := summarizer.GetRuntimeConfig()
	defer func() {
		_ = summarizer.UpdateRuntimeConfig(summarizer.RuntimeConfigUpdate{
			Model:                    &original.Model,
			SystemPromptEN:           &original.SystemPromptEN,
			SystemPromptDE:           &original.SystemPromptDE,
			SystemPromptKZ:           &original.SystemPromptKZ,
			EnableParallelGeneration: &original.EnableParallelGeneration,
		})
	}()

	disableParallel := false
	if err := summarizer.UpdateRuntimeConfig(summarizer.RuntimeConfigUpdate{
		EnableParallelGeneration: &disableParallel,
	}); err != nil {
		t.Fatalf("failed to set runtime config: %v", err)
	}

	dataSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/work-experience", "/education", "/project", "/skill", "/certificate":
			_ = json.NewEncoder(w).Encode([]map[string]any{})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer dataSrv.Close()

	llmSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := model.SummarizerAPIResponse{
			Choices: []model.Choice{
				{Message: model.Message{Content: "generated summary"}},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer llmSrv.Close()

	cfg := &configs.Config{
		Services: configs.ServicesConfig{
			WorkExperienceURL: dataSrv.URL + "/work-experience",
			EducationURL:      dataSrv.URL + "/education",
			ProjectsURL:       dataSrv.URL + "/project",
			SkillsURL:         dataSrv.URL + "/skill",
			CertificatesURL:   dataSrv.URL + "/certificate",
		},
		Summarizer: configs.SummarizerConfig{
			APIKey: "dummy-key",
			APIURL: llmSrv.URL,
		},
	}

	h := New(cfg, newSummarizerCacheMock())
	r := ginTestRouter()
	r.GET("/summarize", h.GetSummarizer)

	req := httptest.NewRequest(http.MethodGet, "/summarize?lang=en", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got["message"] != "generated summary" {
		t.Fatalf("unexpected message: %#v", got["message"])
	}
}

func ginTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}
