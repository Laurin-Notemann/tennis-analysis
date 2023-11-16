package api

import (
	"net/http"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TeamRouter struct {
	PlayerHandler handler.PlayerHandler
	TeamHandler   handler.TeamHandler
	UserHandler   handler.UserHandler
}

func newTeamRouter(
	p handler.PlayerHandler,
	t handler.TeamHandler,
	u handler.UserHandler,
) *TeamRouter {
	return &TeamRouter{PlayerHandler: p, TeamHandler: t, UserHandler: u}
}

func (r *TeamRouter) CreateTeam(ctx echo.Context) (err error) {
	request := new(db.CreateTeamWithTwoPlayersParams)
	if err = ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	team, err := r.TeamHandler.CreateTeamWithTwoPlayers(ctx.Request().Context(), *request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, team)
}

func (r *TeamRouter) GetAllTeamsByUserId(ctx echo.Context) (err error) {
	param := ctx.Param("userId")
	userId, err := uuid.Parse(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	teams, err := r.TeamHandler.GetAllTeamsByUserId(ctx.Request().Context(), userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, teams)
}

func (r *TeamRouter) DeleteTeamById(ctx echo.Context) (err error) {
	param := ctx.Param("id")
	id, err := uuid.Parse(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	team, err := r.TeamHandler.DeleteTeamById(ctx.Request().Context(), id)
	return ctx.JSON(http.StatusOK, team)
}

func (r *TeamRouter) UpdateTeamById(ctx echo.Context) (err error) {
	request := new(db.UpdateTeamByIdParams)
	if err = ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	team, err := r.TeamHandler.UpdateTeamById(ctx.Request().Context(), *request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, team)
}

func RegisterTeamRoute(baseUrl string, e *echo.Echo, r TeamRouter, middleware Middleware) {
	e.POST(baseUrl+"/teams", r.CreateTeam, middleware.AuthMiddleware)
	e.GET(baseUrl+"/teams/:userId", r.GetAllTeamsByUserId, middleware.AuthMiddleware)
	e.DELETE(baseUrl+"/teams/:id", r.DeleteTeamById, middleware.AuthMiddleware)
	e.PUT(baseUrl+"/teams", r.UpdateTeamById, middleware.AuthMiddleware)
}
