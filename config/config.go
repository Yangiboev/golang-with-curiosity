package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	GRPC_PORT = "GRPC_PORT"
	HTTP_PORT = "HTTP_PORT"
)

type Config struct {
	AppVersion string
	Server     Server
	Logger     Logger
	Jaegar     Jaegar
	Metrics    Metrics
	MongoDB    MongoDB
	Kafka      Kafka
	Http       Http
	Redis      Redis
}

type Server struct {
	Port              string
	Development       bool
	Timeout           time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	MaxConnectionIdle time.Duration
	MaxConnectionAge  time.Duration
	Kafka             Kafka
}
type Http struct {
	Port              string
	PprofPort         string
	Timeout           time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	CookieLifeTime    int
	SessionCookieName string
}

// Logger config
type Logger struct {
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}
type Metrics struct {
	Port        string
	URL         string
	ServiceName string
}

//Jaegar config
type Jaegar struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

type MongoDB struct {
	URI      string
	User     string
	Password string
	DB       string
}

type Kafka struct {
	Brokers []string
}
type Redis struct {
	RedisAddress   string
	RedisPassword  string
	RedisDB        string
	RedisDefaultDB string
	Password       string
	MinIdleConn    int
	PoolSize       int
	PoolTimeout    int
	DB             int
}

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if os.Getenv("MODE") == "DOCKER" {
		viper.SetConfigName("config-docker.yml")
	} else {
		viper.SetConfigName("config.yaml")
	}
	err := viper.ReadInConfig()
	return err
}

func LoadConfig() (Config, error) {
	if err := exportConfig(); err != nil {
		return Config{}, err
	}
	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return c, err
	}
	grpcPort := os.Getenv(GRPC_PORT)
	if grpcPort != "" {
		c.Server.Port = grpcPort
	}
	httpPort := os.Getenv(HTTP_PORT)
	if grpcPort != "" {
		c.Http.Port = httpPort
	}
	return c, nil
}
