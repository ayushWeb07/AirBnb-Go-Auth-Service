package config

import (
	"time"
)

type AppEnvType string

const (
	Development AppEnvType = "development"
	Production  AppEnvType = "production"
)

type ServerConfig struct {
	Addr         string        `validate:"required"`
	ReadTimeout  time.Duration `validate:"required"`
	WriteTimeout time.Duration `validate:"required"`
	IdleTimeout  time.Duration `validate:"required"`
	AppEnv       AppEnvType    `validate:"required"`
}
