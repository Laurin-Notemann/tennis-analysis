package api

import (
	"database/sql"
	"net/http"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PlayerRouter struct {
	PlayerHandler handler.PlayerHandler
	TeamHandler   handler.TeamHandler
	UserHandler   handler.UserHandler
}

func newPlayerRouter(
	p handler.PlayerHandler,
	t handler.TeamHandler,
	u handler.UserHandler,
) *PlayerRouter {
	return &PlayerRouter{PlayerHandler: p, TeamHandler: t, UserHandler: u}
}

type CreatePlayerRequest struct {
	FirstName string
	LastName  string
	UserId    uuid.UUID
}

func (r *PlayerRouter) CreatePlayer(ctx echo.Context) (err error) {
	request := new(CreatePlayerRequest)
	if err = ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	teamParams := db.CreateNewTeamWithOnePlayerParams{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Name: sql.NullString{
			String: request.FirstName + " " + request.LastName,
			Valid:  true,
		},
		UserID: request.UserId,
	}

	team, err := r.TeamHandler.CreateTeamWithOnePlayer(ctx.Request().Context(), teamParams)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	player, err := r.PlayerHandler.GetPlayerById(ctx.Request().Context(), team.PlayerOne)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, player)
}

func RegisterPlayersRoute(baseUrl string, e *echo.Echo, r PlayerRouter, middleware Middleware) {
	e.POST(baseUrl+"/players", r.CreatePlayer)
}
