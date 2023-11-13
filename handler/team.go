package handler

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
)

type TeamHandler struct {
	DB db.Querier
}

func NewTeamHandler(DB *db.Queries) *PlayerHandler {
	return &PlayerHandler{
		DB: DB,
	}
}

func (h *TeamHandler) CreateTeamWithOnePlayer(ctx context.Context, args db.CreateNewTeamWithOnePlayerParams) (db.Team, error) {
	team, err := h.DB.CreateNewTeamWithOnePlayer(ctx, args)
	if err != nil {
		return db.Team{}, err
	}
	return team, nil
}
