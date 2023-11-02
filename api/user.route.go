package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Laurin-Notemann/tennis-analysis/handler"
)
type userRouter struct {
  Handler handler.UserHandler
}

func newUserRouter (h handler.UserHandler) *userRouter {
  return &userRouter{Handler: h}
}

func(r userRouter) CreateUser (ctx echo.Context) error{

  //r.Handler.CreateUser(ctx.Request().Context(), ctx.Request().Body)

  return ctx.JSON(http.StatusOK, "Hello World")
}

func RegisterUserRoute(e *echo.Echo, r userRouter) {

  e.POST("/users", r.CreateUser)
  e.GET("/users", func(c echo.Context) error {
    return c.JSON(http.StatusOK, "Hi")
  })
}


