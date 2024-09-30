package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/thoriqdharmawan/be-question-generator/cmd"
	"github.com/thoriqdharmawan/be-question-generator/config"
	"github.com/thoriqdharmawan/be-question-generator/db"
)

func main() {
	if godotenv.Load(".env") != nil {
		log.Fatal("Error loading .env file")
	}

	confVars, configErr := config.New()

	if configErr != nil {
		log.Fatal("err conf: %w", configErr)
	}

	dbErr := db.Init()

	if dbErr != nil {
		log.Fatal("err db: %w", dbErr)
	}

	app := cmd.InitApp()

	app.Listen(confVars.Port)
}
