package summarizer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"llm-service/internal/model"

	"github.com/stretchr/testify/require"
)

func TestClient_doLLMRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

			_ = json.NewEncoder(w).Encode(model.SummarizerAPIResponse{
				Choices: []model.Choice{{Message: model.Message{Content: "ok"}}},
			})
		}))
		defer srv.Close()

		s := &Summarizer{
			API_URL: srv.URL,
			Client:  srv.Client(),
		}

		resp, err := s.doLLMRequest(context.Background(), map[string]interface{}{
			"model": "test-model",
		}, map[string]string{
			"Authorization": "Bearer test-key",
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.Choices, 1)
		require.Equal(t, "ok", resp.Choices[0].Message.Content)
	})

	t.Run("non-200 response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "upstream failure", http.StatusBadGateway)
		}))
		defer srv.Close()

		s := &Summarizer{
			API_URL: srv.URL,
			Client:  srv.Client(),
		}

		_, err := s.doLLMRequest(context.Background(), map[string]interface{}{}, map[string]string{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "summarizer API error")
		require.Contains(t, err.Error(), "upstream failure")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{invalid json"))
		}))
		defer srv.Close()

		s := &Summarizer{
			API_URL: srv.URL,
			Client:  srv.Client(),
		}

		_, err := s.doLLMRequest(context.Background(), map[string]interface{}{}, map[string]string{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to decode summarizer response")
	})
}

func TestClient_fetchAllData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/work-experience":
				_, _ = w.Write([]byte(`[{"id":1,"company":"A","location":"L","startDate":"2020","endDate":"2021","position":"SE","description":"D","displayOrder":1,"techStack":"Go"}]`))
			case "/education":
				_, _ = w.Write([]byte(`[]`))
			case "/project":
				_, _ = w.Write([]byte(`[]`))
			case "/skill":
				_, _ = w.Write([]byte(`[]`))
			case "/certificate":
				_, _ = w.Write([]byte(`[]`))
			default:
				http.NotFound(w, r)
			}
		}))
		defer srv.Close()

		s := &Summarizer{
			DATA_URLS: []string{
				srv.URL + "/work-experience",
				srv.URL + "/education",
				srv.URL + "/project",
				srv.URL + "/skill",
				srv.URL + "/certificate",
			},
		}

		pd, err := s.fetchAllData()
		require.NoError(t, err)
		require.NotNil(t, pd)
		require.Len(t, pd.WorkExperience, 1)
		require.Equal(t, "A", pd.WorkExperience[0].Company)
	})

	t.Run("unknown endpoint key", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`[]`))
		}))
		defer srv.Close()

		s := &Summarizer{
			DATA_URLS: []string{srv.URL + "/unknown"},
		}
		_, err := s.fetchAllData()
		require.Error(t, err)
		require.Contains(t, err.Error(), "unknown key")
	})

	t.Run("invalid JSON for known key", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/work-experience") {
				_, _ = w.Write([]byte(`not-json`))
				return
			}
			_, _ = w.Write([]byte(`[]`))
		}))
		defer srv.Close()

		s := &Summarizer{
			DATA_URLS: []string{srv.URL + "/work-experience"},
		}

		_, err := s.fetchAllData()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to unmarshal work-experience")
	})
}
