package main

import (
	"account-service/internal/controller"
	"account-service/internal/security"
	"account-service/internal/service"
	"time"
)

type app struct {
	config         config
	securityConfig *security.AuthConfig
	controller     *controller.AccountController
}

type config struct {
	addr           string
	port           int
	apiBasePath    string
	keycloakConfig *keycloakConfig
}

type keycloakConfig struct {
	baseURL             string
	realm               string
	adminRealm          string
	adminClientID       string
	adminClientSecret   string
	jwkSetURL           string
	expectedIssuer      string
	expectedAudience    string
	httpTimeout         time.Duration
	frontendLogoutURL   string
	frontendClientID    string
	frontendRedirectURL string
}

func newApp(config config, securityConfig *security.AuthConfig, accountService *service.AccountService) *app {
	return &app{
		config:         config,
		securityConfig: securityConfig,
		controller:     controller.NewAccountController(accountService),
	}
}
