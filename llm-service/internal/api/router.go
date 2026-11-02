package api

import (
	"llm-service/internal/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) registerRoutes(r *gin.Engine) {
	api := r.Group(s.config.App.Endpoint)
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.POST("/summaries", s.handler.GetSummarizer)
	r.GET("/healthz", s.handler.GetHealth)
	configuration := r.Group("/v1/admin/llm/configuration", middleware.RequireGatewayIdentity())
	configuration.GET("", s.handler.GetConfiguration)
	configuration.PUT("", s.handler.PutConfiguration)
}
