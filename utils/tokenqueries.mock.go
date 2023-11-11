package utils

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
)

func (d *DBQueriesMock) CreateToken(ctx context.Context, arg db.CreateTokenParams) (db.User, error) {
	return db.User{}, nil
}

func (d *DBQueriesMock) UpdateTokenByUserId(ctx context.Context, arg db.UpdateTokenByUserIdParams) (db.RefreshToken, error) {
	return db.RefreshToken{}, nil
}

func (d *DBQueriesMock) GetTokenByUserId(ctx context.Context, userId uuid.UUID) (db.RefreshToken, error) {
	return db.RefreshToken{}, nil
}

func (d *DBQueriesMock) DeleteTokenByUserId(ctx context.Context, id uuid.UUID) error {
	return nil
}
