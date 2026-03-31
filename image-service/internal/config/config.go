package config

import (
	"fmt"
	"image-service/internal/env"
	"time"
)

type Config struct {
	Server struct {
		Addr string
		Port int
	}
	API struct {
		BasePath string
	}
	Storage struct {
		Path string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
		TTL      time.Duration
	}
	Security struct {
		JWKSetURL        string
		ExpectedIssuer   string
		ExpectedAudience string
		BackendClientID  string
	}
	Image ImageConfig
}

type ImageConfig struct {
	MaxUploadBytes int64
	ResizeWidth    uint
	JPEGQuality    int
}

func Load() *Config {
	cfg := &Config{}

	// Server
	port := env.GetInt("PORT", 8085)
	cfg.Server.Port = port
	cfg.Server.Addr = fmt.Sprintf(":%d", port)

	// API
	cfg.API.BasePath = env.GetString("API_BASE_PATH", "/v1/album")

	// Storage
	cfg.Storage.Path = env.GetString("STORAGE_PATH", "var/images")

	// Redis
	cfg.Redis.Addr = fmt.Sprintf(
		"%s:%d",
		env.GetString("REDIS_HOST", "localhost"),
		env.GetInt("REDIS_PORT", 6379),
	)
	cfg.Redis.Password = env.GetString("REDIS_PASSWORD", "")
	cfg.Redis.DB = env.GetInt("REDIS_DB", 0)
	cfg.Redis.TTL = time.Duration(env.GetInt("REDIS_TTL_HOURS", 24)) * time.Hour

	// Security
	cfg.Security.JWKSetURL = env.GetString(
		"KEYCLOAK_JWK_SET_URL",
		"http://localhost:8081/realms/backend-auth-dev/protocol/openid-connect/certs",
	)
	cfg.Security.ExpectedIssuer = env.GetString("KEYCLOAK_EXPECTED_ISSUER", "")
	cfg.Security.ExpectedAudience = env.GetString("KEYCLOAK_EXPECTED_AUDIENCE", "")
	cfg.Security.BackendClientID = env.GetString("KEYCLOAK_BACKEND_CLIENT_ID", "backend-auth-client")

	// Image processing
	cfg.Image = ImageConfig{
		MaxUploadBytes: int64(env.GetInt("IMAGE_MAX_UPLOAD_MB", 25)) * 1024 * 1024,
		ResizeWidth:    uint(env.GetInt("IMAGE_RESIZE_WIDTH", 800)),
		JPEGQuality:    env.GetInt("IMAGE_JPEG_QUALITY", 95),
	}

	return cfg
}
