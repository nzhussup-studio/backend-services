package main

import (
	"fmt"
	"llm-service/internal/cache"
	"llm-service/internal/env"
	"log/slog"
	"time"
)

// @title Summarizer Service API
// @version 1.0.0
// @description This is the API for generating summary from user profile data.

// @contact.name Nurzhanat Zhussup
// @contact.url https://www.linkedin.com/in/nurzhanat-zhussup/
// @contact.url https://github.com/nzhussup

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8086

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	slog.Info("Starting LLM service...")

	API_GATEWAY_URL := env.GetString("API_GATEWAY_URL", "http://api-gateway.default.svc.cluster.local:8082")

	config := &config{
		appConfig: &appConfig{
			port:     "8086",
			endpoint: "/v1/llm",
		},
		ratelimiterConfig: &ratelimiterConfig{
			rate:     5,
			burst:    10,
			interval: 3 * time.Minute,
		},
		servicesConfig: &servicesConfig{
			WORK_EXPERIENCE_URL: API_GATEWAY_URL + "/v1/base/work-experience",
			EDUCATION_URL:       API_GATEWAY_URL + "/v1/base/education",
			PROJECTS_URL:        API_GATEWAY_URL + "/v1/base/project",
			SKILLS_URL:          API_GATEWAY_URL + "/v1/base/skill",
			CERTIFICATES_URL:    API_GATEWAY_URL + "/v1/base/certificate",
		},
		summarizerConfig: &summarizerConfig{
			API_KEY: env.GetString("OPENROUTER_API_KEY", ""),
			API_URL: env.GetString("OPENROUTER_API_URL", "https://openrouter.ai/api/v1/chat/completions"),
		},
		redisConfig: &redisConfig{
			addr: fmt.Sprintf(
				"%s:%d",
				env.GetString("REDIS_HOST", "redis-service.default.svc.cluster.local"),
				env.GetInt("REDIS_PORT", 6379),
			),
			duration: 30 * 24 * time.Hour, // 30 days
		},
	}

	rdb := cache.NewRedisClient(
		config.redisConfig.addr,
		"",
		0,
		config.redisConfig.duration,
	)

	slog.Info("Configuration loaded", slog.String("port", config.appConfig.port), slog.String("endpoint", config.appConfig.endpoint))

	app := NewApp(config, rdb)
	if err := app.run(); err != nil {
		slog.Error("Failed to run the application", slog.String("error", err.Error()))
	} else {
		slog.Info("Application is running", slog.String("port", config.appConfig.port))
	}
}
