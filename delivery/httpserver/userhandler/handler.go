package userhandler

import (
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/presenceservice"
	"gocasts/gameapp/service/userservice"
	"gocasts/gameapp/validator/uservalidator"
)

type Handler struct {
	authConfig    authservice.Config
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	presenceSvc   presenceservice.Service
}

func New(authSvc authservice.Service,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	authConfig authservice.Config, presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:    authConfig,
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
		presenceSvc:   presenceSvc,
	}
}
