package configs

import (
	"fmt"
	"llm-service/internal/env"
	"time"

	"golang.org/x/time/rate"
)

type Config struct {
	App         AppConfig
	RateLimiter RateLimiterConfig
	Services    ServicesConfig
	Summarizer  SummarizerConfig
	Redis       RedisConfig
}

type AppConfig struct {
	Port            string
	Endpoint        string
	ShutdownTimeout time.Duration
}

type RateLimiterConfig struct {
	Rate     rate.Limit
	Burst    int
	Interval time.Duration
}

type ServicesConfig struct {
	WorkExperienceURL string
	EducationURL      string
	ProjectsURL       string
	SkillsURL         string
	CertificatesURL   string
}

type SummarizerConfig struct {
	APIKey string
	APIURL string
}

type RedisConfig struct {
	Addr     string
	Duration time.Duration
}

func NewConfigFromEnv() *Config {
	nginxGatewayURL := env.GetString("NGINX_GATEWAY_URL", "http://localhost:8082")

	return &Config{
		App: AppConfig{
			Port:            "8086",
			Endpoint:        "/v1/llm",
			ShutdownTimeout: 5 * time.Second,
		},
		RateLimiter: RateLimiterConfig{
			Rate:     5,
			Burst:    10,
			Interval: 3 * time.Minute,
		},
		Services: ServicesConfig{
			WorkExperienceURL: nginxGatewayURL + "/v1/base/work-experience",
			EducationURL:      nginxGatewayURL + "/v1/base/education",
			ProjectsURL:       nginxGatewayURL + "/v1/base/project",
			SkillsURL:         nginxGatewayURL + "/v1/base/skill",
			CertificatesURL:   nginxGatewayURL + "/v1/base/certificate",
		},
		Summarizer: SummarizerConfig{
			APIKey: env.GetString("OPENROUTER_API_KEY", ""),
			APIURL: env.GetString("OPENROUTER_API_URL", "https://openrouter.ai/api/v1/chat/completions"),
		},
		Redis: RedisConfig{
			Addr: fmt.Sprintf(
				"%s:%d",
				env.GetString("REDIS_HOST", "localhost"),
				env.GetInt("REDIS_PORT", 6379),
			),
			Duration: 30 * 24 * time.Hour,
		},
	}
}
