package services

import (
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser(userPayload *dtos.CreateUser) *utils.AppError
	LoginUser(userPayload *dtos.LoginUser) (string, *utils.AppError)
	GetAllUsers() ([]*models.UserModel, *utils.AppError)
	GetUserById(userPayload *dtos.GetUserById) (*models.UserModel, *utils.AppError)
	DeleteUserById(userPayload *dtos.DeleteUserById) *utils.AppError
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (us *UserService) CreateUser(userPayload *dtos.CreateUser) *utils.AppError {
	us.logger.Info("Create user service called...")

	// check if the user already exists
	_, repositoryErr := us.UserRepository.GetUserByUsernameAndEmail(&dtos.GetUserByUsernameAndEmail{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	})

	if repositoryErr == nil {
		return utils.BadRequest("User with such username and email, already exists")
	}

	// hash the password
	bytes, hashErr := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		us.logger.Fatal("Something went wrong while hashing the password",
			zap.String("error", hashErr.Error()))

		return utils.InternalServerError("Something went wrong while hashing the password: " + hashErr.Error())
	}

	userPayload.Password = string(bytes)

	// call the create user repository
	repositoryErr = us.UserRepository.CreateUser(userPayload)

	if repositoryErr != nil {
		return repositoryErr
	}

	us.logger.Info("Create user service was successful")

	return nil
}

func (us *UserService) LoginUser(userPayload *dtos.LoginUser) (string, *utils.AppError) {
	us.logger.Info("Login user service called...")

	// fetch the user by username and email repository
	existingUserModel, repositoryErr := us.UserRepository.GetUserByUsernameAndEmail(&dtos.GetUserByUsernameAndEmail{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	})

	if repositoryErr != nil {
		return "", repositoryErr
	}

	// check if passwords match
	compareErr := bcrypt.CompareHashAndPassword([]byte(existingUserModel.Password), []byte(userPayload.Password))

	if compareErr != nil {
		us.logger.Error("Invalid password has been provided",
			zap.String("error", compareErr.Error()))

		return "", utils.BadRequest("Invalid password has been provided")
	}

	// generate the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name":  existingUserModel.Username,
		"user_email": existingUserModel.Email,
		"user_id":    existingUserModel.ID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, tokenErr := token.SignedString([]byte(us.serverConfig.JwtSecretKey))

	if tokenErr != nil {
		us.logger.Fatal("Something went wrong while generating the token",
			zap.String("error", tokenErr.Error()))

		return "", utils.InternalServerError("Something went wrong while generating the token: " + tokenErr.Error())
	}

	us.logger.Info("Login user service was successful",
		zap.String("token", tokenString))

	return tokenString, nil
}

func (us *UserService) GetAllUsers() ([]*models.UserModel, *utils.AppError) {
	us.logger.Info("Get all users service called...")

	// call the fetch all users repository
	userModels, repositoryErr := us.UserRepository.GetAllUsers()
	return userModels, repositoryErr
}

func (us *UserService) GetUserById(userPayload *dtos.GetUserById) (*models.UserModel, *utils.AppError) {
	us.logger.Info("Get by id user service called...")

	// call the fetch user by id repository
	userModel, repositoryErr := us.UserRepository.GetUserById(userPayload)
	return userModel, repositoryErr
}

func (us *UserService) DeleteUserById(userPayload *dtos.DeleteUserById) *utils.AppError {
	us.logger.Info("Delete user service called...")

	// call the delete user by id repository
	repositoryErr := us.UserRepository.DeleteUserById(userPayload)
	return repositoryErr
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newUserService
}
