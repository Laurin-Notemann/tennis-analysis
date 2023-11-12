package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TokenGenInput struct {
	UserId        uuid.UUID
	Username      string
	Email         string
	ExpiryDate    time.Time
	SigningKey    string
	IsAccessToken bool
}

type TokenGenerator interface {
	GenerateNewJwtToken(
		input TokenGenInput,
	) (string, error)
}

type CustomTokenClaim struct {
	UserID   uuid.UUID `json:"userId"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type ProdTokenGenerator struct{}

const (
	OneDay   = 24 * time.Hour
	OneMonth = 24 * 30 * time.Hour
)

func (t *ProdTokenGenerator) GenerateNewJwtToken(
	input TokenGenInput,
) (string, error) {
	tokenClaim := CustomTokenClaim{
		input.UserId,
		input.Username,
		input.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(input.ExpiryDate),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)

	token, err := tokenObj.SignedString([]byte(input.SigningKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}
