package matchinghandler

import (
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/matchingservice"
	"gocasts/gameapp/service/presenceservice"
	"gocasts/gameapp/validator/matchingvalidator"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	matchingValidator matchingvalidator.Validator
	matchingSvc       matchingservice.Service
	presenceSvc       presenceservice.Service
}

func New(authConfig authservice.Config, authSvc authservice.Service,
	matchingValidator matchingvalidator.Validator,
	matchingSvc matchingservice.Service, presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		matchingValidator: matchingValidator,
		matchingSvc:       matchingSvc,
		presenceSvc:       presenceSvc,
	}
}
