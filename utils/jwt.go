package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomTokenClaim struct {
	UserID   uuid.UUID `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

const (
	OneDay   = 24 * time.Hour
	OneMonth = 24 * 30 * time.Hour
)

func GenerateNewJwtToken(
	userId uuid.UUID,
	username string,
	email string,
	expiryDate time.Time,
	signingKey string,
) (string, error) {
	tokenClaim := CustomTokenClaim{
		userId,
		username,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryDate),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)

	token, err := tokenObj.SignedString([]byte(signingKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}

func GenerateJwtTokenBasedOnExistingClaim(claim CustomTokenClaim, signingKey string) (string, error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := tokenObj.SignedString([]byte(signingKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}
