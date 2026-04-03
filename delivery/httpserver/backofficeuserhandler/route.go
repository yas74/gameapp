package backofficeuserhandler

import (
	"gocasts/gameapp/delivery/httpserver/middleware"
	"gocasts/gameapp/entity"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRouts(e *echo.Echo) {
	userGroup := e.Group("/backoffice/users")

	userGroup.GET("/", h.listUsers, middleware.Auth(h.authSvc, h.authConfig),
		middleware.AccessCheck(h.authorizationSvc, entity.UserListPermission))
}
