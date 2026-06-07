package main

import (
	"log"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
)
import "github.com/ayushWeb07/AirBnb-Go-Api-Gateway/cmd/app"

func main() {
	// create the server config instance
	serverCfg, err := config.LoadServerConfig()

	if err != nil {
		log.Fatal(err)
	}

	// create the db config instance
	dbCfg, err := config.LoadDbConfig()

	if err != nil {
		log.Fatal(err)
	}

	// create the storage instance
	storage := repositories.NewStorage()

	// create the server instance
	serverApp := &app.App{
		ServerConfig: serverCfg,
		DbConfig:     dbCfg,
		Storage:      storage,
	}

	// run the app
	serverApp.Run()
}
