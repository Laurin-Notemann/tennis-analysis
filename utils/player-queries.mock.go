package utils

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
)

func (d *DBQueriesMock) GetPlayerById(ctx context.Context, id uuid.UUID) (db.Player, error) {
	return db.Player{}, nil
}

func (d *DBQueriesMock) DeletePlayerById(ctx context.Context, id uuid.UUID) (db.Player, error) {
	return db.Player{}, nil
}

func (d *DBQueriesMock) UpdatePlayerById(ctx context.Context, arg db.UpdatePlayerByIdParams) (db.Player, error) {
	return db.Player{}, nil
}
