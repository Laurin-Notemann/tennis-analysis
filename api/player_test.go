package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestPlayerInput struct {
	name  string
	error TestError
	input db.CreateNewTeamWithOnePlayerParams
}

var testDb = utils.DbQueriesTest()
var playerHandler = handler.NewPlayerHandler(testDb)
var teamHandler = handler.NewTeamHandler(testDb)
var playRouter = newPlayerRouter(*playerHandler, *teamHandler, *userHandler)

func TestCreatePlayer(t *testing.T) {
	e := echo.New()
	user := DummyUser(t, e)
	userId := user.ID
	testCreatePlayerInput := []TestPlayerInput{
		{
			name: "success new player",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
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
		{
			name: "error new player",
			error: TestError{
				IsError:       true,
				ExpectedError: &echo.HTTPError{Code: 400, Message: "missing first or last name", Internal: error(nil)},
			},
			input: db.CreateNewTeamWithOnePlayerParams{
				FirstName: "",
				LastName:  "Kuech",
				Name: sql.NullString{
					String: "",
					Valid:  false,
				},
				UserID: userId,
			},
		},
		{
			name: "error new player",
			error: TestError{
				IsError:       true,
				ExpectedError: &echo.HTTPError{Code: 400, Message: "missing first or last name", Internal: error(nil)},
			},
			input: db.CreateNewTeamWithOnePlayerParams{
				FirstName: "Oskar",
				LastName:  "",
				Name: sql.NullString{
					String: "",
					Valid:  false,
				},
				UserID: userId,
			},
		},
		{
			name: "error new player",
			error: TestError{
				IsError:       true,
				ExpectedError: &echo.HTTPError{Code: 409, Message: "player already exists", Internal: error(nil)},
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
		{
			name: "success new player",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: db.CreateNewTeamWithOnePlayerParams{
				FirstName: "Oskar",
				LastName:  "Test",
				Name: sql.NullString{
					String: "",
					Valid:  false,
				},
				UserID: userId,
			},
		},
		{
			name: "success new player",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: db.CreateNewTeamWithOnePlayerParams{
				FirstName: "Laurin",
				LastName:  "Test",
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
			"/api/players",
			string(input),
			playRouter.CreatePlayer,
		)
		if data.error.IsError {
      if assert.Error(t, err) {
        assert.Equal(t, data.error.ExpectedError, err)
      }
		} else {
			if assert.NoError(t, err, "Error with CreatePlayer route") {
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

