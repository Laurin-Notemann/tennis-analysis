package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestCreateTeamInput struct {
	name  string
	error TestError
	input db.CreateTeamWithTwoPlayersParams
}

var teamRouter = newTeamRouter(*playerHandler, *teamHandler, *userHandler)

func TestCreateTeam(t *testing.T) {
	e := echo.New()
	userId, playerOneId, playerTwoId := registerAndCreateTwoPlayers(t, e)

	testInput := []TestCreateTeamInput{
		{
			name: "successful creation",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: db.CreateTeamWithTwoPlayersParams{
				UserID:    userId,
				PlayerOne: playerOneId,
				PlayerTwo: playerTwoId,
				Name:      "Test Name",
			},
		},
	}
	for _, data := range testInput {
		t.Run("create team "+data.name, func(t *testing.T) {
			encodedData, err := json.Marshal(data.input)
			assert.NoError(t, err, "Problem with encoding the team")

			err, rec, _ := DummyRequest(
				t,
				e,
				http.MethodPost,
				"/api/teams",
				string(encodedData),
				teamRouter.CreateTeam,
				"",
			)
			if data.error.IsError {
			} else {
				if assert.NoError(t, err, "Problem with adding new team") {
					team := new(db.Team)
					err = json.Unmarshal(rec.Body.Bytes(), team)
					assert.NoError(t, err, "Couldn't decode returned team")

					assert.Equal(t, team.PlayerOne, data.input.PlayerOne)
					assert.Equal(t, team.PlayerTwo, data.input.PlayerTwo)
					assert.Equal(t, team.Name, data.input.Name)
				}
			}
		})
	}
	_, err := userHandler.DeleteUserById(context.Background(), userId)
	assert.NoError(t, err)
}

type TestUpdateTeamInput struct {
	name  string
	error TestError
	input db.UpdateTeamByIdParams
}

func TestUpdateTeam(t *testing.T) {
	e := echo.New()
	user := DummyUser(t, e)
	userId := user.ID

	teamIn, _, _ := DummyTeam(t, e, userId)

	inputPlayerOne := db.CreateNewTeamWithOnePlayerParams{
		FirstName: "Leonard",
		LastName:  "Hopp",
		Name:      "",
		UserID:    userId,
	}
	newPlayerOne := addNewPlayer(t, e, inputPlayerOne)

	inputPlayerTwo := db.CreateNewTeamWithOnePlayerParams{
		FirstName: "Dongs",
		LastName:  "Dings",
		Name:      "",
		UserID:    userId,
	}
	newPlayerTwo := addNewPlayer(t, e, inputPlayerTwo)

	testInput := []TestUpdateTeamInput{
		{
			name: "successful update",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: db.UpdateTeamByIdParams{
				ID:        teamIn.ID,
				PlayerOne: newPlayerOne.ID,
				PlayerTwo: &newPlayerTwo.ID,
				Name:      "change name",
			},
		},
	}
	for _, data := range testInput {
		t.Run("update team "+data.name, func(t *testing.T) {
			encodedData, err := json.Marshal(data.input)
			assert.NoError(t, err, "Problem with encoding the team update")

			err, rec, _ := DummyRequest(
				t,
				e,
				http.MethodPost,
				"/api/teams/:id",
				string(encodedData),
				teamRouter.UpdateTeamById,
				teamIn.ID.String(),
			)
			if data.error.IsError {
			} else {
				if assert.NoError(t, err, "Problem with updating team") {
					team := new(db.Team)
					err = json.Unmarshal(rec.Body.Bytes(), team)
					assert.NoError(t, err, "Couldn't decode returned team")

					assert.Equal(t, teamIn.ID, team.ID)
					assert.NotEqual(t, teamIn, team)
				}
			}
		})
	}
	_, err := userHandler.DeleteUserById(context.Background(), userId)
	assert.NoError(t, err)
}

func registerAndCreateTwoPlayers(t *testing.T, e *echo.Echo) (uuid.UUID, uuid.UUID, *uuid.UUID) {

	user := DummyUser(t, e)
	playerOne := DummyPlayer(t, e, user.ID)
	playerTwo := DummyPlayerTwo(t, e, user.ID)
	return user.ID, playerOne.ID, &playerTwo.ID
}

func DummyTeam(t *testing.T, e *echo.Echo, userId uuid.UUID) (db.Team, db.Player, db.Player) {
	playerOne := DummyPlayer(t, e, userId)
	playerTwo := DummyPlayerTwo(t, e, userId)

	input := db.CreateTeamWithTwoPlayersParams{
		PlayerOne: playerOne.ID,
		PlayerTwo: &playerTwo.ID,
		Name:      "Test Team one",
		UserID:    userId,
	}
	team := addNewteam(t, e, input)
	return team, playerOne, playerTwo
}

func addNewteam(t *testing.T, e *echo.Echo, data db.CreateTeamWithTwoPlayersParams) db.Team {
	encodedData, err := json.Marshal(data)
	assert.NoError(t, err, "Problem with encoding the team")

	err, rec, _ := DummyRequest(
		t,
		e,
		http.MethodPost,
		"/api/teams",
		string(encodedData),
		teamRouter.CreateTeam,
		"",
	)
	assert.NoError(t, err, "Problem with adding new team")

	team := new(db.Team)
	err = json.Unmarshal(rec.Body.Bytes(), team)
	assert.NoError(t, err, "Couldn't decode returned team")

	return *team
}
