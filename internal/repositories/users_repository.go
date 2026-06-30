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
	UpdateUserById(userParams *dtos.UpdateUserByIdParams, userPayload *dtos.UpdateUserByIdPayload) *utils.AppError
	DeleteUserById(userParams *dtos.DeleteUserByIdParams) *utils.AppError
	GetUserByEmail(userPayload *dtos.GetUserByEmailPayload) (*models.UserModel, *utils.AppError)
	GetUserByUsernameAndEmail(userPayload *dtos.GetUserByUsernameAndEmailPayload) (*models.UserModel, *utils.AppError)
	CreateSession(sessionPayload *dtos.CreateSessionPayload) *utils.AppError
	FetchSession(sessionPayload *dtos.FetchSessionPayload) (*models.SessionModel, *utils.AppError)
	UpdateSessionById(sessionParams *dtos.UpdateSessionByIdParams, sessionPayload *dtos.UpdateSessionByIdPayload) *utils.AppError
	CreateOtp(otpPayload *dtos.CreateOtpRepoPayload) *utils.AppError
	FetchOtp(otpPayload *dtos.FetchOtpRepoPayload) (*models.OtpModel, *utils.AppError)
	DeleteOtps(otpPayload *dtos.DeleteOtpsRepoPayload) *utils.AppError
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
	query := "SELECT id, username, email, verified, created_at, updated_at FROM users"
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

		rowScanErr := rows.Scan(&userModel.ID, &userModel.Username, &userModel.Email, &userModel.Verified, &userModel.CreatedAt, &userModel.UpdatedAt)

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
	query := "SELECT id, username, email, verified, created_at, updated_at FROM users WHERE id = ?"

	queryErr := ur.db.QueryRow(query, userParams.ID).Scan(&userModel.ID, &userModel.Username, &userModel.Email, &userModel.Verified, &userModel.CreatedAt, &userModel.UpdatedAt)

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

func (ur *UserRepository) UpdateUserById(userParams *dtos.UpdateUserByIdParams, userPayload *dtos.UpdateUserByIdPayload) *utils.AppError {
	// prepare and execute the query
	query := "UPDATE users SET verified = ? WHERE id = ?"
	result, queryExecErr := ur.db.Exec(query, userPayload.Verified, userParams.ID)

	if queryExecErr != nil {
		ur.logger.Error("Failed to update user from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to update user from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		ur.logger.Error("Failed to update user from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to update user from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		ur.logger.Error("No user has been updated from the database",
			zap.Int("user_id", userParams.ID))

		return utils.NotFound("User with such id not found")
	}

	ur.logger.Info("Successfully updated the user from the database",
		zap.Int("user_id", userParams.ID))

	return nil
}

func (ur *UserRepository) GetUserByEmail(userPayload *dtos.GetUserByEmailPayload) (*models.UserModel, *utils.AppError) {
	existingUserModel := &models.UserModel{}

	// fetch from the db
	query := "SELECT id, username, email, password, verified, created_at, updated_at FROM users WHERE email = ?"

	queryErr := ur.db.QueryRow(query, userPayload.Email).Scan(&existingUserModel.ID, &existingUserModel.Username, &existingUserModel.Email, &existingUserModel.Password, &existingUserModel.Verified, &existingUserModel.CreatedAt, &existingUserModel.UpdatedAt)

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

func (ur *UserRepository) GetUserByUsernameAndEmail(userPayload *dtos.GetUserByUsernameAndEmailPayload) (*models.UserModel, *utils.AppError) {
	existingUserModel := &models.UserModel{}

	// fetch from the db
	query := "SELECT id, username, email, password, verified, created_at, updated_at FROM users WHERE username = ? AND email = ?"

	queryErr := ur.db.QueryRow(query, userPayload.Username, userPayload.Email).Scan(&existingUserModel.ID, &existingUserModel.Username, &existingUserModel.Email, &existingUserModel.Password, &existingUserModel.Verified, &existingUserModel.CreatedAt, &existingUserModel.UpdatedAt)

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

func (ur *UserRepository) CreateSession(sessionPayload *dtos.CreateSessionPayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO sessions (user_id, refresh_token_hash) VALUES (?, ?)"
	result, queryExecErr := ur.db.Exec(query, sessionPayload.UserID, sessionPayload.RefreshTokenHash)

	if queryExecErr != nil {
		ur.logger.Error("Failed to insert session into the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to insert session into the database: " + queryExecErr.Error())
	}

	id, insertErr := result.LastInsertId()

	if insertErr != nil {
		ur.logger.Error("Failed to insert session into the database",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to insert session into the database: " + insertErr.Error())
	}

	ur.logger.Info("Successfully inserted session into the database",
		zap.Int64("session_id", id))

	return nil
}

func (ur *UserRepository) FetchSession(sessionPayload *dtos.FetchSessionPayload) (*models.SessionModel, *utils.AppError) {
	// create the dummy instance
	sessionModel := &models.SessionModel{}

	// fetch from the db
	query := "SELECT id FROM sessions WHERE refresh_token_hash = ? AND user_id = ? AND revoked = ?"

	queryErr := ur.db.QueryRow(query, sessionPayload.RefreshTokenHash, sessionPayload.UserID, sessionPayload.Revoked).Scan(&sessionModel.ID)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			ur.logger.Error("Such session not found")

			return nil, utils.NotFound("Such session not found")
		}

		ur.logger.Error("Failed to fetch the session from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the session from the database: " + queryErr.Error())
	}

	ur.logger.Info("Successfully fetched the session from the database",
		zap.Int("session_id", sessionModel.ID),
	)

	return sessionModel, nil
}

func (ur *UserRepository) UpdateSessionById(sessionParams *dtos.UpdateSessionByIdParams, sessionPayload *dtos.UpdateSessionByIdPayload) *utils.AppError {
	// prepare and execute the query
	query := "UPDATE sessions SET revoked = ? WHERE id = ?"
	result, queryExecErr := ur.db.Exec(query, sessionPayload.Revoked, sessionParams.ID)

	if queryExecErr != nil {
		ur.logger.Error("Failed to update session from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to update session from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		ur.logger.Error("Failed to update session from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to update session from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		ur.logger.Error("No session has been updated from the database",
			zap.Int("session_id", sessionParams.ID))

		return utils.NotFound("Session with such id not found")
	}

	ur.logger.Info("Successfully updated the session from the database",
		zap.Int("session_id", sessionParams.ID))

	return nil
}

func (ur *UserRepository) CreateOtp(otpPayload *dtos.CreateOtpRepoPayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO otps (user_id, user_email, otp_hash) VALUES (?, ?, ?)"
	result, queryExecErr := ur.db.Exec(query, otpPayload.UserID, otpPayload.UserEmail, otpPayload.OtpHash)

	if queryExecErr != nil {
		ur.logger.Error("Failed to insert otp into the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to insert otp into the database: " + queryExecErr.Error())
	}

	id, insertErr := result.LastInsertId()

	if insertErr != nil {
		ur.logger.Error("Failed to insert otp into the database",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to insert otp into the database: " + insertErr.Error())
	}

	ur.logger.Info("Successfully inserted otp into the database",
		zap.Int64("otp_id", id))

	return nil
}

func (ur *UserRepository) FetchOtp(otpPayload *dtos.FetchOtpRepoPayload) (*models.OtpModel, *utils.AppError) {
	existingOtpModel := &models.OtpModel{}

	// fetch from the db
	query := "SELECT id FROM otps WHERE user_id = ? AND user_email = ? AND otp_hash = ?"

	queryErr := ur.db.QueryRow(query, otpPayload.UserID, otpPayload.UserEmail, otpPayload.OtpHash).Scan(&existingOtpModel.ID)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			ur.logger.Error("No such otp found in the database",
				zap.String("error", queryErr.Error()))

			return nil, utils.NotFound("No such otp found in the database")
		}

		ur.logger.Error("Failed to fetch the otp from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the otp from the database: " + queryErr.Error())
	}

	return existingOtpModel, nil
}

func (ur *UserRepository) DeleteOtps(otpPayload *dtos.DeleteOtpsRepoPayload) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM otps WHERE user_id = ?"
	result, queryExecErr := ur.db.Exec(query, otpPayload.UserID)

	if queryExecErr != nil {
		ur.logger.Error("Failed to delete otps from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to delete otps from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		ur.logger.Error("Failed to delete otps from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to delete otps from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		ur.logger.Error("No otp has been deleted from the database")

		return utils.NotFound("No otps associated with such user")
	}

	ur.logger.Info("Successfully deleted the otps from the database",
		zap.Int("user_id", otpPayload.UserID))

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
