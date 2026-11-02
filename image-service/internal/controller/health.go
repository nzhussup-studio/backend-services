package controller

import (
	"image-service/internal/json"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler godoc
// @Summary Health check endpoint
// @Description Returns 200 OK if the service is up
// @Tags Health
// @Produce json
// @Success 200 {object} model.SuccessResponse "Service is healthy"
func HealthCheckHandler(c *gin.Context) {
	json.ConstructJsonResponseSuccess(c, nil)
}
