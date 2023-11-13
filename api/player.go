package api

import (
	"github.com/Laurin-Notemann/tennis-analysis/handler"
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

func (r *PlayerRouter) CreatePlayer(ctx echo.Context) error {
	return nil
}

func RegisterPlayersRoute(baseUrl string, e *echo.Echo, r PlayerRouter, middleware Middleware) {
	e.POST(baseUrl+"/players", r.CreatePlayer)
}
