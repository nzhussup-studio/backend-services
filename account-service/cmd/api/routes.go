package main

import (
	"net/http"

	"account-service/internal/security"

	"github.com/gin-gonic/gin"
)

func (a *app) GetRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group(a.config.apiBasePath)
	v1.Use(security.AuthMiddleware(a.securityConfig))
	v1.DELETE("", a.controller.DeleteCurrentAccount)

	return router
}

func (a *app) Run() {
	router := a.GetRouter()
	_ = router.Run(a.config.addr)
}
