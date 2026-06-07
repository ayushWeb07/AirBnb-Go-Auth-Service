package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser()
	GetUserById()
	GetAllUsers()
	DeleteUserById()
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
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
		Username: "khabib",
		Email:    "khabib@gmail.com",
		Password: "khabib",
	}

	// hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), 14)

	if err != nil {
		us.logger.Fatal("Something went wrong while hashing the password",
			zap.String("error", err.Error()))

		return
	}

	userModel.Password = string(bytes)

	us.UserRepository.CreateUser(userModel)
}

func (us *UserService) DeleteUserById() {
	us.logger.Info("Delete user service called...")
	us.UserRepository.DeleteUserById()
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
		logger:         logger,
	}

	return newUserService
}
