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
	"github.com/google/uuid"
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

var playerHandler = handler.NewPlayerHandler(utils.DbQueriesTest())
var teamHandler = handler.NewTeamHandler(utils.DbQueriesTest())
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
				"",
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
	userIdString := userId.String()

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
		t.Run("get players "+data.name, func(t *testing.T) {
			addMultiplePlayers(t, e, data.input)
			url := "/api/players/:id"
			t.Log(url)
			err, rec, _ := DummyRequest(
				t,
				e,
				http.MethodGet,
				url,
				string(""),
				playRouter.GetAllPlayersByUserId,
				userIdString,
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
	_, err := userHandler.DeleteUserById(context.Background(), userId)
	assert.NoError(t, err)
}

type TestDeletePlayer struct {
	name  string
	error TestError
	input uuid.UUID
}

func TestDeletePlayerById(t *testing.T) {
	e := echo.New()

	user := DummyUser(t, e)
	userId := user.ID

	player := DummyPlayer(t, e, userId)

	testDataInput := []TestDeletePlayer{
		{
			name: "success",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: player.ID,
		},
	}

	for _, data := range testDataInput {
		t.Run("delete player "+data.name, func(t *testing.T) {
			encodedData, err := json.Marshal(data.input)
			assert.NoError(t, err, "Problem with encoding the id")

			url := "/api/players/:id"
			err, rec, _ := DummyRequest(
				t,
				e,
				http.MethodDelete,
				url,
				string(encodedData),
				playRouter.DeletePlayerById,
				player.ID.String(),
			)
			if data.error.IsError {
				if assert.Error(t, err) {
					assert.Equal(t, data.error.ExpectedError, err)
				}
			} else {
				if assert.NoError(t, err, "Error with DeletePlayer route") {
					deletedPlayer := new(db.Player)
					err := json.Unmarshal(rec.Body.Bytes(), deletedPlayer)
					assert.NoError(t, err, "Couldn't decode deleted Player")

					assert.Equal(t, player, *deletedPlayer)
				}
			}
		})
	}
	_, err := userHandler.DeleteUserById(context.Background(), userId)
	assert.NoError(t, err)
}

type TestUpdatePlayer struct {
	name  string
	error TestError
	input db.UpdatePlayerByIdParams
}

func TestUpdatePlayerById(t *testing.T) {
	e := echo.New()

	user := DummyUser(t, e)
	userId := user.ID

	player := DummyPlayer(t, e, userId)

	testDataInput := []TestUpdatePlayer{
		{
			name: "success",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: db.UpdatePlayerByIdParams{
				ID:        player.ID,
				FirstName: "Oskar",
				LastName:  "Kuech",
			},
		},
	}

	for _, data := range testDataInput {
		t.Run("delete player "+data.name, func(t *testing.T) {
			encodedData, err := json.Marshal(data.input)
			assert.NoError(t, err, "Problem with encoding the update params")

			url := "/api/players"
			err, rec, _ := DummyRequest(
				t,
				e,
				http.MethodPut,
				url,
				string(encodedData),
				playRouter.UpdatePlayerById,
				"",
			)
			if data.error.IsError {
				if assert.Error(t, err) {
					assert.Equal(t, data.error.ExpectedError, err)
				}
			} else {
				if assert.NoError(t, err, "Error with updatePlayer route") {
					updatedPlayer := new(db.Player)
					err := json.Unmarshal(rec.Body.Bytes(), updatedPlayer)
					assert.NoError(t, err, "Couldn't decode updated Player")

					assert.NotEqual(t, player, *updatedPlayer)
					assert.Equal(t, data.input.FirstName, updatedPlayer.FirstName)
					assert.Equal(t, data.input.LastName, updatedPlayer.LastName)
					assert.Equal(t, data.input.ID, updatedPlayer.ID)
				}
			}
		})
	}
	_, err := userHandler.DeleteUserById(context.Background(), userId)
	assert.NoError(t, err)
}

func addMultiplePlayers(t *testing.T, e *echo.Echo, input []db.CreateNewTeamWithOnePlayerParams) {
	for _, data := range input {
		addNewPlayer(t, e, data)
	}
}

func DummyPlayer(t *testing.T, e *echo.Echo, userId uuid.UUID) db.Player {
	seed := db.CreateNewTeamWithOnePlayerParams{
		FirstName: "Laurin",
		LastName:  "Notemann",
		Name: sql.NullString{
			String: "",
			Valid:  false,
		},
		UserID: userId,
	}
	return addNewPlayer(t, e, seed)
}

func addNewPlayer(t *testing.T, e *echo.Echo, data db.CreateNewTeamWithOnePlayerParams) db.Player {
	encodedData, err := json.Marshal(data)
	assert.NoError(t, err, "Problem with encoding the player")

	err, rec, _ := DummyRequest(
		t,
		e,
		http.MethodPost,
		"/api/players",
		string(encodedData),
		playRouter.CreatePlayer,
		"",
	)
	assert.NoError(t, err, "Problem with adding new player")

	player := new(db.Player)
	err = json.Unmarshal(rec.Body.Bytes(), player)
	assert.NoError(t, err, "Couldn't decode list of Players")

	return *player
}
