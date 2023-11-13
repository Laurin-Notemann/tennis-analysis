package utils

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
)

func (d *DBQueriesMock) CreateNewTeamWithOnePlayer(ctx context.Context, args db.CreateNewTeamWithOnePlayerParams) (db.Team, error) {
	return db.Team{}, nil
}

func (d *DBQueriesMock) CreateTeamWithTwoPlayers(ctx context.Context, arg db.CreateTeamWithTwoPlayersParams) (db.Team, error) {
	return db.Team{}, nil
}

func (d *DBQueriesMock) GetTeamById(ctx context.Context, id uuid.UUID) (db.Team, error) {
	return db.Team{}, nil
}

func (d *DBQueriesMock) DeleteTeamById(ctx context.Context, id uuid.UUID) (db.Team, error) {
	return db.Team{}, nil
}
