package handler

import (
	"context"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
)

func TestUserHandler(t *testing.T) {
	newUser := CreateUserInput{
		Username: "laurin",
		Email:    "laurin@test.de",
		Password: "Test",
	}
	dbMock := utils.NewDBQueriesMock()
	userHandler := UserHandler{
		DB: dbMock,
	}

	correctUser := db.User{
		Username: "laurin",
		Email:    "laurin@test.de",
	}

	updateUserArgs := db.UpdateUserByIdParams{
		Username:     "max",
		Email:        "max@test.de",
		PasswordHash: "Max",
	}

	t.Run("CreateUser", func(t *testing.T) {
		user, err := userHandler.CreateUser(context.Background(), newUser)
		if err != nil {
			t.Fatalf("userHandler.CreateUser(%+v) = nil, %v, want match for correct User", newUser, err)
		}

		correctUser.ID = user.ID
		updateUserArgs.ID = user.ID
		correctUser.CreatedAt = user.CreatedAt
		correctUser.UpdatedAt = user.UpdatedAt
		correctUser.PasswordHash = user.PasswordHash

		if user != correctUser {
			t.Fatalf("userHandler.CreateUser(%+v) = %+v, nil, want match for correct %+v, nil", newUser, user, correctUser)
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
			t.Fatalf("userhandler.GetUserById(), err when trying to get user by id(%v): %v", correctUser.ID, err)
		}

		if user.Email != correctUser.Email && user.ID != correctUser.ID && correctUser.Username != user.Username {
			t.Fatalf("userHandler.GetUserById(), want %#v, got %#v", correctUser, user)
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		email := "laurin@test.de"
		user, err := userHandler.GetUserByEmail(context.Background(), email)
		if err != nil {
			t.Fatalf("userhandler.GetUserByEmail(), err when trying to get user by email(%v): %v", email, err)
		}

		if user.Email != correctUser.Email && user.ID != correctUser.ID && correctUser.Username != user.Username {
			t.Fatalf("userHandler.GetUserByEmail(), want %#v, got %#v", correctUser, user)
		}
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		username := "laurin"
		user, err := userHandler.GetUserByUsername(context.Background(), username)
		if err != nil {
			t.Fatalf("userhandler.GetUserByUsername(), err when trying to get user by email(%v): %v", username, err)
		}

		if user.Email != correctUser.Email && user.ID != correctUser.ID && correctUser.Username != user.Username {
			t.Fatalf("userHandler.GetUserByUsername(), want %#v, got %#v", correctUser, user)
		}
	})

	t.Run("UpdateUserById", func(t *testing.T) {
		user, err := userHandler.UpdateUserById(context.Background(), updateUserArgs)
		if err != nil {
			t.Fatalf("userhandler.UpdateUserById(), err when trying to update user by id(%v): %v", updateUserArgs, err)
		}

		fetchedUser, err := userHandler.GetUserById(context.Background(), updateUserArgs.ID)
		if err != nil {
			t.Fatalf("userhandler.GetUserById(), err when trying to get the updated user by id(%v): %v", updateUserArgs.ID, err)
		}

		if user.Email != fetchedUser.Email && user.ID != fetchedUser.ID && fetchedUser.Username != user.Username {
			t.Fatalf("userHandler.UpdateUserById(), want %#v, got %#v", user, fetchedUser)
		}
	})

	t.Run("DeleteUserById", func(t *testing.T) {
		_, err := userHandler.DeleteUserById(context.Background(), correctUser.ID)
		if err != nil {
			t.Fatalf("userhandler.DeleteUserById(), err when trying to delete user by id(%v): %v", correctUser.ID, err)
		}

		users, err := userHandler.GetAllUsers(context.Background())
		if err != nil {
			t.Fatalf("userhandler.GetAllUsers(), err when trying to get all users: %v", err)
		}

		if len(users) > 0 {
			t.Fatalf("userHandler.DeleteUserById(), want empty slice, got %#v", users)
		}
	})
}
