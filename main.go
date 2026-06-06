package main

import (
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
)
import "github.com/ayushWeb07/AirBnb-Go-Api-Gateway/cmd/app"

func main() {
	// create the config instance
	cfg := &config.ServerConfig{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		AppEnv:       config.Development,
	}

	serverApp := &app.App{
		ServerConfig: cfg,
	}

	// run the app
	serverApp.Run()
}
