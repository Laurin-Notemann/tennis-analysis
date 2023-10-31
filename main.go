package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/Laurin-Notemann/tennis-analysis/config"
)

func main() {
	// Parse Config
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v\n", err)
	}

	var cfg config.Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

  log.Printf("env: %v\n", cfg.DB.Url)

}
