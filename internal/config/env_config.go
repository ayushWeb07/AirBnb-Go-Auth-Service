package config

import (
	"fmt"
	"os"
	"strconv"
)

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
