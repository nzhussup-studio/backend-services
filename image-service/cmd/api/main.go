package main

import (
	appconfig "image-service/internal/config"

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
	cfg := appconfig.Load()

	app := newApp(cfg)
	app.Redis.CheckHealth()

	app.Run()
}
