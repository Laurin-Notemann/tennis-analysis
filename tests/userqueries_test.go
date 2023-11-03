package tests

import (
	"context"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/gofrs/uuid"
)

type UserMock struct {
	user       db.User
	createUser db.CreateUserParams
	userById   db.GetUserByIdRow
  id uuid.UUID
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
		userById: db.GetUserByIdRow{},
	}

	t.Run("CreateUser", func(t *testing.T) {
		user, err := testDbQueries.CreateUser(context.Background(), userMock.createUser)
		if err != nil {
			t.Fatalf("CreateUser on testdb couldn't create User %v", err)
		}
    userMock.id = user.ID

		userMock.user = user

		if user.Email != userMock.createUser.Email && user.PasswordHash != userMock.createUser.PasswordHash && user.Username != userMock.createUser.Username {
			t.Fatalf("testDbQueries.CreateUser(%+v) = %+v, nil, output user should have the same content as the input", userMock.createUser, user)
		}
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := testDbQueries.GetAllUsers(context.Background())
		if err != nil {
			t.Fatalf("GetAllUsers on testdb couldn't get any Users %v", err)
		}

		if len(users) != 1 && users[0].Username != userMock.user.Username {
			t.Fatalf("testDbQueries.GetAllUsers() = got %+v, nil should have at least contained: %+v", users, userMock.user)
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		user, err := testDbQueries.GetUserById(context.Background(), userMock.id)
		if err != nil {
			t.Fatalf("GetUserById on testdb couldn't get any User with id %v, %v", userMock.id ,err)
		}

    userMock.userById = user

		if user.Email == userMock.user.Email && user.ID != userMock.id{
			t.Fatalf("testDbQueries.GetUserById() = got %+v, nil should have at gotten user with Email %+v", user, userMock.user.Email)
		}

	})

}
