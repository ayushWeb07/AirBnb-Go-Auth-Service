package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type RoleServiceInterface interface {
	CreateRole(rolePayload *dtos.CreateRolePayload) *utils.AppError
	GetAllRoles() ([]*models.RoleModel, *utils.AppError)
	GetRoleById(roleParams *dtos.GetRoleByIdParams) (*models.RoleModel, *utils.AppError)
	UpdateRoleById(roleParams *dtos.UpdateRoleByIdParams, rolePayload *dtos.UpdateRoleByIdPayload) *utils.AppError
	DeleteRoleById(roleParams *dtos.DeleteRoleByIdParams) *utils.AppError
}

type RoleService struct {
	RoleRepository repositories.RoleRepositoryInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (roleService *RoleService) CreateRole(rolePayload *dtos.CreateRolePayload) *utils.AppError {
	roleService.logger.Info("Create role service called...")

	// call the create role repository
	createRoleRepositoryErr := roleService.RoleRepository.CreateRole(rolePayload)
	return createRoleRepositoryErr
}

func (roleService *RoleService) GetAllRoles() ([]*models.RoleModel, *utils.AppError) {
	roleService.logger.Info("Get all roles service called...")

	// call the fetch all roles repository
	roleModels, getRolesRepositoryErr := roleService.RoleRepository.GetAllRoles()
	return roleModels, getRolesRepositoryErr
}

func (roleService *RoleService) GetRoleById(roleParams *dtos.GetRoleByIdParams) (*models.RoleModel, *utils.AppError) {
	roleService.logger.Info("Get by id role service called...")

	// call the fetch role by id repository
	roleModel, getRoleRepositoryErr := roleService.RoleRepository.GetRoleById(roleParams)
	return roleModel, getRoleRepositoryErr
}

func (roleService *RoleService) UpdateRoleById(roleParams *dtos.UpdateRoleByIdParams, rolePayload *dtos.UpdateRoleByIdPayload) *utils.AppError {
	roleService.logger.Info("Update by id role service called...")

	// call the update role by id repository
	updateRoleRepositoryErr := roleService.RoleRepository.UpdateRoleById(roleParams, rolePayload)
	return updateRoleRepositoryErr
}

func (roleService *RoleService) DeleteRoleById(roleParams *dtos.DeleteRoleByIdParams) *utils.AppError {
	roleService.logger.Info("Delete role service called...")

	// call the delete role by id repository
	deleteRoleRepositoryErr := roleService.RoleRepository.DeleteRoleById(roleParams)
	return deleteRoleRepositoryErr
}

func NewRoleService(roleRepository repositories.RoleRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RoleServiceInterface {
	newRoleService := &RoleService{
		RoleRepository: roleRepository,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newRoleService
}
