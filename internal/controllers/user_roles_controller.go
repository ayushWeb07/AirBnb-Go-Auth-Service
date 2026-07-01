package controllers

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type UserRolesControllerInterface interface {
	GetRolesOfUser(resWriter http.ResponseWriter, req *http.Request)
	AssignRoleToUser(resWriter http.ResponseWriter, req *http.Request)
	RemoveRoleFromUser(resWriter http.ResponseWriter, req *http.Request)
	CheckUserHasRole(resWriter http.ResponseWriter, req *http.Request)
	CheckUserHasAllRoles(resWriter http.ResponseWriter, req *http.Request)
	CheckUserHasAnyRoles(resWriter http.ResponseWriter, req *http.Request)
	GetUserRolesService() services.UserRolesServiceInterface
}

type UserRolesController struct {
	UserRolesService services.UserRolesServiceInterface
	logger           *zap.Logger
	serverConfig     *config.ServerConfig
}

func (userRolesController *UserRolesController) GetRolesOfUser(resWriter http.ResponseWriter, req *http.Request) {
	userRolesPayload := req.Context().Value("params").(*dtos.GetRolesOfUserPayload)

	// call the get user roles service
	roleModels, getRolesServiceErr := userRolesController.UserRolesService.GetRolesOfUser(userRolesPayload)

	if getRolesServiceErr != nil {
		utils.WriteJsonResponse(getRolesServiceErr.StatusCode, resWriter, map[string]any{
			"success": getRolesServiceErr.Success,
			"message": "Something went wrong while getting the roles of the user",
			"error":   getRolesServiceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the roles of the user",
		"roles":   roleModels,
	})
}

func (userRolesController *UserRolesController) AssignRoleToUser(resWriter http.ResponseWriter, req *http.Request) {
	userRolesPayload := req.Context().Value("payload").(*dtos.AssignRoleToUserPayload)

	// call the assign user roles service
	assignRoleServiceErr := userRolesController.UserRolesService.AssignRoleToUser(userRolesPayload)

	if assignRoleServiceErr != nil {
		utils.WriteJsonResponse(assignRoleServiceErr.StatusCode, resWriter, map[string]any{
			"success": assignRoleServiceErr.Success,
			"message": "Something went wrong while assigning the role to the user",
			"error":   assignRoleServiceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully assigned the role to the user",
	})
}

func (userRolesController *UserRolesController) RemoveRoleFromUser(resWriter http.ResponseWriter, req *http.Request) {
	userRolesPayload := req.Context().Value("payload").(*dtos.RemoveRoleFromUserPayload)

	// call the remove user roles service
	removeRoleServiceErr := userRolesController.UserRolesService.RemoveRoleFromUser(userRolesPayload)

	if removeRoleServiceErr != nil {
		utils.WriteJsonResponse(removeRoleServiceErr.StatusCode, resWriter, map[string]any{
			"success": removeRoleServiceErr.Success,
			"message": "Something went wrong while removing the role from the user",
			"error":   removeRoleServiceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully removed the role from the user",
	})
}

func (userRolesController *UserRolesController) CheckUserHasRole(resWriter http.ResponseWriter, req *http.Request) {
	userRolesPayload := req.Context().Value("payload").(*dtos.CheckUserHasRolePayload)

	// call the check user roles service
	hasRole, checkUserHasRoleServiceErr := userRolesController.UserRolesService.CheckUserHasRole(userRolesPayload)

	if checkUserHasRoleServiceErr != nil {
		utils.WriteJsonResponse(checkUserHasRoleServiceErr.StatusCode, resWriter, map[string]any{
			"success": checkUserHasRoleServiceErr.Success,
			"message": "Something went wrong while checking if user has the role",
			"error":   checkUserHasRoleServiceErr.Error(),
		})

		return
	}

	if !hasRole {
		utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
			"success": true,
			"message": "No, user is not assigned with the mentioned role",
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Yes, user is assigned with the mentioned role",
	})
}

func (userRolesController *UserRolesController) CheckUserHasAllRoles(resWriter http.ResponseWriter, req *http.Request) {
	userRolesPayload := req.Context().Value("payload").(*dtos.CheckUserHasAllRolesPayload)

	// call the check user roles service
	hasAllRoles, checkUserHasAllRolesServiceErr := userRolesController.UserRolesService.CheckUserHasAllRoles(userRolesPayload)

	if checkUserHasAllRolesServiceErr != nil {
		utils.WriteJsonResponse(checkUserHasAllRolesServiceErr.StatusCode, resWriter, map[string]any{
			"success": checkUserHasAllRolesServiceErr.Success,
			"message": "Something went wrong while checking if user has all the roles",
			"error":   checkUserHasAllRolesServiceErr.Error(),
		})

		return
	}

	if !hasAllRoles {
		utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
			"success": true,
			"message": "No, user is not assigned with all the mentioned roles",
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Yes, user is assigned with all the mentioned roles",
	})
}

func (userRolesController *UserRolesController) CheckUserHasAnyRoles(resWriter http.ResponseWriter, req *http.Request) {
	userRolesPayload := req.Context().Value("payload").(*dtos.CheckUserHasAnyRolesPayload)

	// call the check user roles service
	hasAnyRoles, checkUserHasAnyRolesServiceErr := userRolesController.UserRolesService.CheckUserHasAnyRoles(userRolesPayload)

	if checkUserHasAnyRolesServiceErr != nil {
		utils.WriteJsonResponse(checkUserHasAnyRolesServiceErr.StatusCode, resWriter, map[string]any{
			"success": checkUserHasAnyRolesServiceErr.Success,
			"message": "Something went wrong while checking if user has any required roles",
			"error":   checkUserHasAnyRolesServiceErr.Error(),
		})

		return
	}

	if !hasAnyRoles {
		utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
			"success": true,
			"message": "No, user is not assigned with any of the mentioned roles",
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Yes, user is assigned with any of the mentioned roles",
	})
}

func (userRolesController *UserRolesController) GetUserRolesService() services.UserRolesServiceInterface {
	return userRolesController.UserRolesService
}

func NewUserRolesController(service services.UserRolesServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserRolesControllerInterface {
	newUserRolesController := &UserRolesController{
		UserRolesService: service,
		logger:           logger,
		serverConfig:     serverConfig,
	}

	return newUserRolesController
}
