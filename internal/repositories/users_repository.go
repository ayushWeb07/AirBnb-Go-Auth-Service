package repositories

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"go.uber.org/zap"
)

type UserRepositoryInterface interface {
	CreateUser()
	GetUserById()
	GetAllUsers()
}

type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func (ur *UserRepository) GetAllUsers() {
	// create the dummy instance
	var userModels []*models.UserModel

	// load the rows
	query := "SELECT id, username, email FROM users"
	rows, err := ur.db.Query(query)

	if err != nil {
		ur.logger.Error("Something went wrong while fetching all the users",
			zap.String("error", err.Error()))
	}
	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		userModel := &models.UserModel{}

		if err := rows.Scan(&userModel.ID, &userModel.Username, &userModel.Email); err != nil {
			ur.logger.Error("Failed to fetch all the users from the database",
				zap.String("error", err.Error()))

			return
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		ur.logger.Error("Failed to fetch all the users from the database",
			zap.String("error", err.Error()))

		return
	}

	ur.logger.Info("Successfully fetched all the users from the database",
		zap.Int("count", len(userModels)))
}

func (ur *UserRepository) GetUserById() {
	// create the dummy instance
	userModel := &models.UserModel{}

	// fetch from the db
	query := "SELECT id, username, email FROM users WHERE id = ?"

	if err := ur.db.QueryRow(query, 1).Scan(&userModel.ID, &userModel.Username, &userModel.Email); err != nil {
		if err == sql.ErrNoRows {
			ur.logger.Error("Failed to fetch the user from the database",
				zap.String("error", err.Error()))

			return
		}

		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", err.Error()))

		return
	}

	ur.logger.Info("Successfully fetched the user from the database",
		zap.String("user_id", userModel.ID),
		zap.String("user_username", userModel.Username),
		zap.String("user_email", userModel.Email),
	)
}

func (ur *UserRepository) CreateUser() {
	// create a dummy user model instance
	userModel := &models.UserModel{
		Username: "ayush",
		Email:    "ayush@gmail.com",
	}

	// insert into the db
	query := "INSERT INTO users (username, email) VALUES (?, ?)"
	result, err := ur.db.Exec(query, userModel.Username, userModel.Email)

	if err != nil {
		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", err.Error()))

		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", err.Error()))

		return
	}

	ur.logger.Info("Successfully inserted user into the database",
		zap.Int64("user_id", id))
}

func NewUserRepository(logger *zap.Logger, db *sql.DB) UserRepositoryInterface {
	newUserRepository := &UserRepository{
		db:     db,
		logger: logger,
	}

	return newUserRepository
}
