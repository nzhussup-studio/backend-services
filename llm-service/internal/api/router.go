package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) registerRoutes(r *gin.Engine) {
	api := r.Group(s.config.App.Endpoint)
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.GET("/health", s.handler.GetHealth)
	api.GET("/summarize", s.handler.GetSummarizer)
}
