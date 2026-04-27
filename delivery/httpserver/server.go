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
		config:                config,
		userhandler:           userhandler.New(authSvc, userSvc, userValidator, config.Auth),
		backofficeUserHandler: backofficeuserhandler.New(authSvc, config.Auth, backofficeUserSvc, authorizationSvc),
		matchinghandler:       matchinghandler.New(config.Auth, authSvc, matchingValidator, matchingSvc),
	}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	e.Use(middleware.Recover())       // recover panics as errors for proper error handling

	// Routes
	e.GET("/health-check", s.healthCheck)

	s.userhandler.SetRouts(e)

	s.backofficeUserHandler.SetRouts(e)

	s.matchinghandler.SetRouts(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
