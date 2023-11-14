package api

import (
	"database/sql"
	"net/http"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
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

	if request.FirstName == "" || request.LastName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing first or last name")
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
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return echo.NewHTTPError(http.StatusConflict, "player already exists")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	player, err := r.PlayerHandler.GetPlayerById(ctx.Request().Context(), team.PlayerOne)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, player)
}

func (r *PlayerRouter) GetAllPlayersByUserId(ctx echo.Context) (err error) {
	param := ctx.Param("id")
	userId, err := uuid.Parse(param)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	teams, err := r.TeamHandler.DB.GetAllTeamsByUserId(ctx.Request().Context(), userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var allPlayerIds []uuid.UUID

	for _, team := range teams {
		if team.PlayerTwo == nil {
			allPlayerIds = append(allPlayerIds, team.PlayerOne)
		}
	}

	var allPlayer []db.Player

	for _, id := range allPlayerIds {
		player, err := r.PlayerHandler.GetPlayerById(ctx.Request().Context(), id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		allPlayer = append(allPlayer, player)
	}

	return ctx.JSON(http.StatusOK, allPlayer)
}

func RegisterPlayersRoute(baseUrl string, e *echo.Echo, r PlayerRouter, middleware Middleware) {
	e.POST(baseUrl+"/players", r.CreatePlayer, middleware.AuthMiddleware)
	e.GET(baseUrl+"/players/:id", r.GetAllPlayersByUserId, middleware.AuthMiddleware)
}

func filter() {

}
