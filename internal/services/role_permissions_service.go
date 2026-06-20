package services

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type RolePermissionsServiceInterface interface {
	GetPermissionsOfUser(userPermissionsPayload *dtos.GetPermissionsOfUserPayload) ([]*models.PermissionModel, *utils.AppError)
	CheckUserHasPermission(userPermissionsPayload *dtos.CheckUserHasPermissionPayload) (bool, *utils.AppError)

	GetPermissionsOfRole(rolePermissionsPayload *dtos.GetPermissionsOfRolePayload) ([]*models.PermissionModel, *utils.AppError)
	AssignPermissionToRole(rolePermissionsPayload *dtos.AssignPermissionToRolePayload) *utils.AppError
	RemovePermissionFromRole(rolePermissionsPayload *dtos.RemovePermissionFromRolePayload) *utils.AppError
	CheckRoleHasPermission(rolePermissionsPayload *dtos.CheckRoleHasPermissionPayload) (bool, *utils.AppError)
}

type RolePermissionsService struct {
	RolePermissionsRepository repositories.RolePermissionsRepositoryInterface
	RoleRepository            repositories.RoleRepositoryInterface
	PermissionRepository      repositories.PermissionRepositoryInterface
	UserRepository            repositories.UserRepositoryInterface
	logger                    *zap.Logger
	serverConfig              *config.ServerConfig
}

func (rolePermissionsService *RolePermissionsService) GetPermissionsOfUser(userPermissionsPayload *dtos.GetPermissionsOfUserPayload) ([]*models.PermissionModel, *utils.AppError) {
	rolePermissionsService.logger.Info("Get user permissions service called...")

	// check if user exists
	_, userRepositoryErr := rolePermissionsService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userPermissionsPayload.UserID})

	if userRepositoryErr != nil {
		return nil, userRepositoryErr
	}

	// call the fetch all user permissions repository
	permissionModels, repositoryErr := rolePermissionsService.RolePermissionsRepository.GetPermissionsOfUser(userPermissionsPayload)
	return permissionModels, repositoryErr
}

func (rolePermissionsService *RolePermissionsService) CheckUserHasPermission(userPermissionsPayload *dtos.CheckUserHasPermissionPayload) (bool, *utils.AppError) {
	rolePermissionsService.logger.Info("Check user role service called...")

	// check if user exists
	_, userRepositoryErr := rolePermissionsService.UserRepository.GetUserById(&dtos.GetUserByIdParams{ID: userPermissionsPayload.UserID})

	if userRepositoryErr != nil {
		return false, userRepositoryErr
	}

	// check if permission exists
	_, permissionRepositoryErr := rolePermissionsService.PermissionRepository.GetPermissionById(&dtos.GetPermissionByIdParams{ID: userPermissionsPayload.PermissionID})

	if permissionRepositoryErr != nil {
		return false, permissionRepositoryErr
	}

	// call the check user permission repository
	hasPermission := rolePermissionsService.RolePermissionsRepository.CheckUserHasPermission(userPermissionsPayload)
	return hasPermission, nil
}

func (rolePermissionsService *RolePermissionsService) GetPermissionsOfRole(rolePermissionsPayload *dtos.GetPermissionsOfRolePayload) ([]*models.PermissionModel, *utils.AppError) {
	rolePermissionsService.logger.Info("Get role permissions service called...")

	// check if role exists
	_, roleRepositoryErr := rolePermissionsService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: rolePermissionsPayload.RoleID})

	if roleRepositoryErr != nil {
		return nil, roleRepositoryErr
	}

	// call the fetch all role permissions repository
	permissionModels, repositoryErr := rolePermissionsService.RolePermissionsRepository.GetPermissionsOfRole(rolePermissionsPayload)
	return permissionModels, repositoryErr
}

func (rolePermissionsService *RolePermissionsService) AssignPermissionToRole(rolePermissionsPayload *dtos.AssignPermissionToRolePayload) *utils.AppError {
	rolePermissionsService.logger.Info("Assign role permissions service called...")

	// check if role exists
	_, roleRepositoryErr := rolePermissionsService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: rolePermissionsPayload.RoleID})

	if roleRepositoryErr != nil {
		return roleRepositoryErr
	}

	// check if permission exists
	_, permissionRepositoryErr := rolePermissionsService.PermissionRepository.GetPermissionById(&dtos.GetPermissionByIdParams{ID: rolePermissionsPayload.PermissionID})

	if permissionRepositoryErr != nil {
		return permissionRepositoryErr
	}

	// check if role already has the permission
	hasPermission := rolePermissionsService.RolePermissionsRepository.CheckRoleHasPermission(&dtos.CheckRoleHasPermissionPayload{
		PermissionID: rolePermissionsPayload.PermissionID,
		RoleID:       rolePermissionsPayload.RoleID,
	})

	if hasPermission {
		return utils.BadRequest("Role already has the permission")
	}

	// call the assign role permissions repository
	repositoryErr := rolePermissionsService.RolePermissionsRepository.AssignPermissionToRole(rolePermissionsPayload)
	return repositoryErr
}

func (rolePermissionsService *RolePermissionsService) RemovePermissionFromRole(rolePermissionsPayload *dtos.RemovePermissionFromRolePayload) *utils.AppError {
	rolePermissionsService.logger.Info("Remove role permissions service called...")

	// check if role exists
	_, roleRepositoryErr := rolePermissionsService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: rolePermissionsPayload.RoleID})

	if roleRepositoryErr != nil {
		return roleRepositoryErr
	}

	// check if permission exists
	_, permissionRepositoryErr := rolePermissionsService.PermissionRepository.GetPermissionById(&dtos.GetPermissionByIdParams{ID: rolePermissionsPayload.PermissionID})

	if permissionRepositoryErr != nil {
		return permissionRepositoryErr
	}

	// check if role already has the permission
	hasPermission := rolePermissionsService.RolePermissionsRepository.CheckRoleHasPermission(&dtos.CheckRoleHasPermissionPayload{
		RoleID:       rolePermissionsPayload.RoleID,
		PermissionID: rolePermissionsPayload.PermissionID,
	})

	if !hasPermission {
		return utils.BadRequest("Role does not have the permission")
	}

	// call the remove role permissions repository
	repositoryErr := rolePermissionsService.RolePermissionsRepository.RemovePermissionFromRole(rolePermissionsPayload)
	return repositoryErr
}

func (rolePermissionsService *RolePermissionsService) CheckRoleHasPermission(rolePermissionsPayload *dtos.CheckRoleHasPermissionPayload) (bool, *utils.AppError) {
	rolePermissionsService.logger.Info("Check role permissions service called...")

	// check if role exists
	_, roleRepositoryErr := rolePermissionsService.RoleRepository.GetRoleById(&dtos.GetRoleByIdParams{ID: rolePermissionsPayload.RoleID})

	if roleRepositoryErr != nil {
		return false, roleRepositoryErr
	}

	// check if permission exists
	_, permissionRepositoryErr := rolePermissionsService.PermissionRepository.GetPermissionById(&dtos.GetPermissionByIdParams{ID: rolePermissionsPayload.PermissionID})

	if permissionRepositoryErr != nil {
		return false, permissionRepositoryErr
	}

	// call the check role permission repository
	hasPermission := rolePermissionsService.RolePermissionsRepository.CheckRoleHasPermission(rolePermissionsPayload)
	return hasPermission, nil
}

func NewRolePermissionsService(rolePermissionsRepository repositories.RolePermissionsRepositoryInterface, roleRepository repositories.RoleRepositoryInterface, permissionRepository repositories.PermissionRepositoryInterface, userRepository repositories.UserRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RolePermissionsServiceInterface {
	newRolePermissionsService := &RolePermissionsService{
		RolePermissionsRepository: rolePermissionsRepository,
		RoleRepository:            roleRepository,
		PermissionRepository:      permissionRepository,
		UserRepository:            userRepository,
		logger:                    logger,
		serverConfig:              serverConfig,
	}

	return newRolePermissionsService
}
