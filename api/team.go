package api

import (
	"net/http"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
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
	return nil
}

func (r *TeamRouter) DeleteTeamById(ctx echo.Context) (err error) {
	return nil
}

func (r *TeamRouter) UpdateTeamById(ctx echo.Context) (err error) {
	return nil
}

func RegisterTeamRoute(baseUrl string, e *echo.Echo, r TeamRouter, middleware Middleware) {
	e.POST(baseUrl+"/teams", r.CreateTeam, middleware.AuthMiddleware)
	e.GET(baseUrl+"/teams/:id", r.GetAllTeamsByUserId, middleware.AuthMiddleware)
	e.DELETE(baseUrl+"/teams/:id", r.DeleteTeamById, middleware.AuthMiddleware)
	e.PUT(baseUrl+"/teams", r.UpdateTeamById, middleware.AuthMiddleware)
}
