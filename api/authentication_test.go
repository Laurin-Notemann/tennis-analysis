package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type (
	TestRegisterInput struct {
		error TestError
		user  handler.RegisterInput
	}
	RefreshInputTest struct {
		name       string
		validation ValidationType
		error      TestError
		durations  TokenDuration
	}

	ValidationType struct {
		validRefresh bool
		vaildAccess  bool
	}

	TokenDuration struct {
		access  time.Duration
		refresh time.Duration
	}

	LoginInputTest struct {
		name      string
		error     TestError
		userInput handler.LoginInput
		durations TokenDuration
	}
)

var tokeGen = utils.MockTokenGenerator{CallOut: 0}
var userHandler = handler.NewUserHandler(utils.DbQueriesTest(), Cfg)
var tokenHandler = handler.NewRefreshTokenHandler(utils.DbQueriesTest(), Cfg, &tokeGen)
var authHandler = handler.NewAuthenticationHandler(utils.DbQueriesTest(), *userHandler, *tokenHandler, &tokeGen)

var authRouter = NewAuthRouter(*userHandler, *tokenHandler, &tokeGen, *authHandler)

func TestRegisterRoute(t *testing.T) {
	e := echo.New()

	testUserInputData := []TestRegisterInput{
		{
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			user: handler.RegisterInput{
				Username: "laurin",
				Email:    "laurin@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Password and confirmation do not match.",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "lennart",
				Email:    "lennart@test.de",
				Password: "Test",
				Confirm:  "TestWrong",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusConflict,
					Message:  "pq: duplicate key value violates unique constraint \"users_username_unique\"",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "laurin",
				Email:    "laurin@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusConflict,
					Message:  "pq: duplicate key value violates unique constraint \"users_email_unique\"",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "laulau",
				Email:    "laurin@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "tim",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "tim",
				Email:    "tim@test.de",
				Password: "Test",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "tim",
				Email:    "tim@test.de",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "Missing inputs.",
					Internal: error(nil),
				},
			},
			user: handler.RegisterInput{
				Username: "",
				Email:    "tim@test.de",
				Confirm:  "Test",
				Password: "Test",
			},
		},
	}

	successAddToDb := 0
	for _, input := range testUserInputData {
		t.Run("/api/register:  "+input.user.Username, func(t *testing.T) {
			encodeUser, err := json.Marshal(input.user)
			if err != nil {
				t.Fatalf("Problem with encoding user %v", err)
			}

			err, rec, req := DummyRequest(t, e, http.MethodPost, "/api/register", string(encodeUser), authRouter.Register, "")

			if input.error.IsError {
				if assert.Error(t, err) {
					if input.user.Username == "lennart" {
						assert.Equal(
							t,
							input.error.ExpectedError,
							err,
						)
					} else if input.user.Username == "laurin" {
						assert.Equal(
							t,
							input.error.ExpectedError,
							err,
						)
					} else if input.user.Username == "laulau" {
						assert.Equal(
							t,
							input.error.ExpectedError,
							err,
						)
					} else if input.user.Username == "tim" {
						assert.Equal(
							t,
							input.error.ExpectedError,
							err,
						)
					} else if input.user.Username == "" {
						assert.Equal(
							t,
							input.error.ExpectedError,
							err,
						)
					}
				}
			} else {
				if assert.NoError(t, err) {
					successAddToDb++
					userRes := new(handler.ResponsePayload)
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

func TestRefreshRoute(t *testing.T) {
	e := echo.New()

	userInput := handler.RegisterInput{
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
				IsError:       false,
				ExpectedError: nil,
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
				IsError:       false,
				ExpectedError: nil,
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
				IsError:       true,
				ExpectedError: echo.NewHTTPError(http.StatusUnauthorized, jwt.ErrTokenInvalidClaims.Error()+": "+jwt.ErrTokenExpired.Error()),
			},
			durations: TokenDuration{
				access:  time.Duration(0),
				refresh: time.Duration(0),
			},
		},
	}

	for _, data := range testInputData {
		t.Run("/api/refresh "+data.name, func(t *testing.T) {

			err, rec, registeredUser, req := refreshUser(t, e, userInput, data.durations.access, data.durations.refresh)

			if !data.validation.vaildAccess && !data.validation.validRefresh {
				if assert.Error(t, err) {
					assert.Equal(t, data.error.ExpectedError, err)
				}
			} else if !data.validation.vaildAccess && data.validation.validRefresh {
				if assert.NoError(t, err) {
					refreshedUser := new(handler.ResponsePayload)

					err := json.Unmarshal(rec.Body.Bytes(), refreshedUser)
					if err != nil {
						t.Fatalf("Couldn't decode User %v", err)
					}

					assert.Equal(t, registeredUser.User, refreshedUser.User)
					assert.NotEqual(t, registeredUser.AccessToken, refreshedUser.AccessToken)
				}
			} else if data.validation.vaildAccess && data.validation.validRefresh {
				if assert.NoError(t, err) {
					refreshedUser := new(handler.ResponsePayload)

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

func TestLoginRoute(t *testing.T) {
	e := echo.New()

	testDataLogin := []LoginInputTest{
		{
			name: "correct login with username",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			userInput: handler.LoginInput{
				UsernameOrEmail: "laurin",
				Password:        "Test",
			},
			durations: TokenDuration{
				access:  5 * time.Minute,
				refresh: 5 * time.Minute,
			},
		},
		{
			name: "correct login with email",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			userInput: handler.LoginInput{
				UsernameOrEmail: "laurin@test.de",
				Password:        "Test",
			},
			durations: TokenDuration{
				access:  5 * time.Minute,
				refresh: 5 * time.Minute,
			},
		},
		{
			name: "wrong login with missmatched pw",
			error: TestError{
				IsError: true,
				ExpectedError: &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  bcrypt.ErrMismatchedHashAndPassword.Error(),
					Internal: error(nil),
				},
			},
			userInput: handler.LoginInput{
				UsernameOrEmail: "laurin",
				Password:        "TestWrong",
			},
			durations: TokenDuration{
				access:  5 * time.Minute,
				refresh: 5 * time.Minute,
			},
		},
	}

	for _, data := range testDataLogin {
		t.Run("/api/login", func(t *testing.T) {
			user := DummyUser(t, e)

			encodeLoginReq, err := json.Marshal(handler.LoginInput{UsernameOrEmail: data.userInput.UsernameOrEmail, Password: data.userInput.Password})
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(string(encodeLoginReq)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err = authRouter.login(c)

			if data.error.IsError {
				if assert.Error(t, err) {
					assert.Equal(t, data.error.ExpectedError.Error(), err.Error())
				}
			} else {
				if assert.NoError(t, err) {
					loggedInUser := new(handler.ResponsePayload)
					err := json.Unmarshal(rec.Body.Bytes(), loggedInUser)
					assert.NoError(t, err, "Couldn't decode User")

					assert.NotEqual(t, user, loggedInUser)
					assert.Equal(t, user.ID, loggedInUser.User.ID)
					assert.Equal(t, user.Username, loggedInUser.User.Username)
					assert.Equal(t, user.Email, loggedInUser.User.Email)
					assert.NotEqual(t, user.RefreshTokenID, loggedInUser.User.RefreshTokenID)
				}
			}
			_, err = userHandler.DeleteUserById(req.Context(), user.ID)
			assert.NoError(t, err)
		})
	}

}

func refreshUser(
	t *testing.T,
	e *echo.Echo,
	userInput handler.RegisterInput,
	durAcc time.Duration,
	durRef time.Duration,
) (error, *httptest.ResponseRecorder, *handler.ResponsePayload, *http.Request) {
	user := RegisterDummyUser(t, e, userInput, &tokeGen, durAcc, durRef)

	encodeRefreshReq, err := json.Marshal(handler.RefreshReq{AccessToken: user.AccessToken})
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
