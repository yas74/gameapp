package main

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver"
	"gocasts/gameapp/repository/migrator"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/repository/mysql/mysqlaccesscontrol"
	"gocasts/gameapp/repository/mysql/mysqluser"
	"gocasts/gameapp/service/authorizationservice"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/backofficeuserservice"
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

	authSvc, userSvc, userValidator, backofficeUserSvc, authorizationSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeUserSvc, authorizationSvc)

	fmt.Println("start echo server")
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(mysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc

}
