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
	// call the fetch all users service
	userModels, err := uc.UserService.GetAllUsers()

	if err != nil {
		render.JSON(resWriter, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Something went wrong while getting all users",
			"error":   err.Error(),
		})
	} else {
		render.JSON(resWriter, http.StatusOK, map[string]any{
			"success": true,
			"message": "Successfully fetched all the users",
			"count":   len(userModels),
		})
	}
}

func (uc *UserController) GetUserById(resWriter http.ResponseWriter, req *http.Request) {
	// call the fetch user by id service
	userModel, err := uc.UserService.GetUserById()

	if err != nil {
		render.JSON(resWriter, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Something went wrong while getting the user by id",
			"error":   err.Error(),
		})
	} else {
		render.JSON(resWriter, http.StatusOK, map[string]any{
			"success":  true,
			"message":  "Successfully fetched the user by id",
			"email":    userModel.Email,
			"username": userModel.Username,
		})
	}
}

func (uc *UserController) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	// call the create user service
	userModel, err := uc.UserService.CreateUser()

	if err != nil {
		render.JSON(resWriter, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Something went wrong while creating the user",
			"error":   err.Error(),
		})
	} else {
		render.JSON(resWriter, http.StatusCreated, map[string]any{
			"success":  true,
			"message":  "Successfully created the user",
			"email":    userModel.Email,
			"username": userModel.Username,
		})
	}
}

func (uc *UserController) DeleteUserById(resWriter http.ResponseWriter, req *http.Request) {
	// call the delete user service
	err := uc.UserService.DeleteUserById()

	if err != nil {
		render.JSON(resWriter, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Something went wrong while deleting the user",
			"error":   err.Error(),
		})
	} else {
		render.JSON(resWriter, http.StatusOK, map[string]any{
			"success": true,
			"message": "Successfully deleted the user",
		})
	}
}

func (uc *UserController) LoginUser(resWriter http.ResponseWriter, req *http.Request) {
	// call the login user service
	token, err := uc.UserService.LoginUser()

	if err != nil {
		render.JSON(resWriter, http.StatusInternalServerError, map[string]any{
			"success": false,
			"message": "Something went wrong while logging in",
			"error":   err.Error(),
		})
	} else {
		render.JSON(resWriter, http.StatusOK, map[string]any{
			"success": true,
			"message": "Login was successful",
			"token":   token,
		})
	}
}

func NewUserController(service services.UserServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserControllerInterface {
	newUserController := &UserController{
		UserService:  service,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newUserController
}
