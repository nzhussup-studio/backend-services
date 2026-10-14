package main

import (
	"account-service/internal/controller"
	"account-service/internal/service"
	"time"
)

type app struct {
	config     config
	controller *controller.AccountController
}

type config struct {
	addr           string
	port           int
	apiBasePath    string
	keycloakConfig *keycloakConfig
}

type keycloakConfig struct {
	baseURL           string
	realm             string
	adminRealm        string
	adminClientID     string
	adminClientSecret string
	httpTimeout       time.Duration
}

func newApp(config config, accountService *service.AccountService) *app {
	return &app{
		config:     config,
		controller: controller.NewAccountController(accountService),
	}
}
