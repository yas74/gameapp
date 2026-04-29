package config

import (
	"gocasts/gameapp/adapter/redis"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/matchingservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application     Application            `koanf:"aplication"`
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           mysql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matchingservice"`
	Redis           redis.Config           `koanf:"redis"`
}
