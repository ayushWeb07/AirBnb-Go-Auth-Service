package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type RolePermissionsControllerInterface interface {
	GetPermissionsOfUser(resWriter http.ResponseWriter, req *http.Request)
	CheckUserHasPermission(resWriter http.ResponseWriter, req *http.Request)
	GetPermissionsOfRole(resWriter http.ResponseWriter, req *http.Request)
	AssignPermissionToRole(resWriter http.ResponseWriter, req *http.Request)
	RemovePermissionFromRole(resWriter http.ResponseWriter, req *http.Request)
	CheckRoleHasPermission(resWriter http.ResponseWriter, req *http.Request)
}

type RolePermissionsController struct {
	RolePermissionsService services.RolePermissionsServiceInterface
	logger                 *zap.Logger
	serverConfig           *config.ServerConfig
}

func (rolePermissionsController *RolePermissionsController) GetPermissionsOfUser(resWriter http.ResponseWriter, req *http.Request) {
	userPermissionsPayload := req.Context().Value("payload").(*dtos.GetPermissionsOfUserPayload)

	// call the get user permissions service
	permissionModels, serviceErr := rolePermissionsController.RolePermissionsService.GetPermissionsOfUser(userPermissionsPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the permissions of the user",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success":     true,
		"message":     "Successfully fetched the permissions of the user",
		"permissions": permissionModels,
	})
}

func (rolePermissionsController *RolePermissionsController) CheckUserHasPermission(resWriter http.ResponseWriter, req *http.Request) {
	userPermissionsPayload := req.Context().Value("payload").(*dtos.CheckUserHasPermissionPayload)

	// call the check user permissions service
	hasPermission, serviceErr := rolePermissionsController.RolePermissionsService.CheckUserHasPermission(userPermissionsPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while checking if user has the permission",
			"error":   serviceErr.Error(),
		})

		return
	}

	if !hasPermission {
		utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
			"success": true,
			"message": "No, permission is not assigned with the mentioned user",
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Yes, role is assigned with the mentioned permission",
	})
}

func (rolePermissionsController *RolePermissionsController) GetPermissionsOfRole(resWriter http.ResponseWriter, req *http.Request) {
	rolePermissionsPayload := req.Context().Value("payload").(*dtos.GetPermissionsOfRolePayload)

	// call the get role permissions service
	permissionModels, serviceErr := rolePermissionsController.RolePermissionsService.GetPermissionsOfRole(rolePermissionsPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the permissions of the role",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success":     true,
		"message":     "Successfully fetched the permissions of the role",
		"permissions": permissionModels,
	})
}

func (rolePermissionsController *RolePermissionsController) AssignPermissionToRole(resWriter http.ResponseWriter, req *http.Request) {
	rolePermissionsPayload := req.Context().Value("payload").(*dtos.AssignPermissionToRolePayload)

	// call the assign role permissions service
	serviceErr := rolePermissionsController.RolePermissionsService.AssignPermissionToRole(rolePermissionsPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while assigning the permission to the role",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully assigned the permission to the role",
	})
}

func (rolePermissionsController *RolePermissionsController) RemovePermissionFromRole(resWriter http.ResponseWriter, req *http.Request) {
	rolePermissionsPayload := req.Context().Value("payload").(*dtos.RemovePermissionFromRolePayload)

	// call the remove role permissions service
	serviceErr := rolePermissionsController.RolePermissionsService.RemovePermissionFromRole(rolePermissionsPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while removing the permission from the role",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully removed the permission from the role",
	})
}

func (rolePermissionsController *RolePermissionsController) CheckRoleHasPermission(resWriter http.ResponseWriter, req *http.Request) {
	rolePermissionsPayload := req.Context().Value("payload").(*dtos.CheckRoleHasPermissionPayload)

	// call the check role permissions service
	hasPermission, serviceErr := rolePermissionsController.RolePermissionsService.CheckRoleHasPermission(rolePermissionsPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while checking if role has the permission",
			"error":   serviceErr.Error(),
		})

		return
	}

	if !hasPermission {
		utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
			"success": true,
			"message": "No, role is not assigned with the mentioned permission",
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Yes, role is assigned with the mentioned permission",
	})
}

func NewRolePermissionsController(service services.RolePermissionsServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RolePermissionsControllerInterface {
	newRolePermissionsController := &RolePermissionsController{
		RolePermissionsService: service,
		logger:                 logger,
		serverConfig:           serverConfig,
	}

	return newRolePermissionsController
}
