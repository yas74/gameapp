package matchinghandler

import (
	"gocasts/gameapp/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRouts(e *echo.Echo) {
	matchingGroup := e.Group("/matching")

	matchingGroup.POST("/add-to-waiting-list", h.AddToWaitingList, middleware.Auth(h.authSvc, h.authConfig))
}
