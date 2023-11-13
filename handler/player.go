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
