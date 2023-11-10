package tests

import (
	"context"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/gofrs/uuid"
)

type UserMock struct {
	user           db.User
	createUser     db.CreateUserParams
	id             uuid.UUID
	updateUserArgs db.UpdateUserByIdParams
}

func TestUserDbQueries(t *testing.T) {
	testDbQueries := DbQueriesTest()
	userMock := UserMock{
		user: db.User{},
		createUser: db.CreateUserParams{
			Username:     "laurin",
			Email:        "laurin@test.de",
			PasswordHash: "Test",
		},
		updateUserArgs: db.UpdateUserByIdParams{
			Username:     "max",
			Email:        "max@test.de",
			PasswordHash: "Max",
		},
	}

	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		user, err := testDbQueries.CreateUser(ctx, userMock.createUser)
		if err != nil {
			t.Fatalf("CreateUser on testdb couldn't create User %v", err)
		}
		userMock.id = user.ID

		userMock.user = user

		userMock.updateUserArgs.ID = user.ID

		if user.Email != userMock.createUser.Email && user.PasswordHash != userMock.createUser.PasswordHash && user.Username != userMock.createUser.Username {
			t.Fatalf("testDbQueries.CreateUser(%+v) = %+v, nil, output user should have the same content as the input", userMock.createUser, user)
		}
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := testDbQueries.GetAllUsers(ctx)
		if err != nil {
			t.Fatalf("GetAllUsers on testdb couldn't get any Users %v", err)
		}

		if len(users) != 1 && users[0].Username != userMock.user.Username {
			t.Fatalf("testDbQueries.GetAllUsers() = got %+v, nil should have at least contained: %+v", users, userMock.user)
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		user, err := testDbQueries.GetUserById(ctx, userMock.id)
		if err != nil {
			t.Fatalf("GetUserById on testdb couldn't get any User with id %v, %v", userMock.id, err)
		}

		if user.Email != userMock.user.Email && user.ID != userMock.id {
			t.Fatalf("testDbQueries.GetUserById() = got %+v, nil should have at gotten user with Email %+v", user, userMock.user.Email)
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		email := "laurin@test.de"
		user, err := testDbQueries.GetUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("GetUserByEmail() on testdb couldn't get any User with email %v, %v", email, err)
		}

		if user.Email != userMock.user.Email && user.ID != userMock.id && user.Username != userMock.user.Username {
			t.Fatalf("testDbQueries.GetUserByEmail() = got %+v, nil should have at gotten user with email %+v", user, email)
		}
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		username := "laurin"
		user, err := testDbQueries.GetUserByUsername(ctx, username)
		if err != nil {
			t.Fatalf("GetUserByUsername() on testdb couldn't get any User with username %v, %v", username, err)
		}

		if user.Email != userMock.user.Email && user.ID != userMock.id && user.Username != userMock.user.Username {
			t.Fatalf("testDbQueries.GetUserByUsername() = got %+v, nil should have at gotten user with username %+v", user, username)
		}
	})

	t.Run("UpdateUserById", func(t *testing.T) {
		user, err := testDbQueries.UpdateUserById(ctx, userMock.updateUserArgs)
		if err != nil {
			t.Fatalf("UpdateUserById() on testdb couldn't update any user with id %v, %v", userMock.id, err)
		}

		userMock.user = user

		if user.Email != userMock.updateUserArgs.Email && user.ID != userMock.id && user.Username != userMock.updateUserArgs.Username {
			t.Fatalf("testDbQueries.UpdateUserById() = got %+v, nil should have at gotten user with Email %+v", user, userMock.user.Email)
		}
	})

	t.Run("DeleteUserById", func(t *testing.T) {
		user, err := testDbQueries.DeleteUserById(ctx, userMock.id)
		if err != nil {
			t.Fatalf("DeleteUserById() on testdb couldn't delete any User with id %v, %v", userMock.id, err)
		}

		users, err := testDbQueries.GetAllUsers(ctx)
		if err != nil {
			t.Fatalf("GetAllUsers() couldn't get any users: %v", err)
		}

		if len(users) > 0 {
			t.Fatalf("testDbQueries.DeleteUserById() = %+v, nil, want empty users slice, got %+v", user, users)
		}
	})

}
