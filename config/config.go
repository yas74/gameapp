package config

import (
	"gocasts/gameapp/adapter/redis"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/scheduler"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/matchingservice"
	"gocasts/gameapp/service/presenceservice"
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
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
	PresenceService presenceservice.Config `koanf:"presence_service"`
	Scheduler       scheduler.Config       `koanf:"scheduler"`
}
