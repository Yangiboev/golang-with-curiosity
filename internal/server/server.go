package server

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/Yangiboev/golang-with-curiosity/config"
	"github.com/Yangiboev/golang-with-curiosity/internal/interceptors"
	"github.com/Yangiboev/golang-with-curiosity/internal/middlewares"
	productsHttpV1 "github.com/Yangiboev/golang-with-curiosity/internal/product/delivery/http/v1"
	"github.com/Yangiboev/golang-with-curiosity/internal/product/delivery/kafka"
	"github.com/Yangiboev/golang-with-curiosity/internal/product/repository"
	"github.com/Yangiboev/golang-with-curiosity/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	certFile        = "ssl/server.crt"
	keyFile         = "ssl/server.pem"
	maxHeaderBytes  = 1 << 20 // 1mb
	gzipLevel       = 5
	stackSize       = 1 << 10 //1kb
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
	kafkaGroupID    = "product_group"
)

type ServerOptions struct {
	Log     logger.Logger
	Config  config.Config
	Tracer  opentracing.Tracer
	MongoDB *mongo.Client
	Echo    *echo.Echo
	Redis   *redis.Client
}
type server struct {
	log     logger.Logger
	cfg     config.Config
	tracer  opentracing.Tracer
	mongoDB *mongo.Client
	echo    *echo.Echo
	redis   *redis.Client
}

func NewServer(opts *ServerOptions) *server {
	return &server{
		log:     opts.Log,
		cfg:     opts.Config,
		tracer:  opts.Tracer,
		mongoDB: opts.MongoDB,
		echo:    echo.New(),
		redis:   opts.Redis,
	}
}

// Run Start server
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	validate := validator.New()
	productsProducer := kafka.NewProductsProducer(s.log, s.cfg)
	productsProducer.Run()
	productMongoRepo := repository.NewProductMongoRepo(s.mongoDB)
	productRedisRepo := repository.NewProductRedisRepository(s.redis)
	productUC := usecase.NewProductUC(productMongoRepo, productRedisRepo, s.log, productsProducer)

	im := interceptors.NewInterceptorManager(s.log, s.cfg)
	mw := middlewares.NewMiddlewareManager(s.log, s.cfg)
	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		s.log.Fatalf("failed to load key pair: %s", err)
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
			Timeout:           s.cfg.Server.Timeout * time.Second,
			MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
			Time:              s.cfg.Server.Timeout * time.Minute,
		}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
			im.Logger,
		),
	)
	productService := product.NewProductService(s.log, productUC, validate)
	productsService.RegisterProductsServiceServer(grpcServer, productService)
	grpc_prometheus.Register(grpcServer)
	v1 := s.echo.Group("/api/v1")

	productHandlers := productsHttpV1.NewProductHandlers(s.log, productUC, validate, v1.Group("/products"), mw)
	productHandlers.MapRoutes()
	productCG := kafka.NewProductsConsumerGroup(s.cfg.Kafka.Brokers, kafkaGroupID, s.log, s.cfg, productUC, validate)
	productCG.RunConsumers(ctx, cancel)
	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.Http.Port)
		s.runHttpServer()
	}()
}
