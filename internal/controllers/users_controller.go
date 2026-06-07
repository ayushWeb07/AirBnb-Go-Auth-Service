package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
)

type UserControllerInterface interface {
	CreateUser(resWriter http.ResponseWriter, req *http.Request)
}

type UserController struct {
	UserService services.UserServiceInterface
}

func (uc *UserController) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Create user endpoint working fine!"))
	uc.UserService.CreateUser()
}

func NewUserController(service services.UserServiceInterface) UserControllerInterface {
	newUserController := &UserController{
		UserService: service,
	}

	return newUserController
}
