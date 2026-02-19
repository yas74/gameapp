package httpserver

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver/userhandler"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/userservice"
	"gocasts/gameapp/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config      config.Config
	userhandler userhandler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:      config,
		userhandler: userhandler.New(authSvc, userSvc, userValidator, config.Auth),
	}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger()) // use the RequestLogger middleware with slog logger
	e.Use(middleware.Recover())       // recover panics as errors for proper error handling

	// Routes
	e.GET("/health-check", s.healthCheck)

	s.userhandler.SetUserRouts(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
