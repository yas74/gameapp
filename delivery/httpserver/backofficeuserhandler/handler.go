package backofficeuserhandler

import (
	"gocasts/gameapp/service/authorizationservice"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/backofficeuserservice"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	authorizationSvc  authorizationservice.Service
	backofficeUserSvc backofficeuserservice.Service
}

func New(authSvc authservice.Service,
	authConfig authservice.Config,
	backofficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		backofficeUserSvc: backofficeUserSvc,
		authorizationSvc:  authorizationSvc,
	}
}
