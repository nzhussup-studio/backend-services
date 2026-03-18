package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llm-service/internal/cache"
	"llm-service/internal/model"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const wrapper = "%w: %v"

var cacheKeys = map[string]string{
	"kz": "summary_kz",
	"de": "summary_de",
	"en": "summary_en",
}

var promptBaseMap = map[string]string{
	"kz": SYSTEM_PROMPT_KZ,
	"de": SYSTEM_PROMPT_DE,
	"en": SYSTEM_PROMPT_EN,
}

var (
	ErrFailedToFetchPersonalData = fmt.Errorf("failed to fetch personal data")
	ErrFailedToFetchPromptBase   = fmt.Errorf("failed to fetch prompt base data")
	ErrFailedToGetSummary        = fmt.Errorf("failed to get summary from LLM")
	ErrNoSummaryReceived         = fmt.Errorf("no summary received from the API")
	ErrCacheFailure              = fmt.Errorf("cache failure: %w", cache.ErrCache)
	ErrFailedToGetCachedSummary  = fmt.Errorf("failed to get cached summary")
)

type SummarizerAPI interface {
	Summarize() (string, error)
}

type Summarizer struct {
	API_KEY   string
	API_URL   string
	DATA_URLS []string
	Client    *http.Client
	redis     cache.Cacher
	lang      string
}

func NewSummarizer(apiKey, apiURL string, dataURLS []string, client *http.Client, redis cache.Cacher, lang string) *Summarizer {
	return &Summarizer{
		API_KEY:   apiKey,
		API_URL:   apiURL,
		DATA_URLS: dataURLS,
		Client:    client,
		redis:     redis,
		lang:      lang,
	}
}

// Summarize fetches personal data from configured endpoints, generates prompts, and requests a summary from the LLM API
func (s *Summarizer) Summarize(ctx context.Context) (string, error) {
	pd, err := s.fetchAllData()
	if err != nil {
		return "", fmt.Errorf(wrapper, ErrFailedToFetchPersonalData, err)
	}

	dataUnchanged := s.checkDataUnchanged(ctx, pd)
	if dataUnchanged {
		if cachedSummary := s.fetchCachedSummary(); cachedSummary != "" {
			return cachedSummary, nil
		}
	} else {
		s.flushAllSummaryCache()

		if err := s.redis.Set("previous_personal_data", pd); err != nil {
			slog.Error("Failed to cache previous personal data", slog.String("error", err.Error()))
		}
	}

	// Immediate summarization for the current language
	currentSummary, err := s.doFullSummarization(ctx, pd)
	if err != nil {
		return "", fmt.Errorf(wrapper, ErrFailedToGetSummary, err)
	}

	// Background summaries for other languages
	go s.backgroundSummarizeOthers(context.Background(), pd)

	return currentSummary, nil
}

// backgroundSummarizeOthers generates summaries for languages other than s.lang, caches them
func (s *Summarizer) backgroundSummarizeOthers(ctx context.Context, pd *model.PersonalData) {
	langs := []string{"en", "kz", "de"}

	for _, lang := range langs {
		if lang == s.lang {
			continue // Skip current language, already processed
		}

		time.Sleep(3 * time.Second) // Add delay between goroutines, to avoid rate limitation issues

		go func(lang string) {
			// Create a new Summarizer instance with the other language
			otherSummarizer := NewSummarizer(s.API_KEY, s.API_URL, s.DATA_URLS, s.Client, s.redis, lang)

			// Run summarization ignoring errors (just log)
			summary, err := otherSummarizer.doFullSummarization(ctx, pd)
			if err != nil {
				slog.Error("Background summarization failed",
					slog.String("lang", lang),
					slog.String("error", err.Error()))
				return
			}

			// Cache the summary (doFullSummarization caches it too, but we do again here just in case)
			if err := s.redis.Set(otherSummarizer.getCacheKey(), summary); err != nil {
				slog.Error("Failed to cache background summary",
					slog.String("lang", lang),
					slog.String("error", err.Error()))
			}
		}(lang)
	}
}

func (s *Summarizer) fetchCachedSummary() string {
	var cachedSummary string

	if err := s.redis.Get(s.getCacheKey(), &cachedSummary); err != nil {
		return ""
	}
	return cachedSummary
}

func (s *Summarizer) doFullSummarization(ctx context.Context, pd *model.PersonalData) (string, error) {
	systemPrompt, userPrompt, err := s.getPromptBase(pd)
	if err != nil {
		return "", fmt.Errorf(wrapper, ErrFailedToFetchPromptBase, err)
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", s.API_KEY),
	}

	payload := map[string]interface{}{
		"model": "tngtech/deepseek-r1t2-chimera:free",
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
	}

	result, err := s.doLLMRequest(ctx, payload, headers)
	if err != nil {
		return "", fmt.Errorf(wrapper, ErrFailedToGetSummary, err)
	}

	if result == nil || len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
		return "", ErrNoSummaryReceived
	}

	// Cache the summary for future requests
	if err := s.redis.Set(s.getCacheKey(), result.Choices[0].Message.Content); err != nil {
		slog.Error("Failed to cache summary", slog.String("error", err.Error()))
	}
	return result.Choices[0].Message.Content, nil
}

// checks if the personal data has NOT changed by comparing it with the cached version.
// Returns true if data is unchanged, false otherwise. Defaults to false if cache is not found.
func (s *Summarizer) checkDataUnchanged(ctx context.Context, currentData *model.PersonalData) bool {
	const cacheKey = "previous_personal_data"

	var previousData model.PersonalData
	if err := s.redis.Get(cacheKey, &previousData); err != nil {
		return false // If cache is not found, assume data has changed → return false (NOT unchanged)
	}

	prevJSON, err := json.Marshal(previousData)
	if err != nil {
		return false // Can't marshal previous data, assume changed → false
	}
	currJSON, err := json.Marshal(currentData)
	if err != nil {
		return false // Can't marshal current data, assume changed → false
	}

	if bytes.Equal(prevJSON, currJSON) {
		return true // Data is unchanged
	}

	if err = s.redis.Set(cacheKey, currentData); err != nil {
		return false // Can't set cache, assume changed → false
	}

	return false // Data changed
}

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

func (s *Summarizer) flushAllSummaryCache() {
	if err := s.redis.Del("summary_en"); err != nil {
		slog.Error("error deleting summary cache for English", slog.Any("error", err))
	}
	if err := s.redis.Del("summary_kz"); err != nil {
		slog.Error("error deleting summary cache for Kazakh", slog.Any("error", err))
	}
	if err := s.redis.Del("summary_de"); err != nil {
		slog.Error("error deleting summary cache for German", slog.Any("error", err))
	}
}

// gets right cache based on the language setting
func (s *Summarizer) getCacheKey() string {
	if key, ok := cacheKeys[s.lang]; ok {
		return key
	}
	slog.Warn("unsupported language, defaulting to English", slog.String("lang", s.lang))
	return "summary_en"
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

// generates system and user prompts based on the personal data
func (s *Summarizer) getPromptBase(pd *model.PersonalData) (string, string, error) {
	systemPrompt, ok := promptBaseMap[s.lang]
	if !ok {
		slog.Warn("unsupported language, defaulting to English prompt", slog.String("lang", s.lang))
		systemPrompt = SYSTEM_PROMPT_EN
	}

	dataJSON, err := json.MarshalIndent(pd, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal personal data: %w", err)
	}

	userPrompt := fmt.Sprintf(USER_PROMPT, string(dataJSON))

	return systemPrompt, userPrompt, nil
}
