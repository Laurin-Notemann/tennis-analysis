package tests

import (
	"context"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
)

func Test_CreateUser(t *testing.T) {
	newUser := db.CreateUserParams{
		Username:     "laurin",
		Email:        "laurin@test.de",
		PasswordHash: "Test",
	}
	dbMock := NewDBQueriesMock()

	userHandler := handler.UserHandler{
		DB: dbMock,
	}

	user, err := userHandler.CreateUser(context.Background(), newUser)
	if err != nil {
		t.Fatalf("userHanlder.CreateUser(%+v) = nil, %v, want match for correct User", newUser, err)
	}
	correctUser := db.User{
		ID:           user.ID,
		Username:     "laurin",
		Email:        "laurin@test.de",
		PasswordHash: "Test",
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	if user != correctUser {
		t.Fatalf("userHanlder.CreateUser(%+v) = %+v, nil, want match for correct %+v, nil", newUser, user, correctUser)
	}
}
