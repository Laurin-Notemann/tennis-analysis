package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestRegisterInput struct {
	error TestError
	user  RegisterInput
}

type TestError struct {
	isError       bool
	expectedError error
}

var cfg = config.Config{
	DB: config.DBConfig{
		Url:     "",
		TestUrl: "",
	},
	JWT: config.JwtConfig{
		AccessToken:  "Test",
		RefreshToken: "Test",
	},
}

var tokeGen = utils.MockTokenGenerator{CallOut: 0}
var userHandler = handler.NewUserHandler(utils.DbQueriesTest(), cfg)
var tokenHandler = handler.NewRefreshTokenHandler(utils.DbQueriesTest(), cfg, &tokeGen)

var authRouter = newAuthRouter(*userHandler, *tokenHandler, &tokeGen)

func TestRegisterRoute(t *testing.T) {
	testUserInputData := []TestRegisterInput{
		{
			error: TestError{
				isError:       false,
				expectedError: nil,
			},
			user: RegisterInput{
				Username: "laurin",
				Email:    "laurin@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Password and confirmation do not match.",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "lennart",
				Email:    "lennart@test.de",
				Password: "Test",
				Confirm:  "TestWrong",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusConflict,
					Message:  "pq: duplicate key value violates unique constraint \"users_username_unique\"",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "laurin",
				Email:    "laurin@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusConflict,
					Message:  "pq: duplicate key value violates unique constraint \"users_email_unique\"",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "laulau",
				Email:    "laurin@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "tim",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "tim",
				Email:    "tim@test.de",
				Password: "Test",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "tim",
				Email:    "tim@test.de",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError: true,
				expectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: RegisterInput{
				Username: "",
				Email:    "tim@test.de",
				Confirm:  "Test",
				Password: "Test",
			},
		},
	}

	e := echo.New()

	successAddToDb := 0
	for _, input := range testUserInputData {
		t.Run("register:  "+input.user.Username, func(t *testing.T) {
			encodeUser, err := json.Marshal(input.user)
			if err != nil {
				t.Fatalf("Problem with encoding user %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(encodeUser)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = authRouter.register(c)

			if input.error.isError {
				if assert.Error(t, err) {
					if input.user.Username == "lennart" {
						assert.Equal(
							t,
							input.error.expectedError,
							err,
						)
					} else if input.user.Username == "laurin" {
						assert.Equal(
							t,
							input.error.expectedError,
							err,
						)
					} else if input.user.Username == "laulau" {
						assert.Equal(
							t,
							input.error.expectedError,
							err,
						)
					} else if input.user.Username == "tim" {
						assert.Equal(
							t,
							input.error.expectedError,
							err,
						)
					} else if input.user.Username == "" {
						assert.Equal(
							t,
							input.error.expectedError,
							err,
						)
					}
				}
			} else {
				if assert.NoError(t, err) {
					successAddToDb++
					userRes := new(ResponsePayload)
					err := json.Unmarshal(rec.Body.Bytes(), userRes)
					if err != nil {
						t.Fatalf("Couldn't decode User %v", err)
					}

					allUsers, err := userHandler.GetAllUsers(req.Context())
					if err != nil {
						t.Fatalf("Couldn't get users from db %v", err)
					}

					assert.Equal(t, http.StatusCreated, rec.Code)
					assert.Equal(t, input.user.Username, userRes.User.Username)
					assert.Equal(t, input.user.Email, userRes.User.Email)
					assert.Equal(t, len(allUsers), successAddToDb)
				}
			}
		})
	}
	allUsers, err := userHandler.GetAllUsers(context.Background())
	for _, val := range allUsers {
		_, err = userHandler.DeleteUserById(context.Background(), val.ID)
		assert.NoError(t, err)
	}
}

type RefreshInputTest struct {
	name       string
	validation ValidationType
	error      TestError
	durations  TokenDuration
}

type ValidationType struct {
	validRefresh bool
	vaildAccess  bool
}

type TokenDuration struct {
	access  time.Duration
	refresh time.Duration
}

func TestRefreshRoute(t *testing.T) {
	e := echo.New()

	userInput := RegisterInput{
		Username: "laurin",
		Email:    "laurin@test.de",
		Password: "Test",
		Confirm:  "Test",
	}
	testInputData := []RefreshInputTest{
		{
			name: "valid tokens",
			validation: ValidationType{
				validRefresh: true,
				vaildAccess:  true,
			},
			error: TestError{
				isError:       false,
				expectedError: nil,
			},
			durations: TokenDuration{
				access:  5 * time.Minute,
				refresh: 5 * time.Minute,
			},
		},
		{
			name: "valid refresh token",
			validation: ValidationType{
				validRefresh: true,
				vaildAccess:  false,
			},
			error: TestError{
				isError:       false,
				expectedError: nil,
			},
			durations: TokenDuration{
				access:  time.Duration(0),
				refresh: 5 * time.Minute,
			},
		},
		{
			name: "invalid tokens",
			validation: ValidationType{
				validRefresh: false,
				vaildAccess:  false,
			},
			error: TestError{
				isError:       true,
				expectedError: echo.NewHTTPError(http.StatusUnauthorized, jwt.ErrTokenInvalidClaims.Error()+": "+jwt.ErrTokenExpired.Error()),
			},
			durations: TokenDuration{
				access:  time.Duration(0),
				refresh: time.Duration(0),
			},
		},
	}

	for _, data := range testInputData {
		t.Run("/refresh "+data.name, func(t *testing.T) {

			err, rec, registeredUser, req := refreshUser(t, e, userInput, data.durations.access, data.durations.refresh)

			if !data.validation.vaildAccess && !data.validation.validRefresh {
				if assert.Error(t, err) {
					assert.Equal(t, data.error.expectedError, err)
				}
			} else if !data.validation.vaildAccess && data.validation.validRefresh {
				if assert.NoError(t, err) {
					refreshedUser := new(ResponsePayload)

					err := json.Unmarshal(rec.Body.Bytes(), refreshedUser)
					if err != nil {
						t.Fatalf("Couldn't decode User %v", err)
					}

					assert.Equal(t, registeredUser.User, refreshedUser.User)
					assert.NotEqual(t, registeredUser.AccessToken, refreshedUser.AccessToken)
				}
			} else if data.validation.vaildAccess && data.validation.validRefresh {
				if assert.NoError(t, err) {
					refreshedUser := new(ResponsePayload)

					err := json.Unmarshal(rec.Body.Bytes(), refreshedUser)
					if err != nil {
						t.Fatalf("Couldn't decode User %v", err)
					}
					assert.Equal(t, registeredUser.User, refreshedUser.User)
					assert.Equal(t, registeredUser.AccessToken, refreshedUser.AccessToken)
				}
			}
			_, err = userHandler.DeleteUserById(req.Context(), registeredUser.User.ID)
			assert.NoError(t, err)
		})
	}

}

func registerNewUser(t *testing.T, e *echo.Echo, userData RegisterInput, durAcc time.Duration, durRef time.Duration) *ResponsePayload {
	encodeUser, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Problem with encoding user %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(encodeUser)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	tokeGen.ExpiryDateAccess = durAcc
	tokeGen.ExpiryDateRefresh = durRef

	err = authRouter.register(c)
	assert.NoError(t, err)

	user := new(ResponsePayload)
	err = json.Unmarshal(rec.Body.Bytes(), user)
	if err != nil {
		t.Fatalf("Couldn't decode User %v", err)
	}

	return user
}

func refreshUser(
	t *testing.T,
	e *echo.Echo,
	userInput RegisterInput,
	durAcc time.Duration,
	durRef time.Duration,
) (error, *httptest.ResponseRecorder, *ResponsePayload, *http.Request) {
	user := registerNewUser(t, e, userInput, durAcc, durRef)

	encodeRefreshReq, err := json.Marshal(RefreshReq{AccessToken: user.AccessToken})
	assert.NoError(t, err)

	if durAcc == time.Duration(0) {
		tokeGen.ExpiryDateAccess = 5 * time.Minute
	}

	req := httptest.NewRequest(http.MethodPost, "/api/refresh", strings.NewReader(string(encodeRefreshReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = authRouter.refresh(c)

	return err, rec, user, req
}