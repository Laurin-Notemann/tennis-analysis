package handler

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	DB db.Querier
}

func NewPlayerHandler(DB *db.Queries) *PlayerHandler {
	return &PlayerHandler{
		DB: DB,
	}
}

func (h *PlayerHandler) GetPlayerById(ctx context.Context, id uuid.UUID) (db.Player, error) {
	player, err := h.DB.GetPlayerById(ctx, id)
	if err != nil {
		return db.Player{}, err
	}
	return player, nil
}

func (h *PlayerHandler) DeletePlayerById(ctx context.Context, id uuid.UUID) (db.Player, error) {
	player, err := h.DB.DeletePlayerById(ctx, id)
	if err != nil {
		return db.Player{}, err
	}
	return player, nil
}

func (h *PlayerHandler) UpdatePlayerById(ctx context.Context, arg db.UpdatePlayerByIdParams) (db.Player, error) {
	player, err := h.DB.UpdatePlayerById(ctx, arg)
	if err != nil {
		return db.Player{}, err
	}
	return player, nil
}
