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
	logger              *zap.Logger
	serverConfig        *config.ServerConfig
}

func (userRolesService *UserRolesService) GetRolesOfUser(userRolesPayload *dtos.GetRolesOfUserPayload) ([]*models.RoleModel, *utils.AppError) {
	userRolesService.logger.Info("Get user roles service called...")

	// call the fetch all user roles repository
	roleModels, repositoryErr := userRolesService.UserRolesRepository.GetRolesOfUser(userRolesPayload)
	return roleModels, repositoryErr
}

func (userRolesService *UserRolesService) AssignRoleToUser(userRolesPayload *dtos.AssignRoleToUserPayload) *utils.AppError {
	userRolesService.logger.Info("Assign user roles service called...")

	// call the assign user roles repository
	repositoryErr := userRolesService.UserRolesRepository.AssignRoleToUser(userRolesPayload)
	return repositoryErr
}

func (userRolesService *UserRolesService) RemoveRoleFromUser(userRolesPayload *dtos.RemoveRoleFromUserPayload) *utils.AppError {
	userRolesService.logger.Info("Remove user roles service called...")

	// call the remove user roles repository
	repositoryErr := userRolesService.UserRolesRepository.RemoveRoleFromUser(userRolesPayload)
	return repositoryErr
}

func (userRolesService *UserRolesService) CheckUserHasRole(userRolesPayload *dtos.CheckUserHasRolePayload) (bool, *utils.AppError) {
	userRolesService.logger.Info("Check user role service called...")

	// call the check user role repository
	hasRole, repositoryErr := userRolesService.UserRolesRepository.CheckUserHasRole(userRolesPayload)
	return hasRole, repositoryErr
}

func NewUserRolesService(userRolesRepository repositories.UserRolesRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserRolesServiceInterface {
	newUserRolesService := &UserRolesService{
		UserRolesRepository: userRolesRepository,
		logger:              logger,
		serverConfig:        serverConfig,
	}

	return newUserRolesService
}
