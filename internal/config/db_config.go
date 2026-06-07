package config

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type DbConfig struct {
	DbUsername        string `validate:"required"`
	DbPassword        string `validate:"required"`
	DbNet             string `validate:"required"`
	DbAddress         string `validate:"required"`
	DbName            string `validate:"required"`
	GooseDriver       string `validate:"required"`
	GooseDbString     string `validate:"required"`
	GooseMigrationDir string `validate:"required"`
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
		DbUsername:        LoadSingleEnvVar("DB_USERNAME", "admin"),
		DbPassword:        LoadSingleEnvVar("DB_PASSWORD", "admin"),
		DbNet:             LoadSingleEnvVar("DB_NET", "tcp"),
		DbAddress:         LoadSingleEnvVar("DB_ADDRESS", "127.0.0.1:3306"),
		DbName:            LoadSingleEnvVar("DB_NAME", "dev_db"),
		GooseDriver:       LoadSingleEnvVar("GOOSE_DRIVER", "mysql"),
		GooseDbString:     LoadSingleEnvVar("GOOSE_DBSTRING", "admin@admin@tcp(127.0.0.1:3306)/dev_db"),
		GooseMigrationDir: LoadSingleEnvVar("GOOSE_MIGRATION_DIR", "internal/database/migrations"),
	}

	return cfg, nil
}

func SetupDB(dbConfig *DbConfig, logger *zap.Logger) (*sql.DB, error) {
	// capture connection properties
	cfg := mysql.NewConfig()
	cfg.User = dbConfig.DbUsername
	cfg.Passwd = dbConfig.DbPassword
	cfg.Net = dbConfig.DbNet
	cfg.Addr = dbConfig.DbAddress
	cfg.DBName = dbConfig.DbName

	// open a database connection
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		logger.Fatal("Something went wrong while opening database connection",
			zap.String("error", err.Error()))

		return nil, err
	}

	// ping the database
	pingErr := db.Ping()
	if pingErr != nil {
		logger.Fatal("Something went wrong while pinging database",
			zap.String("error", pingErr.Error()))

		return nil, pingErr
	}

	logger.Info("Successfully connected to the DB",
		zap.String("db_name", dbConfig.DbName))

	return db, nil
}
