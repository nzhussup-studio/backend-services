package api

import (
	"context"
	"errors"
	"llm-service/configs"
	"llm-service/internal/api/handlers"
	"llm-service/internal/api/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "llm-service/docs"

	"github.com/gin-gonic/gin"
)

type Cache interface {
	Set(key string, value any) error
	Get(key string, dest any) error
	Del(key string) error
	Ping() error
}

type Server struct {
	config  *configs.Config
	redis   Cache
	handler *handlers.Handler
}

func NewServer(config *configs.Config, rdb Cache) *Server {
	return &Server{
		config:  config,
		redis:   rdb,
		handler: handlers.New(config, rdb),
	}
}

func (s *Server) Router() *gin.Engine {
	r := gin.Default()
	rl := middleware.NewRateLimiter(
		s.config.RateLimiter.Rate,
		s.config.RateLimiter.Burst,
		s.config.RateLimiter.Interval,
	)
	r.Use(rl.Middleware())
	s.registerRoutes(r)
	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:    ":" + s.config.App.Port,
		Handler: s.Router(),
	}

	done := make(chan struct{})
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), s.config.App.ShutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(ctx)
		close(done)
	}()

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-done
	return nil
}
