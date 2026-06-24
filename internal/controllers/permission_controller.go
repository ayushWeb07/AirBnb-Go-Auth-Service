package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type PermissionControllerInterface interface {
	CreatePermission(resWriter http.ResponseWriter, req *http.Request)
	GetAllPermissions(resWriter http.ResponseWriter, req *http.Request)
	GetPermissionById(resWriter http.ResponseWriter, req *http.Request)
	UpdatePermissionById(resWriter http.ResponseWriter, req *http.Request)
	DeletePermissionById(resWriter http.ResponseWriter, req *http.Request)
}

type PermissionController struct {
	PermissionService services.PermissionServiceInterface
	logger            *zap.Logger
	serverConfig      *config.ServerConfig
}

func (permissionController *PermissionController) CreatePermission(resWriter http.ResponseWriter, req *http.Request) {
	permissionPayload := req.Context().Value("payload").(*dtos.CreatePermissionPayload)

	// call the create permission service
	serviceErr := permissionController.PermissionService.CreatePermission(permissionPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while creating the permission",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusCreated, resWriter, map[string]any{
		"success": true,
		"message": "Successfully created the permission",
	})
}

func (permissionController *PermissionController) GetAllPermissions(resWriter http.ResponseWriter, req *http.Request) {
	// call the fetch permissions service
	permissionModels, serviceErr := permissionController.PermissionService.GetAllPermissions()

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the permissions",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success":     true,
		"message":     "Successfully fetched the permissions",
		"permissions": permissionModels,
	})
}

func (permissionController *PermissionController) GetPermissionById(resWriter http.ResponseWriter, req *http.Request) {
	permissionParams := req.Context().Value("params").(*dtos.GetPermissionByIdParams)

	// call the fetch permission by id service
	permissionModel, serviceErr := permissionController.PermissionService.GetPermissionById(permissionParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the permission by id",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success":    true,
		"message":    "Successfully fetched the permission by id",
		"permission": permissionModel,
	})
}

func (permissionController *PermissionController) UpdatePermissionById(resWriter http.ResponseWriter, req *http.Request) {
	permissionParams := req.Context().Value("params").(*dtos.UpdatePermissionByIdParams)
	permissionPayload := req.Context().Value("payload").(*dtos.UpdatePermissionByIdPayload)

	// call the delete permission service
	serviceErr := permissionController.PermissionService.UpdatePermissionById(permissionParams, permissionPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while updating the permission",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully updated the permission",
	})
}

func (permissionController *PermissionController) DeletePermissionById(resWriter http.ResponseWriter, req *http.Request) {
	permissionParams := req.Context().Value("params").(*dtos.DeletePermissionByIdParams)

	// call the delete permission service
	serviceErr := permissionController.PermissionService.DeletePermissionById(permissionParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while deleting the permission",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully deleted the permission",
	})
}

func NewPermissionController(service services.PermissionServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) PermissionControllerInterface {
	newPermissionController := &PermissionController{
		PermissionService: service,
		logger:            logger,
		serverConfig:      serverConfig,
	}

	return newPermissionController
}
