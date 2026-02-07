package main

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/userservice"
	"time"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	fmt.Println("start echo server")

	authConfig := authservice.Config{
		SignKey:               JwtSignKey,
		AccessExpirationTime:  AccessTokenExpireDuration,
		RefreshExpirationTime: RefreshTokenExpireDuration,
		AccessSubject:         AccessTokenSubject,
		RefreshSubject:        RefreshTokenSubject,
	}

	mysqlCfg := mysql.Config{
		Username: "gameapp",
		Password: "gameappt0lk2o20",
		Port:     3308,
		Host:     "localhost",
		DBName:   "gameapp_db",
	}

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth:       authConfig,
		Mysql:      mysqlCfg,
	}

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.New(authSvc, mysqlRepo)

	return authSvc, userSvc

}
