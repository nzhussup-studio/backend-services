package main

import (
	"fmt"
	"image-service/internal/env"
	"time"

	_ "image-service/docs"
)

// @title Image Service API
// @version 1.0.0
// @description This is the API for managing image albums and uploads.

// @contact.name Nurzhanat Zhussup
// @contact.url https://www.linkedin.com/in/nurzhanat-zhussup/
// @contact.url https://github.com/nzhussup

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8085

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	discoveryConfig := &discoveryConfig{
		eurekaURL:   "http://discovery-server.default.svc.cluster.local:8761/eureka",
		appName:     "image-service",
		refreshRate: 30,
		servicesConfig: &servicesConfig{
			authService: "auth-service",
		},
	}

	var port int = 8085

	cfg := config{
		addr:            fmt.Sprintf(":%d", port),
		port:            port,
		storagePath:     "var/images",
		apiBasePath:     "/v1/album",
		discoveryConfig: discoveryConfig,
		redisConfig: &redisConfig{
			addr: fmt.Sprintf(
				"%s:%d",
				env.GetString("REDIS_HOST", "redis-service.default.svc.cluster.local"),
				env.GetInt("REDIS_PORT", 6379)),
			password: "",
			db:       0,
			duration: 24 * time.Hour,
		},
		apiGatewayURL: "https://api.nzhussup.com",
		kafkaConfig:   &kafkaConfig{
			// brokerList: []string{
			// 	env.GetString("KAFKA_BROKER_1", "kafka-broker-1.default.svc.cluster.local:29092"),
			// },
			// topic: env.GetString("KAFKA_TOPIC", "image-service"),
		},
	}

	secuirityCfg := GetSecurityConfig(&cfg)

	app := newApp(cfg, secuirityCfg)
	app.Redis.CheckHealth()

	app.Run()
}
