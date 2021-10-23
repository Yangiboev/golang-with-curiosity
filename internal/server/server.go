package server

import (
	"context"

	"github.com/Yangiboev/golang-with-curiosity/config"
	"github.com/Yangiboev/golang-with-curiosity/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	certFile        = "ssl/server.crt"
	keyFile         = "ssl/server.pem"
	maxHeaderBytes  = 1 << 20 // 1mb
	gzipLevel       = 5
	stackSize       = 1 << 10 //1kb
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
	kafkaGroupID    = "storage_group"
)

type ServerOptions struct {
	Log     logger.Logger
	Config  config.Config
	Tracer  opentracing.Tracer
	MongoDB *mongo.Client
	Gin     *gin.Engine
	Redis   *redis.Client
}
type server struct {
	log     logger.Logger
	cfg     config.Config
	tracer  opentracing.Tracer
	mongoDB *mongo.Client
	gin     *gin.Engine
	redis   *redis.Client
}

func NewServer(opts *ServerOptions) *server {
	return &server{
		log:     opts.Log,
		cfg:     opts.Config,
		tracer:  opts.Tracer,
		mongoDB: opts.MongoDB,
		gin:     gin.New(),
		redis:   opts.Redis,
	}
}

// Run Start server
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	validate := validator.New()
	storageProducer := kafka.NewPro
}
