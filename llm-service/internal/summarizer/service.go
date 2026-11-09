package summarizer

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"llm-service/internal/cache"
	"llm-service/internal/model"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	generationLease = 2 * time.Minute
	generationLimit = 90 * time.Second
	failedStateTTL  = time.Minute
)

var supportedLanguages = []string{"en", "kk", "de"}

type Summary struct {
	Text        string    `json:"text"`
	DataHash    string    `json:"data_hash"`
	ConfigHash  string    `json:"config_hash"`
	GeneratedAt time.Time `json:"generated_at"`
}

type Generation struct {
	JobID      string    `json:"job_id"`
	Status     string    `json:"status"`
	DataHash   string    `json:"data_hash"`
	ConfigHash string    `json:"config_hash"`
	StartedAt  time.Time `json:"started_at"`
	Error      string    `json:"error,omitempty"`
}

type Result struct {
	Status      string
	Language    string
	Summary     string
	PollAfterMS int
}

type Summarizer struct {
	API_KEY   string
	API_URL   string
	DATA_URLS []string
	Client    *http.Client
	redis     cache.Cacher
}

func NewSummarizer(apiKey, apiURL string, dataURLs []string, client *http.Client, redis cache.Cacher) *Summarizer {
	if client == nil {
		client = &http.Client{Timeout: generationLimit}
	}
	return &Summarizer{API_KEY: apiKey, API_URL: apiURL, DATA_URLS: dataURLs, Client: client, redis: redis}
}

func (s *Summarizer) Start(ctx context.Context, requestedLanguage string) (Result, error) {
	pd, err := s.fetchAllData(ctx)
	if err != nil {
		return Result{}, fmt.Errorf("failed to fetch personal data: %w", err)
	}
	dataHash, err := hashValue(pd)
	if err != nil {
		return Result{}, fmt.Errorf("failed to hash personal data: %w", err)
	}

	for _, language := range supportedLanguages {
		cfgHash := configHash(language)
		var cached Summary
		if err := s.redis.Get(summaryKey(language), &cached); err == nil && cached.DataHash == dataHash && cached.ConfigHash == cfgHash {
			slog.Info("summary_cache_hit", slog.String("language", language), slog.String("data_hash", shortHash(dataHash)))
			continue
		}

		_ = s.redis.Del(summaryKey(language))
		var active Generation
		if err := s.redis.Get(generationKey(language), &active); err == nil &&
			(active.Status == "failed" || active.DataHash != dataHash || active.ConfigHash != cfgHash) {
			_ = s.redis.Del(generationKey(language))
		}
		job := Generation{JobID: newJobID(), Status: "generating", DataHash: dataHash, ConfigHash: cfgHash, StartedAt: time.Now().UTC()}
		acquired, err := s.redis.SetNX(generationKey(language), job, generationLease)
		if err != nil {
			return Result{}, fmt.Errorf("failed to acquire %s generation: %w", language, err)
		}
		if !acquired {
			slog.Info("generation_already_running", slog.String("language", language))
			continue
		}

		slog.Info("generation_scheduled", slog.String("language", language), slog.String("job_id", job.JobID), slog.String("data_hash", shortHash(dataHash)))
		go s.generate(language, job, pd, dataHash, cfgHash)
	}

	return s.Get(requestedLanguage)
}

func (s *Summarizer) Get(language string) (Result, error) {
	var summary Summary
	if err := s.redis.Get(summaryKey(language), &summary); err == nil {
		return Result{Status: "completed", Language: language, Summary: summary.Text}, nil
	} else if !errors.Is(err, cache.ErrNotFound) {
		return Result{}, fmt.Errorf("failed to read cached summary: %w", err)
	}

	var generation Generation
	if err := s.redis.Get(generationKey(language), &generation); err == nil {
		result := Result{Status: generation.Status, Language: language}
		if generation.Status == "generating" {
			result.PollAfterMS = 1500
		}
		return result, nil
	} else if !errors.Is(err, cache.ErrNotFound) {
		return Result{}, fmt.Errorf("failed to read generation state: %w", err)
	}

	return Result{Status: "missing", Language: language}, nil
}

func (s *Summarizer) generate(language string, job Generation, pd *model.PersonalData, dataHash, cfgHash string) {
	started := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), generationLimit)
	defer cancel()

	slog.Info("generation_started", slog.String("language", language), slog.String("job_id", job.JobID))
	text, err := s.generateWithRetry(ctx, language, pd, job.JobID)
	if err != nil {
		slog.Error("generation_failed", slog.String("language", language), slog.String("job_id", job.JobID), slog.Duration("duration", time.Since(started)), slog.String("error", err.Error()))
		if s.ownsJob(language, job.JobID) {
			job.Status = "failed"
			job.Error = "generation_failed"
			_ = s.redis.SetWithTTL(generationKey(language), job, failedStateTTL)
		}
		return
	}

	if !s.ownsJob(language, job.JobID) {
		slog.Warn("generation_result_discarded", slog.String("language", language), slog.String("job_id", job.JobID))
		return
	}
	record := Summary{Text: text, DataHash: dataHash, ConfigHash: cfgHash, GeneratedAt: time.Now().UTC()}
	if err := s.redis.SetWithTTL(summaryKey(language), record, 0); err != nil {
		slog.Error("generation_cache_failed", slog.String("language", language), slog.String("job_id", job.JobID), slog.String("error", err.Error()))
		return
	}
	_ = s.redis.Del(generationKey(language))
	slog.Info("generation_completed", slog.String("language", language), slog.String("job_id", job.JobID), slog.Duration("duration", time.Since(started)))
}

func (s *Summarizer) generateWithRetry(ctx context.Context, language string, pd *model.PersonalData, jobID string) (string, error) {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		text, retry, err := s.generateOnce(ctx, language, pd)
		if err == nil {
			return text, nil
		}
		lastErr = err
		if !retry || attempt == 3 {
			break
		}
		delay := time.Duration(attempt) * time.Second
		slog.Warn("generation_retry", slog.String("language", language), slog.String("job_id", jobID), slog.Int("attempt", attempt), slog.Duration("delay", delay), slog.String("error", err.Error()))
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
		}
	}
	return "", lastErr
}

func (s *Summarizer) generateOnce(ctx context.Context, language string, pd *model.PersonalData) (string, bool, error) {
	systemPrompt, userPrompt, err := getPromptBase(language, pd)
	if err != nil {
		return "", false, err
	}
	payload := map[string]interface{}{
		"model":    GetRuntimeConfig().Model,
		"messages": []map[string]string{{"role": "system", "content": systemPrompt}, {"role": "user", "content": userPrompt}},
	}
	result, err := s.doLLMRequest(ctx, payload, map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + s.API_KEY})
	if err != nil {
		var upstream *UpstreamError
		retry := errors.As(err, &upstream) && (upstream.Status == http.StatusTooManyRequests || upstream.Status >= 500)
		if errors.Is(err, context.DeadlineExceeded) {
			retry = true
		}
		return "", retry, err
	}
	if result == nil || len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "", false, fmt.Errorf("no summary received from the API")
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), false, nil
}

func (s *Summarizer) ownsJob(language, jobID string) bool {
	var current Generation
	return s.redis.Get(generationKey(language), &current) == nil && current.JobID == jobID
}

func hashValue(value any) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

func configHash(language string) string {
	cfg := GetRuntimeConfig()
	prompt := cfg.SystemPromptEN
	if language == "kk" {
		prompt = cfg.SystemPromptKK
	}
	if language == "de" {
		prompt = cfg.SystemPromptDE
	}
	hash, _ := hashValue(struct{ Model, Prompt string }{cfg.Model, prompt})
	return hash
}

func newJobID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func summaryKey(language string) string    { return "llm:summary:" + language }
func generationKey(language string) string { return "llm:generation:" + language }
func shortHash(hash string) string {
	if len(hash) > 8 {
		return hash[:8]
	}
	return hash
}
