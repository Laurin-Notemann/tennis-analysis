package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
)

type AuthenticationRouter struct {
	UserHandler  handler.UserHandler
	TokenHandler handler.RefreshTokenHandler
	TokenGen     utils.TokenGenerator
	AuthHandler  handler.AuthenticationHandler
}

func newAuthRouter(
	h handler.UserHandler,
	t handler.RefreshTokenHandler,
	tg utils.TokenGenerator,
	a handler.AuthenticationHandler,
) *AuthenticationRouter {
	return &AuthenticationRouter{
		UserHandler:  h,
		TokenHandler: t,
		TokenGen:     tg,
		AuthHandler:  a,
	}
}

func (r AuthenticationRouter) register(ctx echo.Context) (err error) {
	registerInput := new(handler.RegisterInput)
	if err = ctx.Bind(registerInput); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = r.AuthHandler.ValidateUserInputAndGetJwt(*registerInput)
	if err != nil {
		return err
	}

	user, err := r.AuthHandler.CreateUserAndToken(ctx, *registerInput)
	if err != nil {
		return err
	}

	accessToken, err := r.AuthHandler.GenAccessToken(&user)
	if err != nil {
		return err
	}

	registerPayload := handler.ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}
	return ctx.JSON(http.StatusCreated, registerPayload)
}

func (r AuthenticationRouter) refresh(ctx echo.Context) (err error) {
	req := new(handler.RefreshReq)
	if err = ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	accessToken := req.AccessToken

	user, validToken, err := r.AuthHandler.ParseTokenGetUser(accessToken, ctx)
	if err != nil {
		return err
	}

	payload, err := r.AuthHandler.ValidateAccessToken(accessToken, validToken, user)
	if !errors.Is(err, handler.AccessTokenInvalid) {
		return ctx.JSON(http.StatusOK, payload)
	}

	err = r.AuthHandler.ValidateRefreshToken(ctx, user)
	if err != nil {
		return err
	}

	accessToken, err = r.AuthHandler.GenAccessToken(&user)
	if err != nil {
		return err
	}

	args := handler.TokenHandlerInput{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	_, err = r.TokenHandler.UpdateTokenByUserId(ctx.Request().Context(), args)

	payload = handler.ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}
	return ctx.JSON(http.StatusOK, payload)
}

func (r AuthenticationRouter) login(ctx echo.Context) (err error) {
	loginReq := new(handler.LoginInput)
	if err = ctx.Bind(loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user db.User
	if strings.Contains(loginReq.UsernameOrEmail, "@") {
		user, err = r.UserHandler.GetUserByEmail(ctx.Request().Context(), loginReq.UsernameOrEmail)
	} else {
		user, err = r.UserHandler.GetUserByUsername(ctx.Request().Context(), loginReq.UsernameOrEmail)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, err = r.AuthHandler.LoginRefreshToken(ctx, &user)

	accessToken, err := r.AuthHandler.GenAccessToken(&user)

	payload := handler.ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}

	return ctx.JSON(http.StatusOK, payload)
}

func RegisterAuthRoute(baseUrl string, e *echo.Echo, r AuthenticationRouter) {

	e.POST(baseUrl+"/register", r.register)
	e.POST(baseUrl+"/refresh", r.refresh)
	e.POST(baseUrl+"/login", r.login)
}
