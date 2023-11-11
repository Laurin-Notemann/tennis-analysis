package utils

import "github.com/Laurin-Notemann/tennis-analysis/db"

type DBQueriesMock struct {
	users []db.User
  tokens []db.RefreshToken
}

func NewDBQueriesMock() *DBQueriesMock {
	return &DBQueriesMock{
		users: []db.User{},
    tokens: []db.RefreshToken{},
	}
}
