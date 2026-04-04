package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llm-service/internal/model"
	"net/http"
	"strings"
)

// sends a request to the summarizer API with the given payload and headers and returns the response
func (s *Summarizer) doLLMRequest(ctx context.Context, payload map[string]interface{}, headers map[string]string) (*model.SummarizerAPIResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.API_URL, io.NopCloser(bytes.NewReader(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to summarizer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("summarizer API error: %s", string(bodyBytes))
	}

	result := &model.SummarizerAPIResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode summarizer response: %w", err)
	}
	return result, nil
}

// fetches all personal data from the configured endpoints
func (s *Summarizer) fetchAllData() (*model.PersonalData, error) {
	pd := &model.PersonalData{}

	for _, endpoint := range s.DATA_URLS {
		keysSlice := strings.Split(endpoint, "/")
		if len(keysSlice) == 0 {
			return nil, fmt.Errorf("invalid endpoint format: %s", endpoint)
		}
		key := keysSlice[len(keysSlice)-1]

		resp, err := http.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s data: %w", key, err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for %s: %w", key, err)
		}

		switch key {
		case "work-experience":
			if err := json.Unmarshal(bodyBytes, &pd.WorkExperience); err != nil {
				return nil, fmt.Errorf("failed to unmarshal work-experience: %w", err)
			}
		case "education":
			if err := json.Unmarshal(bodyBytes, &pd.Education); err != nil {
				return nil, fmt.Errorf("failed to unmarshal education: %w", err)
			}
		case "project":
			if err := json.Unmarshal(bodyBytes, &pd.Projects); err != nil {
				return nil, fmt.Errorf("failed to unmarshal projects: %w", err)
			}
		case "skill":
			if err := json.Unmarshal(bodyBytes, &pd.Skills); err != nil {
				return nil, fmt.Errorf("failed to unmarshal skills: %w", err)
			}
		case "certificate":
			if err := json.Unmarshal(bodyBytes, &pd.Certificates); err != nil {
				return nil, fmt.Errorf("failed to unmarshal certificates: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	return pd, nil
}
