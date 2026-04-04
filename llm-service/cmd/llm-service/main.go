package main

import (
	"llm-service/configs"
	"llm-service/internal/api"
	"llm-service/internal/cache"
	"log/slog"
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

	config := configs.NewConfigFromEnv()

	rdb := cache.NewRedisClient(
		config.Redis.Addr,
		"",
		0,
		config.Redis.Duration,
	)

	slog.Info("Configuration loaded", slog.String("port", config.App.Port), slog.String("endpoint", config.App.Endpoint))

	server := api.NewServer(config, rdb)
	if err := server.Run(); err != nil {
		slog.Error("Failed to run the application", slog.String("error", err.Error()))
	} else {
		slog.Info("Application is running", slog.String("port", config.App.Port))
	}
}
