package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/gofrs/uuid"
	"golang.org/x/exp/slices"
)

func (d *DBQueriesMock) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("could not create uuid: %v", err)
	}
	newUser := db.User{
		ID:           id,
		Username:     arg.Username,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	d.users = append(d.users, newUser)

	return newUser, nil
}

func (d *DBQueriesMock) GetAllUsers(ctx context.Context) ([]db.User, error) {
	return d.users, nil
}

func (d *DBQueriesMock) GetUserById(ctx context.Context, id uuid.UUID) (db.User, error) {
	idx := slices.IndexFunc(d.users, func(u db.User) bool { return u.ID == id })
	if idx == -1 {
		return db.User{}, errors.New("User not found")
	}
	user := d.users[idx]
	return user, nil
}

func (d *DBQueriesMock) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	idx := slices.IndexFunc(d.users, func(u db.User) bool { return u.Email == email })
	if idx == -1 {
		return db.User{}, errors.New("User not found")
	}
	user := d.users[idx]
	return user, nil
}

func (d *DBQueriesMock) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	idx := slices.IndexFunc(d.users, func(u db.User) bool { return u.Username == username })
	if idx == -1 {
		return db.User{}, errors.New("User not found")
	}
	user := d.users[idx]
	return user, nil
}

func (d *DBQueriesMock) DeleteUserById(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := d.GetUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	var tempUsers []db.User

	for _, item := range d.users {
		if item.ID != id {
			tempUsers = append(tempUsers, item)
		}
	}

	d.users = tempUsers

	return user, nil
}

func (d *DBQueriesMock) UpdateUserById(ctx context.Context, args db.UpdateUserByIdParams) (db.User, error) {
	id := args.ID
	user, err := d.GetUserById(ctx, id)
	if err != nil {
		return db.User{}, err
	}

	user.Email = args.Email
	user.Username = args.Username
	user.PasswordHash = args.PasswordHash

	for idx, item := range d.users {
		if item.ID == id {
			d.users[idx] = user
		}
	}

	return user, nil
}
