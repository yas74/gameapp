package middleware

import (
	"gocasts/gameapp/dto"
	"gocasts/gameapp/pkg/claim"
	"gocasts/gameapp/pkg/errmsg"
	"gocasts/gameapp/pkg/timestamp"
	"gocasts/gameapp/service/presenceservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpsertPresence(service presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := claim.GetClaimsFromEchoContext(c)

			req := dto.UpsertPresenceRequest{
				TimeStamp: timestamp.Now(),
				UserID:    claims.UserID,
			}

			_, err := service.Upsert(c.Request().Context(), req)
			if err != nil {
				// we can just log the error and continue
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorMsgSomethingWentWrong,
				})

			}

			return next(c)
		}
	}

}
