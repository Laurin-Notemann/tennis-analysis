package api

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApi(ctx context.Context, resource handler.ResourceHandlers) *echo.Echo {
  baseUrl := "/api"
  authRouter := newAuthRouter(resource.UserHandler, resource.TokenHandler)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	RegisterAuthRoute(baseUrl, e, *authRouter)
  RegisterHtmlPageRoutes(e)

	return e
}

type UserServerInterface interface {
	CreateUser(ctx echo.Context) error
}

type UserServerInterfaceWrapper struct {
	Handler UserServerInterface
}

func (w *UserServerInterfaceWrapper) CreateUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateUser(ctx)
	return err
}
