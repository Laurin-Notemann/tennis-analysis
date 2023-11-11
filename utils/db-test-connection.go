package utils

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func DbQueriesTest() *db.Queries {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("error loading .env file: %v\n", err)
	}

	var cfg config.Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	dbCon, err := sql.Open("postgres", cfg.DB.TestUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := db.New(dbCon)

	return dbQueries
}
