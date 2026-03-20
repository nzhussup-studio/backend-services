package main

import (
	"account-service/internal/controller"
	"account-service/internal/security"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *app) GetRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", controller.HealthCheckHandler)

	v1 := router.Group(a.config.apiBasePath)
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.Use(security.AuthMiddleware(a.securityConfig))
	v1.DELETE("", a.controller.DeleteCurrentAccount)

	return router
}

func (a *app) Run() {
	router := a.GetRouter()
	_ = router.Run(a.config.addr)
}
