package httpserver

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver/backofficeuserhandler"
	"gocasts/gameapp/delivery/httpserver/matchinghandler"
	"gocasts/gameapp/delivery/httpserver/userhandler"
	"gocasts/gameapp/service/authorizationservice"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/backofficeuserservice"
	"gocasts/gameapp/service/matchingservice"
	"gocasts/gameapp/service/userservice"
	"gocasts/gameapp/validator/matchingvalidator"
	"gocasts/gameapp/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Router                *echo.Echo
	config                config.Config
	userhandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchinghandler       matchinghandler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service, matchingValidator matchingvalidator.Validator) Server {
	return Server{
		Router:                echo.New(),
		config:                config,
		userhandler:           userhandler.New(authSvc, userSvc, userValidator, config.Auth),
		backofficeUserHandler: backofficeuserhandler.New(authSvc, config.Auth, backofficeUserSvc, authorizationSvc),
		matchinghandler:       matchinghandler.New(config.Auth, authSvc, matchingValidator, matchingSvc),
	}
}

func (s Server) Serve() {
	// Middleware
	s.Router.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	s.Router.Use(middleware.Recover())       // recover panics as errors for proper error handling

	// Routes
	s.Router.GET("/health-check", s.healthCheck)

	s.userhandler.SetRouts(s.Router)

	s.backofficeUserHandler.SetRouts(s.Router)

	s.matchinghandler.SetRouts(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	s.Router.Logger.Fatal(s.Router.Start(address))
}
