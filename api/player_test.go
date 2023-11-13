package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestPlayerInput struct {
	name  string
	error TestError
	input db.CreateNewTeamWithOnePlayerParams
}

var playRouter = newPlayerRouter(userHandler)

func TestCreatePlayer(t *testing.T) {
	e := echo.New()
	user := testUser(t, e)
	userId := user.ID
	testCreatePlayerInput := []TestPlayerInput{
		{
			name: "success new player",
			error: TestError{
				isError:       false,
				expectedError: nil,
			},
			input: db.CreateNewTeamWithOnePlayerParams{
				FirstName: "Oskar",
				LastName:  "Kuech",
				Name: sql.NullString{
					String: "",
					Valid:  false,
				},
				UserID: userId,
			},
		},
	}
	for _, data := range testCreatePlayerInput {
		input, err := json.Marshal(data.input)
		assert.NoError(t, err, "Problem with encoding the user")

		err, rec, _ := DummyRequest(
			t,
			e,
			http.MethodPost,
			"api/players",
			string(input),
			playRouter.CreatePlayer,
		)
		if data.error.isError {
		} else {
			if assert.NoError(t, err) {
				newPlayer := new(db.Player)
				err := json.Unmarshal(rec.Body.Bytes(), newPlayer)
				assert.NoError(t, err, "Couldn't decode Player")

				assert.Equal(t, data.input.FirstName, newPlayer.FirstName)
				assert.Equal(t, data.input.LastName, newPlayer.LastName)
			}
		}
	}
  _, err := userHandler.DeleteUserById(context.Background(), userId)
  assert.NoError(t, err)
}

func testUser(t *testing.T, e *echo.Echo) db.User {
	testUserInput := handler.RegisterInput{
		Username: "laurin",
		Email:    "laurin@test",
		Password: "Test",
		Confirm:  "Test",
	}
	user := registerNewUser(t, e, testUserInput, 5*time.Minute, 5*time.Minute)
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
