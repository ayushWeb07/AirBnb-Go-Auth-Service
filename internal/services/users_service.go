package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"go.uber.org/zap"
)

type UserServiceInterface interface {
	CreateUser()
	GetUserById()
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
}

func (us *UserService) GetUserById() {
	us.logger.Info("Get by id user service called...")
	us.UserRepository.GetUserById()
}

func (us *UserService) CreateUser() {
	us.logger.Info("Create user service called...")
	us.UserRepository.CreateUser()
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
		logger:         logger,
	}

	return newUserService
}
