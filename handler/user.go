package handler

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
)

type UserHandler struct {
	DB db.Querier
}

func NewUserHandler(DBTX *db.Queries) *UserHandler {
	return &UserHandler{
    DB: DBTX,
	}
}

func (u *UserHandler) CreateUser(ctx context.Context, newUser db.CreateUserParams) (db.User, error) {
	row, err := u.DB.CreateUser(ctx, newUser)
	if err != nil {
		return db.User{}, err
	}

	return row, nil
}
