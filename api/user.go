package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
)

type UserRouter struct {
	UserHandler handler.UserHandler
}

type UserResponse struct {
	Status string  `json:"status"`
	Data   db.User `json:"data"`
}

func newUserRouter(h handler.UserHandler) *UserRouter {
	return &UserRouter{UserHandler: h}
}

func (r *UserRouter) getUserById(ctx echo.Context) error {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	user, err := r.UserHandler.GetUserById(ctx.Request().Context(), userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	res := UserResponse{
		Status: "success",
		Data:   user,
	}
	return ctx.JSON(http.StatusOK, res)
}

func RegisterUserRoute(baseUrl string, e *echo.Echo, r UserRouter, middleware Middleware) {

	e.GET(baseUrl+"/users/:id", r.getUserById, middleware.AuthMiddleware)
}
