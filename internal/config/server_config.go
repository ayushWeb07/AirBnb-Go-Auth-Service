package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Addr         string        `validate:"required"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
	IdleTimeout  time.Duration `validate:"required"`
	AppEnv       string        `validate:"required"`
	JwtSecretKey string        `validate:"required"`
}

func LoadServerConfig() (*ServerConfig, error) {
	// load the env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Something went wrong while loading the env vars:", err.Error())
		return nil, err
	}

	// load the envs & create the config instance
	cfg := &ServerConfig{
		Addr:         LoadSingleEnvVar("ADDR", ":8181"),
		ReadTimeout:  time.Duration(LoadSingleEnvVar("READ_TIMEOUT", 15)) * time.Second,
		WriteTimeout: time.Duration(LoadSingleEnvVar("WRITE_TIMEOUT", 15)) * time.Second,
		IdleTimeout:  time.Duration(LoadSingleEnvVar("IDLE_TIMEOUT", 180)) * time.Second,
		AppEnv:       LoadSingleEnvVar("APP_ENV", "development"),
		JwtSecretKey: LoadSingleEnvVar("JWT_SECRET_KEY", ""),
	}

	return cfg, nil
}
