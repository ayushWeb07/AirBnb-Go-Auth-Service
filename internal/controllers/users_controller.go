package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	renderPkg "github.com/unrolled/render"
	"go.uber.org/zap"
)

var render *renderPkg.Render

func init() {
	render = renderPkg.New()
}

type UserControllerInterface interface {
	CreateUser(resWriter http.ResponseWriter, req *http.Request)
	LoginUser(resWriter http.ResponseWriter, req *http.Request)
	GetAllUsers(resWriter http.ResponseWriter, req *http.Request)
	GetUserById(resWriter http.ResponseWriter, req *http.Request)
	DeleteUserById(resWriter http.ResponseWriter, req *http.Request)
}

type UserController struct {
	UserService  services.UserServiceInterface
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (uc *UserController) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(dtos.CreateUser)

	// call the create user service
	serviceErr := uc.UserService.CreateUser(&userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while creating the user",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusCreated, resWriter, map[string]any{
		"success":  true,
		"message":  "Successfully created the user",
		"email":    userPayload.Email,
		"username": userPayload.Username,
	})
}

func (uc *UserController) LoginUser(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(dtos.LoginUser)

	// call the login user service
	token, serviceErr := uc.UserService.LoginUser(&userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while logging in",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Login was successful",
		"token":   token,
	})
}

func (uc *UserController) GetAllUsers(resWriter http.ResponseWriter, req *http.Request) {
	// call the fetch all users service
	userModels, serviceErr := uc.UserService.GetAllUsers()

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting all users",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched all the users",
		"count":   len(userModels),
	})
}

func (uc *UserController) GetUserById(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(dtos.GetUserById)

	// call the fetch user by id service
	userModel, serviceErr := uc.UserService.GetUserById(&userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the user by id",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success":  true,
		"message":  "Successfully fetched the user by id",
		"email":    userModel.Email,
		"username": userModel.Username,
	})
}

func (uc *UserController) DeleteUserById(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(dtos.DeleteUserById)

	// call the delete user service
	serviceErr := uc.UserService.DeleteUserById(&userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while deleting the user",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully deleted the user",
	})
}

func NewUserController(service services.UserServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserControllerInterface {
	newUserController := &UserController{
		UserService:  service,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newUserController
}
