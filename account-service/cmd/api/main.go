package main

import (
	"fmt"
	"time"

	"account-service/internal/env"
	"account-service/internal/keycloak"
	"account-service/internal/service"

	_ "account-service/docs"
)

// @title Account Service API
// @version 1.0.0
// @description API for authenticated self-service account management through Keycloak.
// @contact.name Nurzhanat Zhussup
// @contact.url https://github.com/nzhussup
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8087
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	port := env.GetInt("PORT", 8087)

	cfg := config{
		addr:        fmt.Sprintf(":%d", port),
		port:        port,
		apiBasePath: env.GetString("API_BASE_PATH", "/v1/account"),
		keycloakConfig: &keycloakConfig{
			baseURL:           env.GetString("KEYCLOAK_BASE_URL", "http://localhost:8081"),
			realm:             env.GetString("KEYCLOAK_REALM", "backend-auth-dev"),
			adminRealm:        env.GetString("KEYCLOAK_ADMIN_REALM", "master"),
			adminClientID:     env.GetString("KEYCLOAK_ADMIN_CLIENT_ID", "admin-cli"),
			adminClientSecret: env.GetString("KEYCLOAK_ADMIN_CLIENT_SECRET", ""),
			httpTimeout:       10 * time.Second,
		},
	}

	adminClient := keycloak.NewAdminClient(keycloak.AdminClientConfig{
		BaseURL:           cfg.keycloakConfig.baseURL,
		Realm:             cfg.keycloakConfig.realm,
		AdminRealm:        cfg.keycloakConfig.adminRealm,
		AdminClientID:     cfg.keycloakConfig.adminClientID,
		AdminClientSecret: cfg.keycloakConfig.adminClientSecret,
		HTTPTimeout:       cfg.keycloakConfig.httpTimeout,
	})

	accountService := service.NewAccountService(adminClient)
	app := newApp(cfg, accountService)
	app.Run()
}
