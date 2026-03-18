package main

import "github.com/gin-gonic/gin"

// handleGetHealth godoc
// @Summary      Health check endpoint
// @Description  Checks the connectivity and health of dependent services, particularly Redis.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string  "Status OK"
// @Failure      500  {object}  map[string]string  "Redis connection failed"
// @Router       /v1/health [get]
func (app *app) handleGetHealth(ctx *gin.Context) {

	if err := app.redis.Ping(); err != nil {
		ctx.JSON(500, gin.H{"status": "error", "message": "Redis connection failed"})
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}
