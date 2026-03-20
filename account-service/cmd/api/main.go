package main

import (
	"fmt"
	"time"

	"account-service/internal/env"
	"account-service/internal/keycloak"
	"account-service/internal/security"
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
			baseURL:             env.GetString("KEYCLOAK_BASE_URL", "http://localhost:8081"),
			realm:               env.GetString("KEYCLOAK_REALM", "backend-auth-dev"),
			adminRealm:          env.GetString("KEYCLOAK_ADMIN_REALM", "master"),
			adminClientID:       env.GetString("KEYCLOAK_ADMIN_CLIENT_ID", "admin-cli"),
			adminClientSecret:   env.GetString("KEYCLOAK_ADMIN_CLIENT_SECRET", ""),
			jwkSetURL:           env.GetString("KEYCLOAK_JWK_SET_URL", "http://localhost:8081/realms/backend-auth-dev/protocol/openid-connect/certs"),
			expectedIssuer:      env.GetString("KEYCLOAK_EXPECTED_ISSUER", ""),
			expectedAudience:    env.GetString("KEYCLOAK_EXPECTED_AUDIENCE", ""),
			httpTimeout:         10 * time.Second,
			frontendLogoutURL:   env.GetString("KEYCLOAK_FRONTEND_LOGOUT_URL", ""),
			frontendClientID:    env.GetString("KEYCLOAK_FRONTEND_CLIENT_ID", "frontend-admin-auth-client"),
			frontendRedirectURL: env.GetString("KEYCLOAK_FRONTEND_REDIRECT_URL", ""),
		},
	}

	securityConfig := &security.AuthConfig{
		JWKSetURL:        cfg.keycloakConfig.jwkSetURL,
		ExpectedIssuer:   cfg.keycloakConfig.expectedIssuer,
		ExpectedAudience: cfg.keycloakConfig.expectedAudience,
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
	app := newApp(cfg, securityConfig, accountService)
	app.Run()
}
