package handler

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
)

type TeamHandler struct {
	DB db.Querier
}

func NewTeamHandler(DB *db.Queries) *TeamHandler {
	return &TeamHandler{
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

func (h *TeamHandler) CreateTeamWithTwoPlayers(ctx context.Context, args db.CreateTeamWithTwoPlayersParams) (db.Team, error) {
	team, err := h.DB.CreateTeamWithTwoPlayers(ctx, args)
	if err != nil {
		return db.Team{}, err
	}
	return team, nil
}

func (h *TeamHandler) GetAllTeamsByUserId(ctx context.Context, id uuid.UUID) ([]db.Team, error) {
	teams, err := h.DB.GetAllTeamsByUserId(ctx, id)
	if err != nil {
		return []db.Team{}, err
	}
	return teams, nil
}

func (h *TeamHandler) GetTeamById(ctx context.Context, id uuid.UUID) (db.Team, error) {
	team, err := h.DB.GetTeamById(ctx, id)
	if err != nil {
		return db.Team{}, err
	}
	return team, nil
}

func (h *TeamHandler) UpdateTeamById(ctx context.Context, args db.UpdateTeamByIdParams) (db.Team, error) {
	team, err := h.DB.UpdateTeamById(ctx, args)
	if err != nil {
		return db.Team{}, err
	}
	return team, nil
}

func (h *TeamHandler) DeleteTeamById(ctx context.Context, id uuid.UUID) (db.Team, error) {
	team, err := h.DB.DeleteTeamById(ctx, id)
	if err != nil {
		return db.Team{}, err
	}
	return team, nil
}
