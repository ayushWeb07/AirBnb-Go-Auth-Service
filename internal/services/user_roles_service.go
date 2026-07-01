package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type UserRolesServiceInterface interface {
	GetRolesOfUser(userRolesPayload *dtos.GetRolesOfUserPayload) ([]*models.RoleModel, *utils.AppError)
	AssignRoleToUser(userRolesPayload *dtos.AssignRoleToUserPayload) *utils.AppError
	RemoveRoleFromUser(userRolesPayload *dtos.RemoveRoleFromUserPayload) *utils.AppError
	CheckUserHasRole(userRolesPayload *dtos.CheckUserHasRolePayload) (bool, *utils.AppError)

	CheckUserHasAllRoles(userRolesPayload *dtos.CheckUserHasAllRolesPayload) (bool, *utils.AppError)
	CheckUserHasAnyRoles(userRolesPayload *dtos.CheckUserHasAnyRolesPayload) (bool, *utils.AppError)
	GetUserRolesRepository() repositories.UserRolesRepositoryInterface
}

type UserRolesService struct {
	UserRolesRepository repositories.UserRolesRepositoryInterface
	UserRepository      repositories.UserRepositoryInterface
	RoleRepository      repositories.RoleRepositoryInterface
	logger              *zap.Logger
	serverConfig        *config.ServerConfig
}

func (userRolesService *UserRolesService) GetRolesOfUser(userRolesPayload *dtos.GetRolesOfUserPayload) ([]*models.RoleModel, *utils.AppError) {
	userRolesService.logger.Info("Get user roles service called...")

	// check if user exists
	_, getUserRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if getUserRepositoryErr != nil {
		return nil, getUserRepositoryErr
	}

	// call the fetch all user roles repository
	roleModels, getRolesRepositoryErr := userRolesService.UserRolesRepository.GetRolesOfUser(userRolesPayload)
	return roleModels, getRolesRepositoryErr
}

func (userRolesService *UserRolesService) AssignRoleToUser(userRolesPayload *dtos.AssignRoleToUserPayload) *utils.AppError {
	userRolesService.logger.Info("Assign user roles service called...")

	// check if user exists
	_, getUserRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if getUserRepositoryErr != nil {
		return getUserRepositoryErr
	}

	// check if role exists
	_, getRoleRepositoryErr := userRolesService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: userRolesPayload.RoleID})

	if getRoleRepositoryErr != nil {
		return getRoleRepositoryErr
	}

	// check if user already has the role
	hasRole := userRolesService.UserRolesRepository.CheckUserHasRole(&dtos.CheckUserHasRolePayload{
		UserID: userRolesPayload.UserID,
		RoleID: userRolesPayload.RoleID,
	})

	if hasRole {
		return utils.BadRequest("User is already assigned with that role")
	}

	// call the assign user roles repository
	assignRoleRepositoryErr := userRolesService.UserRolesRepository.AssignRoleToUser(userRolesPayload)
	return assignRoleRepositoryErr
}

func (userRolesService *UserRolesService) RemoveRoleFromUser(userRolesPayload *dtos.RemoveRoleFromUserPayload) *utils.AppError {
	userRolesService.logger.Info("Remove user roles service called...")

	// check if user exists
	_, getUserRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if getUserRepositoryErr != nil {
		return getUserRepositoryErr
	}

	// check if role exists
	_, getRoleRepositoryErr := userRolesService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: userRolesPayload.RoleID})

	if getRoleRepositoryErr != nil {
		return getRoleRepositoryErr
	}

	// check if user has the role
	hasRole := userRolesService.UserRolesRepository.CheckUserHasRole(&dtos.CheckUserHasRolePayload{
		UserID: userRolesPayload.UserID,
		RoleID: userRolesPayload.RoleID,
	})

	if !hasRole {
		return utils.BadRequest("User is not assigned with that role")
	}

	// call the remove user roles repository
	removeRoleRepositoryErr := userRolesService.UserRolesRepository.RemoveRoleFromUser(userRolesPayload)
	return removeRoleRepositoryErr
}

func (userRolesService *UserRolesService) CheckUserHasRole(userRolesPayload *dtos.CheckUserHasRolePayload) (bool, *utils.AppError) {
	userRolesService.logger.Info("Check user role service called...")

	// check if user exists
	_, getUserRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if getUserRepositoryErr != nil {
		return false, getUserRepositoryErr
	}

	// check if role exists
	_, getRoleRepositoryErr := userRolesService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: userRolesPayload.RoleID})

	if getRoleRepositoryErr != nil {
		return false, getRoleRepositoryErr
	}

	// call the check user role repository
	hasRole := userRolesService.UserRolesRepository.CheckUserHasRole(userRolesPayload)
	return hasRole, nil
}

func (userRolesService *UserRolesService) CheckUserHasAllRoles(userRolesPayload *dtos.CheckUserHasAllRolesPayload) (bool, *utils.AppError) {
	userRolesService.logger.Info("Check user has all roles service called...")

	// check if user exists
	_, getUserRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if getUserRepositoryErr != nil {
		return false, getUserRepositoryErr
	}

	// call the check user role repository
	hasAllRoles := userRolesService.UserRolesRepository.CheckUserHasAllRoles(userRolesPayload)
	return hasAllRoles, nil
}

func (userRolesService *UserRolesService) CheckUserHasAnyRoles(userRolesPayload *dtos.CheckUserHasAnyRolesPayload) (bool, *utils.AppError) {
	userRolesService.logger.Info("Check user has any roles service called...")

	// check if user exists
	_, getUserRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if getUserRepositoryErr != nil {
		return false, getUserRepositoryErr
	}

	// call the check user role repository
	hasAnyRoles := userRolesService.UserRolesRepository.CheckUserHasAnyRoles(userRolesPayload)
	return hasAnyRoles, nil
}

func (userRolesService *UserRolesService) GetUserRolesRepository() repositories.UserRolesRepositoryInterface {
	userRolesService.logger.Info("Get user roles repository, service called...")

	return userRolesService.UserRolesRepository
}

func NewUserRolesService(userRolesRepository repositories.UserRolesRepositoryInterface, userRepository repositories.UserRepositoryInterface, roleRepository repositories.RoleRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserRolesServiceInterface {
	newUserRolesService := &UserRolesService{
		UserRolesRepository: userRolesRepository,
		UserRepository:      userRepository,
		RoleRepository:      roleRepository,
		logger:              logger,
		serverConfig:        serverConfig,
	}

	return newUserRolesService
}
