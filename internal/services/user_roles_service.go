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
	_, userRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if userRepositoryErr != nil {
		return nil, userRepositoryErr
	}

	// call the fetch all user roles repository
	roleModels, repositoryErr := userRolesService.UserRolesRepository.GetRolesOfUser(userRolesPayload)
	return roleModels, repositoryErr
}

func (userRolesService *UserRolesService) AssignRoleToUser(userRolesPayload *dtos.AssignRoleToUserPayload) *utils.AppError {
	userRolesService.logger.Info("Assign user roles service called...")

	// check if user exists
	_, userRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if userRepositoryErr != nil {
		return userRepositoryErr
	}

	// check if role exists
	_, roleRepositoryErr := userRolesService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: userRolesPayload.RoleID})

	if roleRepositoryErr != nil {
		return roleRepositoryErr
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
	repositoryErr := userRolesService.UserRolesRepository.AssignRoleToUser(userRolesPayload)
	return repositoryErr
}

func (userRolesService *UserRolesService) RemoveRoleFromUser(userRolesPayload *dtos.RemoveRoleFromUserPayload) *utils.AppError {
	userRolesService.logger.Info("Remove user roles service called...")

	// check if user exists
	_, userRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if userRepositoryErr != nil {
		return userRepositoryErr
	}

	// check if role exists
	_, roleRepositoryErr := userRolesService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: userRolesPayload.RoleID})

	if roleRepositoryErr != nil {
		return roleRepositoryErr
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
	repositoryErr := userRolesService.UserRolesRepository.RemoveRoleFromUser(userRolesPayload)
	return repositoryErr
}

func (userRolesService *UserRolesService) CheckUserHasRole(userRolesPayload *dtos.CheckUserHasRolePayload) (bool, *utils.AppError) {
	userRolesService.logger.Info("Check user role service called...")

	// check if user exists
	_, userRepositoryErr := userRolesService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userRolesPayload.UserID})

	if userRepositoryErr != nil {
		return false, userRepositoryErr
	}

	// check if role exists
	_, roleRepositoryErr := userRolesService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: userRolesPayload.RoleID})

	if roleRepositoryErr != nil {
		return false, roleRepositoryErr
	}

	// call the check user role repository
	hasRole := userRolesService.UserRolesRepository.CheckUserHasRole(userRolesPayload)
	return hasRole, nil
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
