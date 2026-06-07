package app

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/routers"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AppInterface interface {
	Run()
}

type App struct {
	ServerConfig *config.ServerConfig
	DbConfig     *config.DbConfig
	Storage      *repositories.Storage
}

func (app *App) Run() {
	// setup logger
	logger := config.GetLogger(app.ServerConfig.AppEnv)

	// validate the server config
	validate := validator.New()

	if validationErr := validate.Struct(app.ServerConfig); validationErr != nil {
		logger.Error("Failed while validating the server config",
			zap.String("error", validationErr.Error()))
	}

	// validate the db config
	if validationErr := validate.Struct(app.DbConfig); validationErr != nil {
		logger.Error("Failed while validating the db config",
			zap.String("error", validationErr.Error()))
	}

	// create the server instance
	server := &http.Server{
		Addr:         app.ServerConfig.Addr,
		ReadTimeout:  app.ServerConfig.ReadTimeout,
		WriteTimeout: app.ServerConfig.WriteTimeout,
		IdleTimeout:  app.ServerConfig.IdleTimeout,
		Handler:      routers.RegisterRouters(),
	}

	// start the server
	logger.Info("Starting the server...",
		zap.String("port", app.ServerConfig.Addr))

	err := server.ListenAndServe()

	if err != nil {
		logger.Error("Something went wrong while starting server")
	}
}
