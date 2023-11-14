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

type TestGetAllPlayer struct {
	name           string
	error          TestError
	input          []db.CreateNewTeamWithOnePlayerParams
	expectedLength int
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
		t.Run("create player "+data.name, func(t *testing.T) {
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
		})
	}
	_, err := userHandler.DeleteUserById(context.Background(), userId)
	assert.NoError(t, err)
}

func TestGetAllPlayersByUserId(t *testing.T) {
	e := echo.New()

	user := DummyUser(t, e)
	userId := user.ID

	testDataInput := []TestGetAllPlayer{
		{
			name: "success",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: []db.CreateNewTeamWithOnePlayerParams{
				{
					FirstName: "Laurin",
					LastName:  "Notemann",
					Name: sql.NullString{
						String: "",
						Valid:  false,
					},
					UserID: userId,
				},
				{
					FirstName: "Max",
					LastName:  "Mustermann",
					Name: sql.NullString{
						String: "",
						Valid:  false,
					},
					UserID: userId,
				},
			},
			expectedLength: 2,
		},
	}

	for _, data := range testDataInput {
		t.Run("create player "+data.name, func(t *testing.T) {
			err, rec, _ := DummyRequest(
				t,
				e,
				http.MethodGet,
				"/api/players",
				string(""),
				playRouter.GetAllPlayersByUserId,
			)
			if data.error.IsError {
				if assert.Error(t, err) {
					assert.Equal(t, data.error.ExpectedError, err)
				}
			} else {
				if assert.NoError(t, err, "Error with CreatePlayer route") {
					allPlayers := new([]db.Player)
					err := json.Unmarshal(rec.Body.Bytes(), allPlayers)
					assert.NoError(t, err, "Couldn't decode list of Players")

					assert.Equal(t, data.expectedLength, len(*allPlayers))
				}
			}
		})
	}
}

func addMultiplePlayers(t *testing.T, e *echo.Echo, input []db.CreateNewTeamWithOnePlayerParams) {
	for _, data := range input {
		encodedData, err := json.Marshal(data)
		assert.NoError(t, err, "Problem with encoding the user")

		err, _, _ = DummyRequest(
			t,
			e,
			http.MethodPost,
			"/api/players",
			string(encodedData),
			playRouter.CreatePlayer,
		)
		assert.NoError(t, err, "Problem with adding new player")
	}

}
