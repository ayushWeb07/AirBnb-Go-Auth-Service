package services

import (
	"fmt"
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser() (*models.UserModel, error)
	GetUserById() (*models.UserModel, error)
	GetAllUsers() ([]*models.UserModel, error)
	DeleteUserById() error
	LoginUser() (string, error)
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (us *UserService) GetAllUsers() ([]*models.UserModel, error) {
	us.logger.Info("Get all users service called...")

	// call the fetch all users repository
	userModels, err := us.UserRepository.GetAllUsers()
	return userModels, err
}

func (us *UserService) GetUserById() (*models.UserModel, error) {
	us.logger.Info("Get by id user service called...")

	// call the fetch user by id repository
	userModel, err := us.UserRepository.GetUserById()
	return userModel, err
}

func (us *UserService) CreateUser() (*models.UserModel, error) {
	us.logger.Info("Create user service called...")

	// create a dummy user model instance
	userModel := &models.UserModel{
		Username: "ronaldo",
		Email:    "ronaldo@gmail.com",
		Password: "ronaldo@2007",
	}

	// check if the user already exists
	_, err := us.UserRepository.GetUserByUsernameAndEmail(userModel)

	if err == nil {
		return nil, fmt.Errorf("User with such username and email, already exists")
	}

	// hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)

	if err != nil {
		us.logger.Fatal("Something went wrong while hashing the password",
			zap.String("error", err.Error()))

		return nil, err
	}

	userModel.Password = string(bytes)

	// call the create user repository
	_, err = us.UserRepository.CreateUser(userModel)

	if err != nil {
		return nil, err
	}

	us.logger.Info("Create user service was successful")

	return userModel, nil
}

func (us *UserService) DeleteUserById() error {
	us.logger.Info("Delete user service called...")

	// call the delete user by id repository
	err := us.UserRepository.DeleteUserById()
	return err
}

func (us *UserService) LoginUser() (string, error) {
	us.logger.Info("Login user service called...")

	// create a dummy user model instance
	userModel := &models.UserModel{
		Username: "conor",
		Email:    "conor@gmail.com",
		Password: "aonor@2007",
	}

	// fetch the user by username and email repository
	existingUserModel, err := us.UserRepository.GetUserByUsernameAndEmail(userModel)

	if err != nil {
		return "", err
	}

	// check if passwords match
	err = bcrypt.CompareHashAndPassword([]byte(existingUserModel.Password), []byte(userModel.Password))

	if err != nil {
		us.logger.Error("Invalid password has been provided",
			zap.String("error", err.Error()))

		return "", err
	}

	// generate the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name":  existingUserModel.Username,
		"user_email": existingUserModel.Email,
		"user_id":    existingUserModel.ID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(us.serverConfig.JwtSecretKey))

	if err != nil {
		us.logger.Fatal("Something went wrong while generating the token",
			zap.String("error", err.Error()))

		return "", err
	}

	us.logger.Info("Login user service was successful",
		zap.String("token", tokenString))

	return tokenString, nil
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newUserService
}
