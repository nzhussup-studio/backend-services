package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"llm-service/internal/model"
	"llm-service/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	data = model.PersonalData{
		WorkExperience: []*model.WorkExperience{
			{
				ID:          1,
				Company:     "Google",
				Location:    "Mountain View",
				StartDate:   "2021-01-01",
				EndDate:     "2023-01-01",
				Position:    "Software Engineer",
				Description: "Worked on scalable backend systems.",
				TechStack:   "Go, Kubernetes, GCP",
			},
		},
		Education: []*model.Education{
			{
				ID:          1,
				Institution: "MIT",
				Location:    "Cambridge, MA",
				StartDate:   time.Date(2016, time.September, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     utils.ToPtrTime(time.Date(2020, time.May, 31, 0, 0, 0, 0, time.UTC)),
				Degree:      "B.Sc. in Computer Science",
				Thesis:      "Efficient Parallel Algorithms",
			},
		},
		Projects: []*model.Project{
			{
				ID:        1,
				Name:      "SmartAI",
				URL:       "https://github.com/user/smartai",
				TechStack: "Go, TensorFlow",
			},
		},
		Skills: []*model.Skill{
			{
				ID:         1,
				Category:   "Programming Languages",
				SkillNames: "Go, Python, JavaScript",
			},
		},
		Certificates: []*model.Certificate{
			{
				ID:   1,
				Name: "AWS Certified Solutions Architect",
				URL:  "https://aws.amazon.com/certification/",
			},
		},
	}
)

// Flexible mockRedis with optional injected behavior
type mockRedis struct {
	store   map[string]string
	setFunc func(key string, value any) error
	getFunc func(key string, dest any) error
	delFunc func(key string) error
}

func newMockRedis() *mockRedis {
	return &mockRedis{store: make(map[string]string)}
}

func (m *mockRedis) Set(key string, value any) error {
	if m.setFunc != nil {
		return m.setFunc(key, value)
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.store[key] = string(bytes)
	return nil
}

func (m *mockRedis) Get(key string, dest any) error {
	if m.getFunc != nil {
		return m.getFunc(key, dest)
	}
	val, ok := m.store[key]
	if !ok {
		return fmt.Errorf("key not found")
	}
	return json.Unmarshal([]byte(val), dest)
}

func (m *mockRedis) Del(key string) error {
	if m.delFunc != nil {
		return m.delFunc(key)
	}
	delete(m.store, key)
	return nil
}

func modifiedData() model.PersonalData {
	d := data
	d.WorkExperience[0].Company = "Microsoft"
	return d
}

func mockDataServer(t *testing.T, data model.PersonalData) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp interface{}
		switch r.URL.Path {
		case "/work-experience":
			resp = data.WorkExperience
		case "/education":
			resp = data.Education
		case "/project":
			resp = data.Projects
		case "/skill":
			resp = data.Skills
		case "/certificate":
			resp = data.Certificates
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}))
}

func mockLLMServer(t *testing.T, response string, fail bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fail {
			http.Error(w, "model error", http.StatusInternalServerError)
			return
		}
		resp := model.SummarizerAPIResponse{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Content: response,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestSummarizer_Summarize_Multilingual_Failures(t *testing.T) {
	tests := []struct {
		name           string
		lang           string
		data           model.PersonalData
		llmContent     string
		llmFail        bool
		setup          func(redis *mockRedis)
		overrideAPIURL bool
		expectErr      bool
		errMatch       string
	}{
		{
			name:      "LLM API returns non-200 response",
			lang:      "en",
			data:      data,
			llmFail:   true,
			expectErr: true,
			errMatch:  "model error",
		},
		{
			name:       "LLM response is empty",
			lang:       "kz",
			data:       data,
			llmContent: "",
			expectErr:  true,
			errMatch:   "no summary received",
		},
		{
			name:           "Broken personal data fetch (wrong URL)",
			lang:           "de",
			data:           data,
			expectErr:      true,
			overrideAPIURL: true,
			errMatch:       "failed to fetch personal data: failed to unmarshal work-experience: invalid character 'p' after top-level value",
		},
		{
			name: "Cache get returns invalid data",
			lang: "en",
			data: data,
			setup: func(redis *mockRedis) {
				redis.store["previous_personal_data"] = "not a json object"
			},
			llmContent: "Fallback summary for broken cache",
			expectErr:  false,
		},
		{
			name: "Cache set fails silently",
			lang: "kz",
			data: data,
			setup: func(redis *mockRedis) {
				redis.setFunc = func(key string, value any) error {
					return fmt.Errorf("forced cache failure")
				}
			},
			llmContent: "Summary with cache write failure",
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := newMockRedis()
			if tt.setup != nil {
				tt.setup(redis)
			}

			var dataSrv *httptest.Server
			if tt.overrideAPIURL {
				dataSrv = httptest.NewServer(http.NotFoundHandler())
			} else {
				dataSrv = mockDataServer(t, tt.data)
			}
			defer dataSrv.Close()

			llmSrv := mockLLMServer(t, tt.llmContent, tt.llmFail)
			defer llmSrv.Close()

			s := NewSummarizer("dummy-key", llmSrv.URL, []string{
				dataSrv.URL + "/work-experience",
				dataSrv.URL + "/education",
				dataSrv.URL + "/project",
				dataSrv.URL + "/skill",
				dataSrv.URL + "/certificate",
			}, llmSrv.Client(), redis, tt.lang)

			summary, err := s.Summarize(context.Background())

			if tt.expectErr {
				require.Error(t, err)
				if tt.errMatch != "" {
					require.Contains(t, err.Error(), tt.errMatch)
				}
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, summary)
			}
		})
	}
}

func TestSummarizer_Summarize_SuccessCases(t *testing.T) {
	tests := []struct {
		name             string
		lang             string
		data             model.PersonalData
		primeCache       bool
		cacheSummary     string
		llmContent       string
		changeData       bool
		expectedContains string
	}{
		{
			name:             "Generate new summary in English",
			lang:             "en",
			data:             data,
			llmContent:       "Nurzhanat Zhussup is a software engineer with experience at Google.",
			expectedContains: "Nurzhanat Zhussup",
		},
		{
			name:             "Generate new summary in Kazakh",
			lang:             "kz",
			data:             data,
			llmContent:       "Нұржанат Жүсіп — Google компаниясында тәжірибесі бар бағдарламашы.",
			expectedContains: "Нұржанат",
		},
		{
			name:             "Generate new summary in German",
			lang:             "de",
			data:             data,
			llmContent:       "Nurzhanat Zhussup ist ein Softwareentwickler mit Erfahrung bei Google.",
			expectedContains: "Softwareentwickler",
		},
		{
			name:             "Use cached summary if data unchanged",
			lang:             "en",
			data:             data,
			primeCache:       true,
			cacheSummary:     "Cached summary of Nurzhanat from Redis",
			expectedContains: "Cached summary",
		},
		{
			name:             "Regenerate summary if data changed",
			lang:             "en",
			data:             modifiedData(),
			llmContent:       "Summary updated after company change to Microsoft.",
			expectedContains: "Microsoft",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := newMockRedis()

			// Prime cache if necessary
			if tt.primeCache {
				_ = redis.Set("previous_personal_data", &data)
				cacheKey := "summary_" + tt.lang
				_ = redis.Set(cacheKey, tt.cacheSummary)
			}

			dataSrv := mockDataServer(t, tt.data)
			defer dataSrv.Close()

			llmSrv := mockLLMServer(t, tt.llmContent, false)
			defer llmSrv.Close()

			s := NewSummarizer("dummy-key", llmSrv.URL, []string{
				dataSrv.URL + "/work-experience",
				dataSrv.URL + "/education",
				dataSrv.URL + "/project",
				dataSrv.URL + "/skill",
				dataSrv.URL + "/certificate",
			}, llmSrv.Client(), redis, tt.lang)

			summary, err := s.Summarize(context.Background())
			require.NoError(t, err)
			require.Contains(t, summary, tt.expectedContains)
		})
	}
}
