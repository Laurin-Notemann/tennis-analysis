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
			name: "",
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

			}
		}
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
