package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RefreshTokenHandler struct {
	DB       db.Querier
	TokenGen utils.TokenGenerator
	Env      config.Config
}

func NewRefreshTokenHandler(
	DBTX *db.Queries,
	env config.Config,
	tokenGen utils.TokenGenerator,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		DB:       DBTX,
		Env:      env,
		TokenGen: tokenGen,
	}
}

type TokenHandlerInput struct {
	UserId   uuid.UUID
	Username string
	Email    string
}

func (h *RefreshTokenHandler) CreateTokenAndReturnUser(ctx context.Context, input TokenHandlerInput) (db.User, error) {
	duration := time.Now().Add(utils.OneMonth)
	tokeGenInput := utils.TokenGenInput{
		UserId:        input.UserId,
		Username:      input.Username,
		Email:         input.Email,
		ExpiryDate:    duration,
		SigningKey:    h.Env.JWT.RefreshToken,
		IsAccessToken: false,
	}
	signedRefreshToken, err := h.TokenGen.GenerateNewJwtToken(tokeGenInput)
	if err != nil {
		return db.User{}, err
	}

	createRefreshToken := db.CreateTokenParams{
		Token:      signedRefreshToken,
		ExpiryDate: duration,
		UserID:     input.UserId,
	}

	user, err := h.DB.CreateToken(ctx, createRefreshToken)
	if err != nil {
		return db.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}

func (h *RefreshTokenHandler) GetTokenByUserId(ctx context.Context, userId uuid.UUID) (db.RefreshToken, error) {
	token, err := h.DB.GetTokenByUserId(ctx, userId)
	if err != nil {
		return db.RefreshToken{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return token, nil
}

func (h *RefreshTokenHandler) UpdateTokenByUserId(ctx context.Context, input TokenHandlerInput) (db.RefreshToken, error) {
	duration := time.Now().Add(utils.OneMonth)
	tokeGenInput := utils.TokenGenInput{
		UserId:        input.UserId,
		Username:      input.Username,
		Email:         input.Email,
		ExpiryDate:    duration,
		SigningKey:    h.Env.JWT.RefreshToken,
		IsAccessToken: false,
	}
	signedRefreshToken, err := h.TokenGen.GenerateNewJwtToken(tokeGenInput)
	if err != nil {
		return db.RefreshToken{}, err
	}

	updatedRefreshToken := db.UpdateTokenByUserIdParams{
		Token:      signedRefreshToken,
		ExpiryDate: duration,
		UserID:     input.UserId,
	}

	token, err := h.DB.UpdateTokenByUserId(ctx, updatedRefreshToken)
	if err != nil {
		return db.RefreshToken{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return token, nil
}

func (h *RefreshTokenHandler) DeleteTokenByUserId(ctx context.Context, userId uuid.UUID) error {
	err := h.DB.DeleteTokenByUserId(ctx, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
