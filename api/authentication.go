package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
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

	AuthenticationRouter struct {
		UserHandler  handler.UserHandler
		TokenHandler handler.RefreshTokenHandler
		TokenGen     utils.TokenGenerator
	}

	ResponsePayload struct {
		AccessToken string  `json:"accessToken"`
		User        db.User `json:"user"`
	}
)

func newAuthRouter(
	h handler.UserHandler,
	t handler.RefreshTokenHandler,
	tg utils.TokenGenerator,
) *AuthenticationRouter {
	return &AuthenticationRouter{UserHandler: h, TokenHandler: t, TokenGen: tg}
}

func (r AuthenticationRouter) register(ctx echo.Context) (err error) {
	registerInput := new(RegisterInput)
	if err = ctx.Bind(registerInput); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = r.validateUserInputAndGetJwt(*registerInput)
	if err != nil {
		return err
	}
	user, err := r.createUserAndToken(ctx, *registerInput)
	if err != nil {
		return err
	}

	accessToken, err := r.genAccessToken(&user)
	if err != nil {
		return err
	}

	registerPayload := ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}
	return ctx.JSON(http.StatusCreated, registerPayload)
}

func (r AuthenticationRouter) refresh(ctx echo.Context) (err error) {
	req := new(RefreshReq)
	if err = ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	accessToken := req.AccessToken

	user, validToken, err := r.parseTokenGetUser(accessToken, ctx)
	if err != nil {
		return err
	}

	payload, err := r.validateAccessToken(accessToken, validToken, user)
	if !errors.Is(err, AccessTokenInvalid) {
		return ctx.JSON(http.StatusOK, payload)
	}

	err = r.validateRefreshToken(ctx, user)
	if err != nil {
		return err
	}

	accessToken, err = r.genAccessToken(&user)
	if err != nil {
		return err
	}

	args := handler.TokenHandlerInput{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	_, err = r.TokenHandler.UpdateTokenByUserId(ctx.Request().Context(), args)

	payload = ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}
	return ctx.JSON(http.StatusOK, payload)
}

func (r AuthenticationRouter) login(ctx echo.Context) (err error) {
	loginReq := new(LoginInput)
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, err = r.loginRefreshToken(ctx, &user)

	accessToken, err := r.genAccessToken(&user)

	payload := ResponsePayload{
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

func (r *AuthenticationRouter) validateUserInputAndGetJwt(input RegisterInput) error {
	if input.Username == "" || input.Email == "" || input.Confirm == "" || input.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing inputs.")
	}
	if input.Password != input.Confirm {
		return echo.NewHTTPError(http.StatusBadRequest, "Password and confirmation do not match.")
	}
	return nil
}

func (r *AuthenticationRouter) validateAccessToken(accessToken string, valid *jwt.Token, user db.User) (ResponsePayload, error) {
	if !valid.Valid {
		return ResponsePayload{}, AccessTokenInvalid
	}
	payload := ResponsePayload{
		AccessToken: accessToken,
		User:        user,
	}
	return payload, nil
}

func (r *AuthenticationRouter) parseTokenGetUser(accessToken string, ctx echo.Context) (db.User, *jwt.Token, error) {
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

func (r *AuthenticationRouter) validateRefreshToken(ctx echo.Context, user db.User) error {
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

func (r *AuthenticationRouter) createUserAndToken(ctx echo.Context, registerInput RegisterInput) (db.User, error) {
	userInput := handler.CreateUserInput{
		Username: registerInput.Username,
		Email:    registerInput.Email,
		Password: registerInput.Password,
	}
	newUser, err := r.UserHandler.CreateUser(ctx.Request().Context(), userInput)
	if err != nil {
		return db.User{}, err
	}

	user, err := r.genRefreshToken(ctx, &newUser)
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (r *AuthenticationRouter) genAccessToken(user *db.User) (string, error) {
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

func (r *AuthenticationRouter) genRefreshToken(ctx echo.Context, user *db.User) (db.User, error) {
	updatedUser, err := r.TokenHandler.CreateTokenAndReturnUser(
		ctx.Request().Context(),
		handler.TokenHandlerInput{
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

func (r *AuthenticationRouter) loginRefreshToken(ctx echo.Context, user *db.User) (db.User, error) {
	err := r.TokenHandler.DeleteTokenByUserId(ctx.Request().Context(), user.ID)
	if err != nil {
		return db.User{}, err
	}

	refUser, err := r.genRefreshToken(ctx, user)
	return refUser, err
}
