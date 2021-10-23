package main

import (
	"context"
	"log"

	"github.com/Yangiboev/golang-with-curiosity/config"
	"github.com/Yangiboev/golang-with-curiosity/internal/server"
	"github.com/Yangiboev/golang-with-curiosity/pkg/jaegar"
	"github.com/Yangiboev/golang-with-curiosity/pkg/kafka"
	"github.com/Yangiboev/golang-with-curiosity/pkg/logger"
	"github.com/Yangiboev/golang-with-curiosity/pkg/mongodb"
	"github.com/Yangiboev/golang-with-curiosity/pkg/redis"
	"github.com/opentracing/opentracing-go"
)

// @title Ecommerce template microservice
// @version 1.0
// @description storage management
// @contact.name Dilmurod Yangiboev
// @contact.url https://github.com/Yangiboev
// @contact.email dimok.aka0771@gmail.com

// @host localhost:5000
// @BasePath /api/v1

func main() {
	log.Println("Starting storage project")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting storage server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)
	tracer, closer, err := jaegar.InitJaegar(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("opentracing connected")
	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, cfg)
	if err != nil {
		appLogger.Fatal("cannot create mongodb", err)
	}
	defer func() {
		if err := mongoDBConn.Disconnect(ctx); err != nil {
			appLogger.Fatal("mongoDBConn.Disconnect", err)
		}
	}()
	appLogger.Infof("MongoDB connected: %v", mongoDBConn.NumberSessionsInProgress())
	kafkaConn, err := kafka.NewKafkaConn(ctx, cfg)
	if err != nil {
		appLogger.Fatal("cannot create kafka", err)
	}
	defer kafkaConn.Close()
	brokers, err := kafkaConn.Brokers()
	if err != nil {
		appLogger.Fatal("NewKafkaConn", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)
	redisClient := redis.NewRedisClient(cfg)
	appLogger.Info("Redis connected")

	s := server.NewServer(&server.ServerOptions{
		Log:     appLogger,
		Config:  cfg,
		Tracer:  tracer,
		MongoDB: mongoDBConn,
		Redis:   redisClient,
	})
	appLogger.Fatal(s.Run())
}
