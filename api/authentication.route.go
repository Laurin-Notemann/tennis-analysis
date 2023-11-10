package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
)

type (
	RefreshReq struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	RegisterInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}

	AuthenticationRouter struct {
		UserHandler handler.UserHandler
	}

	OK  = string
	ERR = string

	Response struct {
		Status string
		Res    interface{}
	}

	ResponsePayload struct {
		AccessToken string  `json:"accessToken"`
		User        db.User `json:"user"`
	}
)

const (
	oneDay   = 24 * time.Hour
	oneMonth = 24 * 30 * time.Hour
)

func newAuthRouter(h handler.UserHandler) *AuthenticationRouter {
	return &AuthenticationRouter{UserHandler: h}
}

func (r AuthenticationRouter) register(ctx echo.Context) (err error) {
	newUser := new(RegisterInput)
	if err = ctx.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	signedAccessToken, err := validateUserInputAndGetJwt(
		*newUser,
		r.UserHandler.Env.JWT.AccessToken,
	)
	if err != nil {
		return err
	}

	input := handler.CreateUserInput{
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
	}
	user, err := r.UserHandler.CreateUser(ctx.Request().Context(), input)
	if err != nil {
		return err
	}

	registerPayload := ResponsePayload{
		AccessToken: signedAccessToken,
		User:        user,
	}
	return ctx.JSON(http.StatusCreated, registerPayload)
}

func (r AuthenticationRouter) refresh(ctx echo.Context) (err error) {
	req := new(RefreshReq)
	if err = ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	refreshToken := req.RefreshToken

	validRefreshToken, err := jwt.ParseWithClaims(refreshToken, &utils.CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.UserHandler.Env.JWT.RefreshToken), nil
	})

	refreshClaim, ok := validRefreshToken.Claims.(*utils.CustomTokenClaim)
	if ok && validRefreshToken.Valid {
		log.Println("User still logged in")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, err := r.UserHandler.GetUserByEmail(ctx.Request().Context(), refreshClaim.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(oneMonth))

	signedRefreshToken, err := utils.GenerateJwtTokenBasedOnExistingClaim(*refreshClaim, r.UserHandler.Env.JWT.RefreshToken)
	if err != nil {
		return err
	}

	accessToken := req.AccessToken

	validAccessToken, err := jwt.ParseWithClaims(accessToken, &utils.CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.UserHandler.Env.JWT.AccessToken), nil
	})

	_, okAcc := validAccessToken.Claims.(*utils.CustomTokenClaim)
	if !okAcc && !validAccessToken.Valid {
		accessToken, err = utils.GenerateNewJwtToken(user.Username, user.Email, oneDay, r.UserHandler.Env.JWT.AccessToken)
		if err != nil {
			return err
		}
	}

	params := db.UpdateUserByIdParams{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RefreshToken: sql.NullString{
			String: signedRefreshToken,
			Valid:  true,
		},
	}

	updatedUser, err := r.UserHandler.UpdateUserById(ctx.Request().Context(), params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	payload := ResponsePayload{
		AccessToken: accessToken,
		User:        updatedUser,
	}

	return ctx.JSON(http.StatusOK, payload)
}

func RegisterAuthRoute(baseUrl string, e *echo.Echo, r AuthenticationRouter) {

	e.POST(baseUrl+"/register", r.register)
	e.POST(baseUrl+"/refresh", r.refresh)
}

func validateUserInputAndGetJwt(input RegisterInput, token string) (string, error) {
	if input.Username == "" || input.Email == "" || input.Confirm == "" || input.Password == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Missing inputs.")
	}

	if input.Password != input.Confirm {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Password and confirmation do not match.")
	}

	signedAccessToken, err := utils.GenerateNewJwtToken(input.Username, input.Email, oneDay, token)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}
