package api

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApi(ctx context.Context, resource handler.ResourceHandlers, tokenGen utils.TokenGenerator) *echo.Echo {
	baseUrl := "/api"
	authRouter := newAuthRouter(resource.UserHandler, resource.TokenHandler, tokenGen, resource.AuthHandler)
	userRouter := newUserRouter(resource.UserHandler)
	playerRouter := newPlayerRouter(resource.PlayerHandler, resource.TeamHandler, resource.UserHandler)

	customMiddleware := NewMiddleware(resource.AuthHandler)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	RegisterAuthRoute(baseUrl, e, *authRouter)

	RegisterUserRoute(baseUrl, e, *userRouter, *customMiddleware)
	RegisterPlayersRoute(baseUrl, e, *playerRouter, *customMiddleware)
	RegisterHtmlPageRoutes(e)

	return e
}
