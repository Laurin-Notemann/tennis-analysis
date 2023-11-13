package handler

import "github.com/Laurin-Notemann/tennis-analysis/db"

type PlayerHandler struct {
	DB db.Querier
}

func NewPlayerHandler(DB *db.Queries) *PlayerHandler {
	return &PlayerHandler{
		DB: DB,
	}
}


