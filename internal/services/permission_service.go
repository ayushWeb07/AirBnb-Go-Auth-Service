package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type PermissionServiceInterface interface {
	CreatePermission(permissionPayload *dtos.CreatePermissionPayload) *utils.AppError
	GetAllPermissions() ([]*models.PermissionModel, *utils.AppError)
	GetPermissionById(permissionParams *dtos.GetPermissionByIdParams) (*models.PermissionModel, *utils.AppError)
	UpdatePermissionById(permissionParams *dtos.UpdatePermissionByIdParams, permissionPayload *dtos.UpdatePermissionByIdPayload) *utils.AppError
	DeletePermissionById(permissionParams *dtos.DeletePermissionByIdParams) *utils.AppError
}

type PermissionService struct {
	PermissionRepository repositories.PermissionRepositoryInterface
	logger               *zap.Logger
	serverConfig         *config.ServerConfig
}

func (permissionService *PermissionService) CreatePermission(permissionPayload *dtos.CreatePermissionPayload) *utils.AppError {
	permissionService.logger.Info("Create permission service called...")

	// call the create permission repository
	repositoryErr := permissionService.PermissionRepository.CreatePermission(permissionPayload)
	return repositoryErr
}

func (permissionService *PermissionService) GetAllPermissions() ([]*models.PermissionModel, *utils.AppError) {
	permissionService.logger.Info("Get all permissions service called...")

	// call the fetch all permissions repository
	permissionModels, repositoryErr := permissionService.PermissionRepository.GetAllPermissions()
	return permissionModels, repositoryErr
}

func (permissionService *PermissionService) GetPermissionById(permissionParams *dtos.GetPermissionByIdParams) (*models.PermissionModel, *utils.AppError) {
	permissionService.logger.Info("Get by id permission service called...")

	// call the fetch permission by id repository
	permissionModel, repositoryErr := permissionService.PermissionRepository.GetPermissionById(permissionParams)
	return permissionModel, repositoryErr
}

func (permissionService *PermissionService) UpdatePermissionById(permissionParams *dtos.UpdatePermissionByIdParams, permissionPayload *dtos.UpdatePermissionByIdPayload) *utils.AppError {
	permissionService.logger.Info("Update by id permission service called...")

	// call the update permission by id repository
	repositoryErr := permissionService.PermissionRepository.UpdatePermissionById(permissionParams, permissionPayload)
	return repositoryErr
}

func (permissionService *PermissionService) DeletePermissionById(permissionParams *dtos.DeletePermissionByIdParams) *utils.AppError {
	permissionService.logger.Info("Delete permission service called...")

	// call the delete permission by id repository
	repositoryErr := permissionService.PermissionRepository.DeletePermissionById(permissionParams)
	return repositoryErr
}

func NewPermissionService(permissionRepository repositories.PermissionRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) PermissionServiceInterface {
	newPermissionService := &PermissionService{
		PermissionRepository: permissionRepository,
		logger:               logger,
		serverConfig:         serverConfig,
	}

	return newPermissionService
}
