package tests

import (
	"context"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
)

func TestUserHandler(t *testing.T) {
	newUser := db.CreateUserParams{
		Username:     "laurin",
		Email:        "laurin@test.de",
		PasswordHash: "Test",
	}
	dbMock := NewDBQueriesMock()
	userHandler := handler.UserHandler{
		DB: dbMock,
	}

	correctUser := db.User{
		Username:     "laurin",
		Email:        "laurin@test.de",
		PasswordHash: "Test",
	}

	t.Run("CreateUser", func(t *testing.T) {
		user, err := userHandler.CreateUser(context.Background(), newUser)
		if err != nil {
			t.Fatalf("userHanlder.CreateUser(%+v) = nil, %v, want match for correct User", newUser, err)
		}

		correctUser.ID = user.ID
		correctUser.CreatedAt = user.CreatedAt
		correctUser.UpdatedAt = user.UpdatedAt

		if user != correctUser {
			t.Fatalf("userHanlder.CreateUser(%+v) = %+v, nil, want match for correct %+v, nil", newUser, user, correctUser)
		}
	})

	t.Run("GetAllUsers", func(t *testing.T) {
		users, err := userHandler.GetAllUsers(context.Background())
		if err != nil {
			t.Fatalf("userHandler.GetAllUsers(), err when trying to get users: %v", err)
		}
		if len(users) != 1 {
			t.Fatalf("userHandler.GetAllUsers(), want list of users with one entry, got %#v", users)
		}
	})

	t.Run("GetUserById", func(t *testing.T) {
		user, err := userHandler.GetUserById(context.Background(), correctUser.ID)
		if err != nil {
			t.Fatalf("userhandler.GetUserById(), err when tryong to get user by id: %v", err)
		}
		if user.Email != correctUser.Email && user.ID != correctUser.ID && correctUser.Username != user.Username {
			t.Fatalf("userHandler.GetAllUsers(), want %#v, got %#v", correctUser, user)
		}
	})

}
