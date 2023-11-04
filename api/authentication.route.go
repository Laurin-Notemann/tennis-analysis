package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Laurin-Notemann/tennis-analysis/handler"
)
type RegiserInput struct {
  Username string `json:"username"`
  Email string `json:"email"`
  Password string `json:"password"`
  Confirm string `json:"confirm"`
}

type authenticationRouter struct {
	UserHandler handler.UserHandler
}

func newAuthRouter(h handler.UserHandler) *authenticationRouter {
	return &authenticationRouter{UserHandler: h}
}

func (r authenticationRouter) register(ctx echo.Context) (err error) {
  input := new(RegiserInput)

  if err = ctx.Bind(input); err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, err.Error())
  }

  log.Printf("%v", input)

	//r.Handler.CreateUser(ctx.Request().Context(), ctx.Request().Body)

	return ctx.JSON(http.StatusOK, input)
}

func RegisterAuthRoute(baseUrl string, e *echo.Echo, r authenticationRouter) {

	e.POST(baseUrl + "/register", r.register)
}
