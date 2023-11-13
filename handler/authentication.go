package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var AccessTokenInvalid = errors.New("Access token is invalid")

type (
	RefreshReq struct {
		AccessToken string `json:"accessToken"`
	}
	RegisterInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}

	LoginInput struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
		Password        string `json:"password"`
	}

	ResponsePayload struct {
		AccessToken string  `json:"accessToken"`
		User        db.User `json:"user"`
	}
)

type AuthenticationHandler struct {
	DB  db.Querier
	UserHandler  UserHandler
	TokenHandler RefreshTokenHandler
	TokenGen     utils.TokenGenerator
}

func NewAuthenticationHandler(
	DBTX *db.Queries,
  userHandler UserHandler,
  tokenHandler RefreshTokenHandler,
	tokenGen utils.TokenGenerator,
) *AuthenticationHandler {
	return &AuthenticationHandler{
		DB:       DBTX,
    UserHandler: userHandler,
    TokenHandler: tokenHandler,
		TokenGen: tokenGen,
	}
}

func (r *AuthenticationHandler) ValidateUserInputAndGetJwt(input RegisterInput) error {
	if input.Username == "" || input.Email == "" || input.Confirm == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing inputs.")
	}
	if input.Password != input.Confirm {
		return echo.NewHTTPError(http.StatusBadRequest, "Password and confirmation do not match.")
	}
	return nil
}

func (r *AuthenticationHandler) ValidateAccessToken(accessToken string, valid *jwt.Token, user db.User) (ResponsePayload, error) {
	if !valid.Valid {
		return ResponsePayload{}, AccessTokenInvalid
	}
	payload := ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}
	return payload, nil
}

func (r *AuthenticationHandler) ParseTokenGetUser(accessToken string, ctx echo.Context) (db.User, *jwt.Token, error) {
	if accessToken == "" {
		return db.User{}, nil, echo.NewHTTPError(http.StatusUnauthorized, "Access Token is empty")
	}
	validAccessToken, err := jwt.ParseWithClaims(accessToken, &utils.CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.UserHandler.Env.JWT.AccessToken), nil
	})

	accessTokenClaim, okAcc := validAccessToken.Claims.(*utils.CustomTokenClaim)
	if !okAcc {
		return db.User{}, validAccessToken, echo.NewHTTPError(http.StatusInternalServerError, "Couldn't parse claim")
	}

	user, err := r.UserHandler.GetUserById(ctx.Request().Context(), accessTokenClaim.UserID)
	if err != nil {
		return db.User{}, validAccessToken, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return user, validAccessToken, err
}

func (r *AuthenticationHandler) ValidateRefreshToken(ctx echo.Context, user db.User) error {
	refreshTokenObj, err := r.TokenHandler.GetTokenByUserId(ctx.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	refreshToken := refreshTokenObj.Token
	validRefreshToken, err := jwt.ParseWithClaims(refreshToken, &utils.CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.UserHandler.Env.JWT.RefreshToken), nil
	})
	_, ok := validRefreshToken.Claims.(*utils.CustomTokenClaim)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Couldn't parse claim.")
	}

	if validRefreshToken.Valid {
		return nil
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (r *AuthenticationHandler) CreateUserAndToken(ctx echo.Context, registerInput RegisterInput) (db.User, error) {
	userInput := CreateUserInput{
		Username: registerInput.Username,
		Email:    registerInput.Email,
		Password: registerInput.Password,
	}
	newUser, err := r.UserHandler.CreateUser(ctx.Request().Context(), userInput)
	if err != nil {
		return db.User{}, err
	}

	user, err := r.GenRefreshToken(ctx, &newUser)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (r *AuthenticationHandler) GenAccessToken(user *db.User) (string, error) {
	expiryDate := time.Now().Add(utils.OneDay)
	tokeGenInput := utils.TokenGenInput{
		UserId:        user.ID,
		Username:      user.Username,
		Email:         user.Email,
		ExpiryDate:    expiryDate,
		SigningKey:    r.UserHandler.Env.JWT.AccessToken,
		IsAccessToken: true,
	}
	signedAccessToken, err := r.TokenGen.GenerateNewJwtToken(tokeGenInput)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func (r *AuthenticationHandler) GenRefreshToken(ctx echo.Context, user *db.User) (db.User, error) {
	updatedUser, err := r.TokenHandler.CreateTokenAndReturnUser(
		ctx.Request().Context(),
		TokenHandlerInput{
			UserId:   user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	)
	if err != nil {
		return db.User{}, err
	}
	return updatedUser, nil
}

func (r *AuthenticationHandler) LoginRefreshToken(ctx echo.Context, user *db.User) (db.User, error) {
	err := r.TokenHandler.DeleteTokenByUserId(ctx.Request().Context(), user.ID)
	if err != nil {
		return db.User{}, err
	}

	refUser, err := r.GenRefreshToken(ctx, user)
	return refUser, err
}
