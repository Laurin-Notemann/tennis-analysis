package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	DB  db.Querier
	Env config.Config
}

func NewUserHandler(DBTX *db.Queries, env config.Config) *UserHandler {
	return &UserHandler{
		DB:  DBTX,
		Env: env,
	}
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

func (u *UserHandler) CreateUser(ctx context.Context, input CreateUserInput) (db.User, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	registeredUser := db.CreateUserParams{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPw),
	}

	user, err := u.DB.CreateUser(ctx, registeredUser)
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return db.User{}, echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	if err != nil {
		return db.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
