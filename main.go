package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"

	"github.com/Laurin-Notemann/tennis-analysis/api"
	"github.com/Laurin-Notemann/tennis-analysis/config"
	"github.com/Laurin-Notemann/tennis-analysis/db"
	"github.com/Laurin-Notemann/tennis-analysis/handler"
	"github.com/Laurin-Notemann/tennis-analysis/utils"
)

func main() {
	ctx := context.Background()

	// Parse Config
	err := godotenv.Load()
	if err != nil {
		log.Printf("error?? loading .env file: %v\n", err)
	}

	var cfg config.Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	log.Printf("env: %v\n", cfg.DB.Url)

	dbCon, err := sql.Open("postgres", cfg.DB.Url)
	if err != nil {
		log.Fatal(err)
	}

	err = dbCon.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := db.New(dbCon)

	tokenGen := utils.ProdTokenGenerator{}
	userHandler := handler.NewUserHandler(dbQueries, cfg)
	tokenHandler := handler.NewRefreshTokenHandler(dbQueries, cfg, &tokenGen)
	authHanlder := handler.NewAuthenticationHandler(dbQueries, *userHandler, *tokenHandler, &tokenGen)
  playerHandler := handler.NewPlayerHandler(dbQueries)
	teamHandler := handler.NewTeamHandler(dbQueries)

	resourceHandler := handler.ResourceHandlers{
		UserHandler:  *userHandler,
		TokenHandler: *tokenHandler,
		AuthHandler:  *authHanlder,
    PlayerHandler: *playerHandler,
    TeamHandler: *teamHandler,
	}

	server := api.NewApi(ctx, resourceHandler, &tokenGen)

  echoPort := cfg.ECHO.Port
  echoHost := cfg.ECHO.Host

  echoString := echoHost +":"+ fmt.Sprint(echoPort)

	err = server.Start(echoString)
	if err != nil {
		log.Fatal(err)
	}

	dbCon.Close()

}
