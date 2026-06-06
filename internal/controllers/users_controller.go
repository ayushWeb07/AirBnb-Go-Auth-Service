package controllers

import "github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"

type UserControllerInterface interface {
	CreateUser() error
}

type UserController struct {
	UserService services.UserServiceInterface
}

func (u *UserController) CreateUser() error {
	return nil
}

func NewUserController(service services.UserServiceInterface) UserControllerInterface {
	newUserController := &UserController{
		UserService: service,
	}

	return newUserController
}
