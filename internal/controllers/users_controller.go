package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"go.uber.org/zap"
)

type UserControllerInterface interface {
	CreateUser(resWriter http.ResponseWriter, req *http.Request)
	GetUserById(resWriter http.ResponseWriter, req *http.Request)
	GetAllUsers(resWriter http.ResponseWriter, req *http.Request)
}

type UserController struct {
	UserService services.UserServiceInterface
	logger      *zap.Logger
}

func (uc *UserController) GetAllUsers(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Get all users endpoint working fine!"))
	uc.UserService.GetAllUsers()
}

func (uc *UserController) GetUserById(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Get by id user endpoint working fine!"))
	uc.UserService.GetUserById()
}

func (uc *UserController) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Create user endpoint working fine!"))
	uc.UserService.CreateUser()
}

func NewUserController(service services.UserServiceInterface, logger *zap.Logger) UserControllerInterface {
	newUserController := &UserController{
		UserService: service,
		logger:      logger,
	}

	return newUserController
}
