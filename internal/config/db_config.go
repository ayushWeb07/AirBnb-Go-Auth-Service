package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	DbUsername string `validate:"required"`
	DbPassword string `validate:"required"`
	DbNet      string `validate:"required"`
	DbAddress  string `validate:"required"`
	DbName     string `validate:"required"`
}

func LoadDbConfig() (*DbConfig, error) {
	// load the env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Something went wrong while loading the env vars:", err.Error())
		return nil, err
	}

	// load the envs & create the config instance
	cfg := &DbConfig{
		DbUsername: LoadSingleEnvVar("DB_USERNAME", "admin"),
		DbPassword: LoadSingleEnvVar("DB_PASSWORD", "admin"),
		DbNet:      LoadSingleEnvVar("DB_NET", "tcp"),
		DbAddress:  LoadSingleEnvVar("DB_ADDRESS", "127.0.0.1:3306"),
		DbName:     LoadSingleEnvVar("DB_NAME", "dev_db"),
	}

	return cfg, nil
}
