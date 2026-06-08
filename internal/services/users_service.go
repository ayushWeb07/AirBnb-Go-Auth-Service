package services

import (
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser()
	GetUserById()
	GetAllUsers()
	DeleteUserById()
	LoginUser()
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (us *UserService) GetAllUsers() {
	us.logger.Info("Get all users service called...")
	us.UserRepository.GetAllUsers()
}

func (us *UserService) GetUserById() {
	us.logger.Info("Get by id user service called...")
	us.UserRepository.GetUserById()
}

func (us *UserService) CreateUser() {
	us.logger.Info("Create user service called...")

	// create a dummy user model instance
	userModel := &models.UserModel{
		Username: "conor",
		Email:    "conor@gmail.com",
		Password: "conor@2007",
	}

	// hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)

	if err != nil {
		us.logger.Fatal("Something went wrong while hashing the password",
			zap.String("error", err.Error()))

		return
	}

	userModel.Password = string(bytes)

	// call the repository endpoint
	err = us.UserRepository.CreateUser(userModel)

	if err != nil {
		return
	}

	// generate the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name":  userModel.Username,
		"user_email": userModel.Email,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(us.serverConfig.JwtSecretKey))

	if err != nil {
		us.logger.Fatal("Something went wrong while generating the token",
			zap.String("error", err.Error()))

		return
	}

	us.logger.Info("Create user service was successful",
		zap.String("token", tokenString))
}

func (us *UserService) DeleteUserById() {
	us.logger.Info("Delete user service called...")
	us.UserRepository.DeleteUserById()
}

func (us *UserService) LoginUser() {
	us.logger.Info("Login user service called...")

	// create a dummy user model instance
	userModel := &models.UserModel{
		Username: "conor",
		Email:    "conor@gmail.com",
		Password: "conor@2007",
	}

	// call the repository endpoint
	err := us.UserRepository.LoginUser(userModel)

	if err != nil {
		return
	}

	// generate the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name":  userModel.Username,
		"user_email": userModel.Email,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(us.serverConfig.JwtSecretKey))

	if err != nil {
		us.logger.Fatal("Something went wrong while generating the token",
			zap.String("error", err.Error()))

		return
	}

	us.logger.Info("Login user service was successful",
		zap.String("token", tokenString))
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newUserService
}
