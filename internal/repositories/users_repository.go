package repositories

import (
	"database/sql"
	"fmt"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"go.uber.org/zap"
)

type UserRepositoryInterface interface {
	GetAllUsers() ([]*models.UserModel, error)
	GetUserById() (*models.UserModel, error)
	CreateUser(userModel *models.UserModel) error
	DeleteUserById() error
}

type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func (ur *UserRepository) GetAllUsers() ([]*models.UserModel, error) {
	// create the dummy instance
	var userModels []*models.UserModel

	// load the rows
	query := "SELECT id, username, email FROM users"
	rows, err := ur.db.Query(query)

	if err != nil {
		ur.logger.Error("Something went wrong while fetching all the users",
			zap.String("error", err.Error()))

		return nil, err
	}
	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		userModel := &models.UserModel{}

		if err := rows.Scan(&userModel.ID, &userModel.Username, &userModel.Email); err != nil {
			ur.logger.Error("Failed to fetch all the users from the database",
				zap.String("error", err.Error()))

			return nil, err
		}

		userModels = append(userModels, userModel)
	}

	if err := rows.Err(); err != nil {
		ur.logger.Error("Failed to fetch all the users from the database",
			zap.String("error", err.Error()))

		return nil, err
	}

	ur.logger.Info("Successfully fetched all the users from the database",
		zap.Int("count", len(userModels)))

	return userModels, nil
}

func (ur *UserRepository) GetUserById() (*models.UserModel, error) {
	// create the dummy instance
	userModel := &models.UserModel{}

	// fetch from the db
	query := "SELECT id, username, email FROM users WHERE id = ?"

	if err := ur.db.QueryRow(query, 1).Scan(&userModel.ID, &userModel.Username, &userModel.Email); err != nil {
		if err == sql.ErrNoRows {
			ur.logger.Error("Failed to fetch the user from the database",
				zap.String("error", err.Error()))

			return nil, err
		}

		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", err.Error()))

		return nil, err
	}

	ur.logger.Info("Successfully fetched the user from the database",
		zap.String("user_id", userModel.ID),
		zap.String("user_username", userModel.Username),
		zap.String("user_email", userModel.Email),
	)

	return userModel, nil
}

func (ur *UserRepository) CreateUser(userModel *models.UserModel) error {
	// insert into the db
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	result, err := ur.db.Exec(query, userModel.Username, userModel.Email, userModel.Password)

	if err != nil {
		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", err.Error()))

		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", err.Error()))

		return err
	}

	ur.logger.Info("Successfully inserted user into the database",
		zap.Int64("user_id", id))

	return nil
}

func (ur *UserRepository) DeleteUserById() error {
	// prepare and execute the query
	query := "DELETE FROM users WHERE id = ?"
	result, err := ur.db.Exec(query, 1)

	if err != nil {
		ur.logger.Error("Failed to delete user from the database",
			zap.String("error", err.Error()))

		return err
	}

	// check if any rows got affected
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		ur.logger.Error("Failed to delete user from the database",
			zap.String("error", err.Error()))

		return err
	}

	if rowsAffected == 0 {
		ur.logger.Error("No user has been deleted from the database",
			zap.Int("user_id", 1))

		return fmt.Errorf("No user has been deleted from the database")
	}

	ur.logger.Info("Successfully deleted the user from the database",
		zap.Int("user_id", 1))

	return nil
}

func NewUserRepository(logger *zap.Logger, db *sql.DB) UserRepositoryInterface {
	newUserRepository := &UserRepository{
		db:     db,
		logger: logger,
	}

	return newUserRepository
}
