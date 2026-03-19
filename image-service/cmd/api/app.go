package main

import (
	"image-service/internal/config/cache"
	"image-service/internal/config/messaging"
	"image-service/internal/config/security"
	"image-service/internal/controller"
	"image-service/internal/model"
	"image-service/internal/repository"
	"image-service/internal/service"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("albumtype", model.ValidateAlbumType)
}

type app struct {
	config         config
	securityConfig *security.AuthConfig
	Controller     *controller.Controller
	Service        *service.Service
	Storage        *repository.Storage
	Redis          *cache.RedisClient
}

type config struct {
	addr          string
	port          int
	storagePath   string
	apiBasePath   string
	redisConfig   *redisConfig
	apiGatewayURL string
	kafkaConfig   *kafkaConfig
}

type redisConfig struct {
	addr     string
	password string
	db       int
	duration time.Duration
}

type kafkaConfig struct {
	brokerList []string
	topic      string
}

func newApp(config config, securityConfig *security.AuthConfig) *app {
	redisClient := cache.NewRedisClient(
		config.redisConfig.addr,
		config.redisConfig.password,
		config.redisConfig.db,
		config.redisConfig.duration,
	)
	producer := messaging.NewKafkaProducer(config.kafkaConfig.brokerList, config.kafkaConfig.topic)
	storage := repository.NewStorage(config.storagePath, config.apiBasePath)
	service := service.NewService(storage, redisClient, securityConfig, validate)
	controller := controller.NewController(service, producer)

	return &app{
		config:         config,
		securityConfig: securityConfig,
		Controller:     controller,
		Service:        service,
		Storage:        storage,
		Redis:          redisClient,
	}
}

func (a *app) Run() {
	router := a.GetRouter()
	router.Run(a.config.addr)
}
