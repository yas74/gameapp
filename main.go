package main

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver"
	"gocasts/gameapp/repository/migrator"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/userservice"
	"gocasts/gameapp/validator/uservalidator"
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

	// TODO - read config path from command line
	cfg2 := config.Load("config.yml")
	fmt.Printf("cfg2: %v\n", cfg2)
	// TODO - merge cfg with cfg2

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth:       authConfig,
		Mysql:      mysqlCfg,
	}
	// TODO - add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	authSvc, userSvc, userValidator := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	fmt.Println("start echo server")
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.New(authSvc, mysqlRepo)

	uV := uservalidator.New(mysqlRepo)

	return authSvc, userSvc, uV

}
