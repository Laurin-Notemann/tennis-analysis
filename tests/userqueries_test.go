package tests

import (
	"context"
	"testing"

	"github.com/Laurin-Notemann/tennis-analysis/db"
)

func TestCreateUserQuery (t *testing.T) {
  testDbQueries := DbQueriesTest()
  createUser := db.CreateUserParams{
    Username: "laurin",
    Email: "laurin@test.de",
    PasswordHash: "Test",
  }
  user, err := testDbQueries.CreateUser(context.Background(), createUser)
  if err != nil {
    t.Fatalf("CreateUser on testdb couldn't create User %v", err)
  }

  users, err := testDbQueries.GetAllUsers(context.Background())
  if err != nil {
    t.Fatalf("GetAllUsers on testdb couldn't get any Users %v", err)
  }

  if len(users) != 1 && users[0].Username != user.Username && user.Username != createUser.Username {
    t.Fatalf("testDbQueries.CreateUser(%+v) = %+v, nil want to create the correct user",createUser, user)
  }
}
