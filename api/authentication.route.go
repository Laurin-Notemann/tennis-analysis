package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/lib/pq"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
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

	CustomTokenClaim struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		jwt.RegisteredClaims
	}

	AuthenticationRouter struct {
		UserHandler handler.UserHandler
		env         config.Config
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

func newAuthRouter(h handler.UserHandler, env config.Config) *AuthenticationRouter {
	return &AuthenticationRouter{UserHandler: h, env: env}
}

func (r AuthenticationRouter) register(ctx echo.Context) (err error) {
	newUser := new(RegisterInput)

	if err = ctx.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if newUser.Username == "" || newUser.Email == "" || newUser.Confirm == "" || newUser.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing inputs.")
	}

	if newUser.Password != newUser.Confirm {
		return echo.NewHTTPError(http.StatusBadRequest, "Password and confirmation do not match.")
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	signedAccessToken, err := generateNewJwtToken(newUser.Username, newUser.Email, oneDay, r.env.JWT.AccessToken)
	if err != nil {
		return err
	}

	signedRefreshToken, err := generateNewJwtToken(newUser.Username, newUser.Email, oneMonth, r.env.JWT.RefreshToken)
	if err != nil {
		return err
	}

	user, err := r.UserHandler.CreateUser(ctx.Request().Context(), db.CreateUserParams{
		Username:     newUser.Username,
		Email:        newUser.Email,
		PasswordHash: string(hashedPw),
		RefreshToken: sql.NullString{
			String: signedRefreshToken,
			Valid:  true,
		},
	})

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {

		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	validRefreshToken, err := jwt.ParseWithClaims(refreshToken, &CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.env.JWT.RefreshToken), nil
	})

	refreshClaim, ok := validRefreshToken.Claims.(*CustomTokenClaim)
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

	signedRefreshToken, err := generateJwtTokenBasedOnExistingClaim(*refreshClaim, r.env.JWT.RefreshToken)
	if err != nil {
		return err
	}

	accessToken := req.AccessToken

	validAccessToken, err := jwt.ParseWithClaims(accessToken, &CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.env.JWT.AccessToken), nil
	})

	_, okAcc := validAccessToken.Claims.(*CustomTokenClaim)
	if !okAcc && !validAccessToken.Valid {
		accessToken, err = generateNewJwtToken(user.Username, user.Email, oneDay, r.env.JWT.AccessToken)
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

func generateNewJwtToken(username string, email string, expiryDate time.Duration, signingKey string) (string, error) {
	tokenClaim := CustomTokenClaim{
		username,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDate)),
		},
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)

	token, err := tokenObj.SignedString([]byte(signingKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}

func generateJwtTokenBasedOnExistingClaim(claim CustomTokenClaim, signingKey string) (string, error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := tokenObj.SignedString([]byte(signingKey))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}
