package services

import (
	"fmt"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
)

type UserServiceInterface interface {
	CreateUser()
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
}

func (us *UserService) CreateUser() {
	fmt.Println("Create user service called...")
	us.UserRepository.CreateUser()
}

func NewUserService(repo repositories.UserRepositoryInterface) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
	}

	return newUserService
}
