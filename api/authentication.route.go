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

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

type CustomTokenClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type AuthenticationRouter struct {
	UserHandler handler.UserHandler
	env         config.Config
}

type OK = string
type ERR = string

type Response struct {
	Status string
	Res    interface{}
}

type ResponsePayload struct {
	AccessToken string  `json:"accessToken"`
	User        db.User `json:"user"`
}

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

	accessTokenClaim := CustomTokenClaim{
		newUser.Username,
		newUser.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)

	signedAccessToken, err := accessTokenObj.SignedString([]byte(r.env.JWT.AccessToken))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshTokenClaim := CustomTokenClaim{
		newUser.Username,
		newUser.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour)),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)

	signedRefreshToken, err := refreshTokenObj.SignedString([]byte(r.env.JWT.AccessToken))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

type RefreshReq struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	refreshClaim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour))

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	// create new refreshToken (so only until the user hasn't been logged in for 30 days will he have to log in again)
	signedRefreshToken, err := refreshTokenObj.SignedString([]byte(r.env.JWT.AccessToken))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	accessToken := req.AccessToken

	validAccessToken, err := jwt.ParseWithClaims(accessToken, &CustomTokenClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(r.env.JWT.AccessToken), nil
	})

	_, okAcc := validAccessToken.Claims.(*CustomTokenClaim)
	if !okAcc && !validAccessToken.Valid {
		accessTokenClaim := CustomTokenClaim{
			user.Username,
			user.Email,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}
		accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaim)

		signedAccessToken, err := accessTokenObj.SignedString([]byte(r.env.JWT.AccessToken))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		accessToken = signedAccessToken
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
