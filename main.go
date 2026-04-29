package main

import (
	"context"
	"fmt"
	"gocasts/gameapp/adapter/redis"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver"
	"gocasts/gameapp/repository/migrator"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/repository/mysql/mysqlaccesscontrol"
	"gocasts/gameapp/repository/mysql/mysqluser"
	"gocasts/gameapp/repository/redis/redismatching"
	"gocasts/gameapp/scheduler"
	"gocasts/gameapp/service/authorizationservice"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/backofficeuserservice"
	"gocasts/gameapp/service/matchingservice"
	"gocasts/gameapp/service/userservice"
	"gocasts/gameapp/validator/matchingvalidator"
	"gocasts/gameapp/validator/uservalidator"
	"os"
	"os/signal"
	"sync"
)

func main() {
	// TODO - read config path from command line
	cfg := config.Load("config.yml")

	// TODO - add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	// TODO - add struct and add these returned items as the struct fields
	authSvc, userSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingValidator, matchingSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingSvc, matchingValidator)

	go func() {
		server.Serve()
	}()

	done := make(chan bool)

	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(matchingSvc)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
	fmt.Println("arrived interupt signal, shutting down gracefully")

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()
	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("http server shutdown error:", err)
	}

	done <- true

	<-ctxWithTimeout.Done()

	wg.Wait()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service, matchingvalidator.Validator, matchingservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(mysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	matchingV := matchingvalidator.New()

	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingV, matchingSvc

}
