package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	AuthHandler handler.AuthenticationHandler
}

func NewMiddleware(
	a handler.AuthenticationHandler,
) *Middleware {
	return &Middleware{
		AuthHandler: a,
	}
}

func (m *Middleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		headers := ctx.Request().Header.Get("Authorization")
		if headers == "" {
			ctx.Error(echo.NewHTTPError(http.StatusUnauthorized, "No bearer of token"))
			return nil
		}

		token := strings.Split(headers, " ")

		user, validToken, err := m.AuthHandler.ParseTokenGetUser(token[1], ctx)
		if err != nil {
			ctx.Error(echo.NewHTTPError(http.StatusUnauthorized, err.Error()))
			return nil
		}

		_, err = m.AuthHandler.ValidateAccessToken(token[1], validToken, user)
		if errors.Is(err, handler.AccessTokenInvalid) {
			ctx.Error(echo.NewHTTPError(http.StatusUnauthorized, err.Error()))
			return nil
		}

		return next(ctx)
	}
}
