package tests

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/gofrs/uuid"
	"golang.org/x/exp/slices"
)
type DBQueriesMock struct {
  users []db.User
}

func NewDBQueriesMock () *DBQueriesMock{
  return &DBQueriesMock{
    users: []db.User{},
  }
}

func (d *DBQueriesMock) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
  id, err := uuid.NewV4()
  if err != nil {
    log.Fatalf("could not create uuid: %v", err)
  }
  newUser := db.User {
    ID: id,
    Username: arg.Username,
    Email: arg.Email,
    PasswordHash: arg.PasswordHash,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  d.users = append(d.users, newUser)

  return newUser, nil
}

func (d *DBQueriesMock) GetAllUsers(ctx context.Context) ([]db.User, error) {
  return d.users, nil
}

func (d *DBQueriesMock) GetUserById(ctx context.Context, id uuid.UUID) (db.GetUserByIdRow, error) {
  idx := slices.IndexFunc(d.users, func(u db.User) bool {return u.ID == id})
  if idx == -1 {
    return db.GetUserByIdRow{}, errors.New("No User found")
  }
  
  user := db.GetUserByIdRow{
    ID: d.users[idx].ID,
    Column2: "?",
    Email: d.users[idx].Email,
  }
  return user, nil
}
