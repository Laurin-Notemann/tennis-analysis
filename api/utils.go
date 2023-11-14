package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestError struct {
	IsError       bool
	ExpectedError error
}

var Cfg = config.Config{
	DB: config.DBConfig{
		Url:     "",
		TestUrl: "",
	},
	JWT: config.JwtConfig{
		AccessToken:  "Test",
		RefreshToken: "Test",
	},
}

var TestDb = utils.DbQueriesTest()

func RegisterDummyUser(t *testing.T, e *echo.Echo, userData handler.RegisterInput, tokenGen *utils.MockTokenGenerator, durAcc time.Duration, durRef time.Duration) *handler.ResponsePayload {
	var userHandler = handler.NewUserHandler(TestDb, Cfg)
	var tokenHandler = handler.NewRefreshTokenHandler(TestDb, Cfg, tokenGen)
	var authHandler = handler.NewAuthenticationHandler(TestDb, *userHandler, *tokenHandler, tokenGen)
	var authRouter = NewAuthRouter(*userHandler, *tokenHandler, tokenGen, *authHandler)

	tokenGen.ExpiryDateAccess = durAcc
	tokenGen.ExpiryDateRefresh = durRef
	encodeUser, err := json.Marshal(userData)
	assert.NoError(t, err, "Problem with encoding the user")

	err, rec, _ := DummyRequest(t, e, http.MethodPost, "/api/register", string(encodeUser), authRouter.Register)
	assert.NoError(t, err, "Problem with registering test user")

	user := new(handler.ResponsePayload)
	err = json.Unmarshal(rec.Body.Bytes(), user)
	if err != nil {
		t.Fatalf("Couldn't decode User %v", err)
	}

	return user
}

func DummyUser(t *testing.T, e *echo.Echo) db.User {
	testUserInput := handler.RegisterInput{
		Username: "laurin",
		Email:    "laurin@test.de",
		Password: "Test",
		Confirm:  "Test",
	}
	user := RegisterDummyUser(t, e, testUserInput, &utils.MockTokenGenerator{}, 5*time.Minute, 5*time.Minute)
	return user.User
}

func DummyRequest(
	t *testing.T,
	e *echo.Echo,
	method string,
	url string,
	encodedInput string,
	routerFunc echo.HandlerFunc,
) (err error, rec *httptest.ResponseRecorder, req *http.Request) {
	req = httptest.NewRequest(method, url, strings.NewReader(string(encodedInput)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = routerFunc(c)
	return err, rec, req
}
