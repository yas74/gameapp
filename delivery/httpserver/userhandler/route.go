package userhandler

import (
	"gocasts/gameapp/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetUserRouts(e *echo.Echo) {
	userGroup := e.Group("/users")

	userGroup.GET("/profile", h.userProfile, middleware.Auth(h.authSvc, h.authConfig))
	userGroup.POST("/login", h.userLogin)
	userGroup.POST("/register", h.userRegister)
}
