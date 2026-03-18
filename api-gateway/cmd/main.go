package main

import (
	"api-gateway/internal/routes"
	"log/slog"
)

func main() {
	slog.Info("Starting API Gateway...", slog.Int64("port", 8082))

	r := routes.SetupRouter()
	r.Run(":8082")
}
