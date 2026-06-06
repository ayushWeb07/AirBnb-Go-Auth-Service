package app

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AppInterface interface {
	Run()
}

type App struct {
	ServerConfig *config.ServerConfig
}

func (app *App) Run() {
	// setup logger
	logger := config.GetLogger(app.ServerConfig.AppEnv)

	// validate the config
	validate := validator.New()
	validationErr := validate.Struct(app.ServerConfig)
	if validationErr != nil {
		logger.Error("Failed while validating the server config",
			zap.String("error", validationErr.Error()))
	}

	// create the server instance
	server := &http.Server{
		Addr:         app.ServerConfig.Addr,
		ReadTimeout:  app.ServerConfig.ReadTimeout,
		WriteTimeout: app.ServerConfig.WriteTimeout,
		IdleTimeout:  app.ServerConfig.IdleTimeout,
		Handler:      nil,
	}

	// start the server
	logger.Info("Starting the server...",
		zap.String("port", app.ServerConfig.Addr))

	err := server.ListenAndServe()

	if err != nil {
		logger.Error("Something went wrong while starting server")
	}
}
