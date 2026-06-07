package config

import (
	"fmt"

	"go.uber.org/zap"
)

func GetLogger(appEnv string) *zap.Logger {
	// setup logger
	logger := zap.Must(zap.NewProduction())
	if appEnv == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Something went wrong while syncing the zap logger:", err.Error())
		}
	}(logger)

	return logger
}
