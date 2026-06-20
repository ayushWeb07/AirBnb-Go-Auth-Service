package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type RoleControllerInterface interface {
	CreateRole(resWriter http.ResponseWriter, req *http.Request)
	GetAllRoles(resWriter http.ResponseWriter, req *http.Request)
	GetRoleById(resWriter http.ResponseWriter, req *http.Request)
	UpdateRoleById(resWriter http.ResponseWriter, req *http.Request)
	DeleteRoleById(resWriter http.ResponseWriter, req *http.Request)
}

type RoleController struct {
	RoleService  services.RoleServiceInterface
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (roleController *RoleController) CreateRole(resWriter http.ResponseWriter, req *http.Request) {
	rolePayload := req.Context().Value("payload").(*dtos.CreateRolePayload)

	// call the create role service
	serviceErr := roleController.RoleService.CreateRole(rolePayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while creating the role",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusCreated, resWriter, map[string]any{
		"success": true,
		"message": "Successfully created the role",
	})
}

func (roleController *RoleController) GetAllRoles(resWriter http.ResponseWriter, req *http.Request) {
	// call the fetch roles service
	roleModels, serviceErr := roleController.RoleService.GetAllRoles()

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the roles",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the roles",
		"roles":   roleModels,
	})
}

func (roleController *RoleController) GetRoleById(resWriter http.ResponseWriter, req *http.Request) {
	roleParams := req.Context().Value("params").(*dtos.GetRoleByIdParams)

	// call the fetch user by id service
	roleModel, serviceErr := roleController.RoleService.GetRoleById(roleParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the role by id",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the role by id",
		"role":    roleModel,
	})
}

func (roleController *RoleController) UpdateRoleById(resWriter http.ResponseWriter, req *http.Request) {
	roleParams := req.Context().Value("params").(*dtos.UpdateRoleByIdParams)
	rolePayload := req.Context().Value("payload").(*dtos.UpdateRoleByIdPayload)

	// call the delete user service
	serviceErr := roleController.RoleService.UpdateRoleById(roleParams, rolePayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while updating the role",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully updated the role",
	})
}

func (roleController *RoleController) DeleteRoleById(resWriter http.ResponseWriter, req *http.Request) {
	roleParams := req.Context().Value("params").(*dtos.DeleteRoleByIdParams)

	// call the delete user service
	serviceErr := roleController.RoleService.DeleteRoleById(roleParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while deleting the role",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully deleted the role",
	})
}

func NewRoleController(service services.RoleServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RoleControllerInterface {
	newRoleController := &RoleController{
		RoleService:  service,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newRoleController
}
