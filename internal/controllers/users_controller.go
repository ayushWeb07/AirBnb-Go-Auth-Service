package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	renderPkg "github.com/unrolled/render"
	"go.uber.org/zap"
)

var render *renderPkg.Render

func init() {
	render = renderPkg.New()
}

type UserControllerInterface interface {
	CreateUser(resWriter http.ResponseWriter, req *http.Request)
	GetUserById(resWriter http.ResponseWriter, req *http.Request)
	GetAllUsers(resWriter http.ResponseWriter, req *http.Request)
	DeleteUserById(resWriter http.ResponseWriter, req *http.Request)
	LoginUser(resWriter http.ResponseWriter, req *http.Request)
}

type UserController struct {
	UserService  services.UserServiceInterface
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (uc *UserController) GetAllUsers(resWriter http.ResponseWriter, req *http.Request) {
	render.JSON(resWriter, http.StatusOK, map[string]any{
		"message": "Get all users endpoint working fine!",
		"status":  "Fine!",
	})

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

func (uc *UserController) DeleteUserById(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Delete user endpoint working fine!"))
	uc.UserService.DeleteUserById()
}

func (uc *UserController) LoginUser(resWriter http.ResponseWriter, req *http.Request) {
	resWriter.Write([]byte("Login user endpoint working fine!"))
	uc.UserService.LoginUser()
}

func NewUserController(service services.UserServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserControllerInterface {
	newUserController := &UserController{
		UserService:  service,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newUserController
}
