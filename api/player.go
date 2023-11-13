package api

import "github.com/Laurin-Notemann/tennis-analysis/handler"

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

func (r *PlayerRouter) CreatePlayer() {

}
