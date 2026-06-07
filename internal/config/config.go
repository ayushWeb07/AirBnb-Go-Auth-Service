package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Addr         string        `validate:"required"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
	IdleTimeout  time.Duration `validate:"required"`
	AppEnv       string        `validate:"required"`
}

func LoadSingleEnvVar[T string | int | float32 | bool](key string, defaultVal T) T {
	val, exists := os.LookupEnv(key)

	if exists {
		switch any(defaultVal).(type) {
		case string:
			return any(val).(T)

		case int:
			i, err := strconv.Atoi(val)

			if err != nil {
				fmt.Println("Something went wrong while parsing the env var:", err.Error())
			}

			return any(i).(T)

		case float32:
			i, err := strconv.ParseFloat(val, 32)

			if err != nil {
				fmt.Println("Something went wrong while parsing the env var:", err.Error())
			}

			return any(i).(T)

		case bool:
			i, err := strconv.ParseBool(val)

			if err != nil {
				fmt.Println("Something went wrong while parsing the env var:", err.Error())
			}

			return any(i).(T)

		default:
			fmt.Println("Unsupported env variable type:")

		}
	}

	return defaultVal
}

func LoadConfig() (*ServerConfig, error) {
	// load the env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Something went wrong while loading the env vars:", err.Error())
		return nil, err
	}

	// load the envs & create the config instances
	cfg := &ServerConfig{
		Addr:         LoadSingleEnvVar("ADDR", ":8181"),
		ReadTimeout:  time.Duration(LoadSingleEnvVar("READ_TIMEOUT", 15)) * time.Second,
		WriteTimeout: time.Duration(LoadSingleEnvVar("WRITE_TIMEOUT", 15)) * time.Second,
		IdleTimeout:  time.Duration(LoadSingleEnvVar("IDLE_TIMEOUT", 180)) * time.Second,
		AppEnv:       LoadSingleEnvVar("APP_ENV", "development"),
	}

	return cfg, nil
}
