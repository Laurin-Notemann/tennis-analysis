package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/tests"
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

func TestAuthRoute(t *testing.T) {
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
				isError:       false,
				expectedError: nil,
			},
			user: RegisterInput{
				Username: "paulo",
				Email:    "paulo@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError:       false,
				expectedError: nil,
			},
			user: RegisterInput{
				Username: "florian",
				Email:    "florian@test.de",
				Password: "Test",
				Confirm:  "Test",
			},
		},
		{
			error: TestError{
				isError:       false,
				expectedError: nil,
			},
			user: RegisterInput{
				Username: "max",
				Email:    "max@test.de",
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
				Username: "max",
				Email:    "max@test.de",
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
				Email: "tim@test.de",
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
				Email: "tim@test.de",
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
				Email: "tim@test.de",
				Confirm:  "Test",
				Password: "Test",
			},
		},
	}

	e := echo.New()

	cfg := config.Config{
		DB: config.DBConfig{
			Url:     "",
			TestUrl: "",
		},
		JWT: config.JwtConfig{
			AccessToken:  "Test",
			RefreshToken: "Test",
		},
	}

	userHandler := handler.NewUserHandler(tests.DbQueriesTest(), cfg)

	authRouter := newAuthRouter(*userHandler)

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
					} else if input.user.Username == "max" {
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
}
