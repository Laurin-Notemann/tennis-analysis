package api

import (
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
		UserHandler  handler.UserHandler
		TokenHandler handler.RefreshTokenHandler
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

func newAuthRouter(h handler.UserHandler, t handler.RefreshTokenHandler) *AuthenticationRouter {
	return &AuthenticationRouter{UserHandler: h, TokenHandler: t}
}

func (r AuthenticationRouter) register(ctx echo.Context) (err error) {
	registerInput := new(RegisterInput)
	if err = ctx.Bind(registerInput); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = validateUserInputAndGetJwt(
		*registerInput,
	)
	if err != nil {
		return err
	}

	userInput := handler.CreateUserInput{
		Username: registerInput.Username,
		Email:    registerInput.Email,
		Password: registerInput.Password,
	}
	newUser, err := r.UserHandler.CreateUser(ctx.Request().Context(), userInput)
	if err != nil {
		return err
	}

	user, err := r.TokenHandler.CreateTokenAndReturnUser(
		ctx.Request().Context(),
		handler.TokenHandlerInput{
			UserId:   newUser.ID,
			Username: newUser.Username,
			Email:    newUser.Email,
		})

	expiryDate := time.Now().Add(utils.OneDay)
	signedAccessToken, err := utils.GenerateNewJwtToken(user.ID, user.Username, user.Email, expiryDate, r.UserHandler.Env.JWT.AccessToken)
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

	refreshClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(utils.OneMonth))

	//	signedRefreshToken, err := utils.GenerateJwtTokenBasedOnExistingClaim(*refreshClaim, r.UserHandler.Env.JWT.RefreshToken)
	//	if err != nil {
	//		return err
	//	}

	accessToken := req.AccessToken

	validAccessToken, err := jwt.ParseWithClaims(accessToken, &utils.CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.UserHandler.Env.JWT.AccessToken), nil
	})

	_, okAcc := validAccessToken.Claims.(*utils.CustomTokenClaim)
	if !okAcc && !validAccessToken.Valid {
		expiryDate := time.Now().Add(utils.OneDay)
		accessToken, err = utils.GenerateNewJwtToken(user.ID, user.Username, user.Email, expiryDate, r.UserHandler.Env.JWT.AccessToken)
		if err != nil {
			return err
		}
	}

	params := db.UpdateUserByIdParams{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
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

func validateUserInputAndGetJwt(input RegisterInput) error {
	if input.Username == "" || input.Email == "" || input.Confirm == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing inputs.")
	}

	if input.Password != input.Confirm {
		return echo.NewHTTPError(http.StatusBadRequest, "Password and confirmation do not match.")
	}

	return nil
}
