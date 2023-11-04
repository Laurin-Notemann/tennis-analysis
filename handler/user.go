package handler

import (
	"context"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/gofrs/uuid"
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
	user, err := u.DB.CreateUser(ctx, newUser)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (u *UserHandler) GetAllUsers(ctx context.Context) ([]db.User, error) {
	users, err := u.DB.GetAllUsers(ctx)
	if err != nil {
		return []db.User{}, err
	}
	return users, nil
}

func (u *UserHandler) GetUserById(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := u.DB.GetUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (u *UserHandler) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := u.DB.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (u *UserHandler) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := u.DB.GetUserByUsername(ctx, username)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (u *UserHandler) UpdateUserById(ctx context.Context, args db.UpdateUserByIdParams) (db.User, error) {
	user, err := u.DB.UpdateUserById(ctx, args)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (u *UserHandler) DeleteUserById(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := u.DB.DeleteUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
