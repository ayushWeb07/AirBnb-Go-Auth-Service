package repositories

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type UserRepositoryInterface interface {
	CreateUser(userPayload *dtos.CreateUserPayload) *utils.AppError
	GetAllUsers() ([]*models.UserModel, *utils.AppError)
	GetUserById(userParams *dtos.GetUserByIdParams) (*models.UserModel, *utils.AppError)
	DeleteUserById(userParams *dtos.DeleteUserByIdParams) *utils.AppError
	GetUserByUsernameAndEmail(userPayload *dtos.GetUserByUsernameAndEmailPayload) (*models.UserModel, *utils.AppError)
}

type UserRepository struct {
	db           *sql.DB
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (ur *UserRepository) CreateUser(userPayload *dtos.CreateUserPayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	result, queryExecErr := ur.db.Exec(query, userPayload.Username, userPayload.Email, userPayload.Password)

	if queryExecErr != nil {
		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to insert user into the database: " + queryExecErr.Error())
	}

	id, insertErr := result.LastInsertId()

	if insertErr != nil {
		ur.logger.Error("Failed to insert user into the database",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to insert user into the database: " + insertErr.Error())
	}

	ur.logger.Info("Successfully inserted user into the database",
		zap.Int64("user_id", id))

	return nil
}

func (ur *UserRepository) GetAllUsers() ([]*models.UserModel, *utils.AppError) {
	// create the dummy instance
	var userModels []*models.UserModel

	// load the rows
	query := "SELECT id, username, email FROM users"
	rows, queryErr := ur.db.Query(query)

	if queryErr != nil {
		ur.logger.Error("Something went wrong while fetching all the users",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the users: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		userModel := &models.UserModel{}

		rowScanErr := rows.Scan(&userModel.ID, &userModel.Username, &userModel.Email)

		if rowScanErr != nil {
			ur.logger.Error("Failed to fetch all the users from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching all the users: " + rowScanErr.Error())
		}

		userModels = append(userModels, userModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		ur.logger.Error("Failed to fetch all the users from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the users: " + rowsErr.Error())
	}

	ur.logger.Info("Successfully fetched all the users from the database",
		zap.Int("count", len(userModels)))

	return userModels, nil
}

func (ur *UserRepository) GetUserById(userParams *dtos.GetUserByIdParams) (*models.UserModel, *utils.AppError) {
	// create the dummy instance
	userModel := &models.UserModel{}

	// fetch from the db
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?"

	queryErr := ur.db.QueryRow(query, userParams.ID).Scan(&userModel.ID, &userModel.Username, &userModel.Email, &userModel.CreatedAt, &userModel.UpdatedAt)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			ur.logger.Error("Such user not found",
				zap.Int("user_id", userParams.ID))

			return nil, utils.NotFound("User with such id not found")
		}

		ur.logger.Error("Failed to fetch the user from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the user from the database: " + queryErr.Error())
	}

	ur.logger.Info("Successfully fetched the user from the database",
		zap.Int("user_id", userModel.ID),
	)

	return userModel, nil
}

func (ur *UserRepository) GetUserByUsernameAndEmail(userPayload *dtos.GetUserByUsernameAndEmailPayload) (*models.UserModel, *utils.AppError) {
	existingUserModel := &models.UserModel{}

	// fetch from the db
	query := "SELECT id, username, email, password FROM users WHERE username = ? AND email = ?"

	queryErr := ur.db.QueryRow(query, userPayload.Username, userPayload.Email).Scan(&existingUserModel.ID, &existingUserModel.Username, &existingUserModel.Email, &existingUserModel.Password)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			ur.logger.Error("No such user found in the database",
				zap.String("error", queryErr.Error()))

			return nil, utils.NotFound("No such user found in the database")
		}

		ur.logger.Error("Failed to fetch the user from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the user from the database: " + queryErr.Error())
	}

	return existingUserModel, nil
}

func (ur *UserRepository) DeleteUserById(userParams *dtos.DeleteUserByIdParams) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM users WHERE id = ?"
	result, queryExecErr := ur.db.Exec(query, userParams.ID)

	if queryExecErr != nil {
		ur.logger.Error("Failed to delete user from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to delete user from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		ur.logger.Error("Failed to delete user from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to delete user from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		ur.logger.Error("No user has been deleted from the database",
			zap.Int("user_id", userParams.ID))

		return utils.NotFound("User with such id not found")
	}

	ur.logger.Info("Successfully deleted the user from the database",
		zap.Int("user_id", userParams.ID))

	return nil
}

func NewUserRepository(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) UserRepositoryInterface {
	newUserRepository := &UserRepository{
		db:           db,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newUserRepository
}
