package main

import (
	"llm-service/internal/cache"
	"llm-service/internal/ratelimiter"
	"time"

	_ "llm-service/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

type app struct {
	config *config
	redis  *cache.RedisClient
}

type config struct {
	appConfig         *appConfig
	ratelimiterConfig *ratelimiterConfig
	servicesConfig    *servicesConfig
	summarizerConfig  *summarizerConfig
	redisConfig       *redisConfig
}

type appConfig struct {
	port     string
	endpoint string
}

type ratelimiterConfig struct {
	rate     rate.Limit
	burst    int
	interval time.Duration
}

type servicesConfig struct {
	WORK_EXPERIENCE_URL string
	EDUCATION_URL       string
	PROJECTS_URL        string
	SKILLS_URL          string
	CERTIFICATES_URL    string
}

type summarizerConfig struct {
	API_KEY string
	API_URL string
}

type redisConfig struct {
	addr     string
	duration time.Duration
}

func NewApp(config *config, rdb *cache.RedisClient) *app {
	return &app{
		config: config,
		redis:  rdb,
	}
}
func (app *app) getRouter() *gin.Engine {
	r := gin.Default()

	rl := ratelimiter.NewRateLimiter(
		app.config.ratelimiterConfig.rate,
		app.config.ratelimiterConfig.burst,
		app.config.ratelimiterConfig.interval,
	)
	r.Use(rl.Middleware())

	api := r.Group(app.config.appConfig.endpoint)

	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.GET("/health", app.handleGetHealth)
	api.GET("/summarize", app.handleGetSummarizer)
	return r
}

func (a *app) run() error {
	router := a.getRouter()
	return router.Run(":" + a.config.appConfig.port)
}
