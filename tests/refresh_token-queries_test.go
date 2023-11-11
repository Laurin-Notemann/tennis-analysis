package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TokenTestData struct {
	action     string
	tokenInput testToken
}

type testToken struct {
	Token      string
	ExpiryDate time.Time
	UserID     uuid.UUID
}

func TestRefreshTokenQueries(t *testing.T) {
	testDbQueries := DbQueriesTest()
	user := getTestUser(t, testDbQueries)
	context := context.Background()
	token := testToken{
		UserID:     user.ID,
		ExpiryDate: time.Now(),
		Token:      "Test",
	}
	testData := []TokenTestData{
		{
			action:     "create",
			tokenInput: token,
		},
		{
			action:     "get",
			tokenInput: token,
		},
		{
			action:     "update",
			tokenInput: token,
		},
		{
			action:     "delete",
			tokenInput: token,
		},
	}

	for _, val := range testData {

		t.Run(val.action, func(t *testing.T) {
			if val.action == "create" {
				tokenParams := getCreateTokenParams(val.tokenInput)

				tokenUser, err := testDbQueries.CreateToken(context, tokenParams)
				if assert.NoError(t, err) {
					updatedUser, err := testDbQueries.GetUserById(context, user.ID)
					assert.NoError(t, err)

					assert.Equal(t, user.ID, tokenUser.ID)
					assert.Equal(t, user.Username, tokenUser.Username)
					assert.Equal(t, user.Email, tokenUser.Email)
					assert.Equal(t, user.PasswordHash, tokenUser.PasswordHash)
					assert.Equal(t, updatedUser, tokenUser)
					assert.NoError(t, err)
				}
			}
		})
	}
	_, err := testDbQueries.DeleteUserById(context, user.ID)
	assert.NoError(t, err)
}

func executeTests(t *testing.T, user db.User, tokenUser db.User) {

}

func getTestUser(t *testing.T, d *db.Queries) db.User {
	t.Helper()

	testUser := db.CreateUserParams{
		Username:     "Max",
		Email:        "max@mustermann.dev",
		PasswordHash: "$2a$12$FCKggXka2kwzUFw4YVfks.cJy3SAyKe.dXvxI.kZTWuFXvtYthBwW",
	}
	user, err := d.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	return user
}

func getCreateTokenParams(input testToken) db.CreateTokenParams {
	return db.CreateTokenParams{
		Token:      input.Token,
		ExpiryDate: input.ExpiryDate,
		UserID:     input.UserID,
	}
}
