package repositories

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type RolePermissionsRepositoryInterface interface {
	GetPermissionsOfUser(userPermissionsPayload *dtos.GetPermissionsOfUserPayload) ([]*models.PermissionModel, *utils.AppError)
	CheckUserHasPermission(userPermissionsPayload *dtos.CheckUserHasPermissionPayload) bool

	GetPermissionsOfRole(rolePermissionsPayload *dtos.GetPermissionsOfRolePayload) ([]*models.PermissionModel, *utils.AppError)
	AssignPermissionToRole(rolePermissionsPayload *dtos.AssignPermissionToRolePayload) *utils.AppError
	RemovePermissionFromRole(rolePermissionsPayload *dtos.RemovePermissionFromRolePayload) *utils.AppError
	CheckRoleHasPermission(rolePermissionsPayload *dtos.CheckRoleHasPermissionPayload) bool
}

type RolePermissionsRepository struct {
	db           *sql.DB
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (rolePermissionsRepository *RolePermissionsRepository) GetPermissionsOfUser(userPermissionsPayload *dtos.GetPermissionsOfUserPayload) ([]*models.PermissionModel, *utils.AppError) {
	// create the dummy instance
	var permissionModels []*models.PermissionModel

	// load the rows
	query :=
		`SELECT p.* FROM user_roles ur 
		 INNER JOIN role_permissions rp ON rp.role_id = ur.role_id 
    	 INNER JOIN permissions p ON p.id = rp.permission_id 
         WHERE ur.user_id = ?`

	rows, queryErr := rolePermissionsRepository.db.Query(query, userPermissionsPayload.UserID)

	if queryErr != nil {
		rolePermissionsRepository.logger.Error("Something went wrong while fetching the permissions",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching the permissions: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		permissionModel := &models.PermissionModel{}

		rowScanErr := rows.Scan(&permissionModel.ID, &permissionModel.Name, &permissionModel.Description, &permissionModel.Resource, &permissionModel.Action, &permissionModel.CreatedAt, &permissionModel.UpdatedAt)

		if rowScanErr != nil {
			rolePermissionsRepository.logger.Error("Failed to fetch the permissions from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching the permissions: " + rowScanErr.Error())
		}

		permissionModels = append(permissionModels, permissionModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		rolePermissionsRepository.logger.Error("Failed to fetch the permissions from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching the permissions: " + rowsErr.Error())
	}

	rolePermissionsRepository.logger.Info("Successfully fetched the permissions from the database",
		zap.Int("count", len(permissionModels)))

	return permissionModels, nil
}

func (rolePermissionsRepository *RolePermissionsRepository) CheckUserHasPermission(userPermissionsPayload *dtos.CheckUserHasPermissionPayload) bool {
	// create the dummy instance
	permissionModel := &models.PermissionModel{}

	// fetch from the db
	query :=
		`SELECT p.id FROM user_roles ur 
		 INNER JOIN role_permissions rp ON rp.role_id = ur.role_id 
    	 INNER JOIN permissions p ON p.id = rp.permission_id 
         WHERE ur.user_id = ? AND rp.permission_id = ?`

	queryErr := rolePermissionsRepository.db.QueryRow(query, userPermissionsPayload.UserID, userPermissionsPayload.PermissionID).Scan(&permissionModel.ID)

	if queryErr != nil {
		rolePermissionsRepository.logger.Error("User does not have the permission",
			zap.Int("user_id", userPermissionsPayload.UserID),
			zap.Int("permission_id", userPermissionsPayload.PermissionID),
		)

		return false
	}

	rolePermissionsRepository.logger.Info("User has the particular permission",
		zap.Int("user_id", userPermissionsPayload.UserID),
		zap.Int("permission_id", userPermissionsPayload.PermissionID),
	)

	return true
}

func (rolePermissionsRepository *RolePermissionsRepository) GetPermissionsOfRole(rolePermissionsPayload *dtos.GetPermissionsOfRolePayload) ([]*models.PermissionModel, *utils.AppError) {
	// create the dummy instance
	var permissionModels []*models.PermissionModel

	// load the rows
	query := "SELECT p.* FROM role_permissions rp INNER JOIN permissions p ON p.id = rp.permission_id WHERE rp.role_id = ?"
	rows, queryErr := rolePermissionsRepository.db.Query(query, rolePermissionsPayload.RoleID)

	if queryErr != nil {
		rolePermissionsRepository.logger.Error("Something went wrong while fetching the permissions",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching the permissions: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		permissionModel := &models.PermissionModel{}

		rowScanErr := rows.Scan(&permissionModel.ID, &permissionModel.Name, &permissionModel.Description, &permissionModel.Resource, &permissionModel.Action, &permissionModel.CreatedAt, &permissionModel.UpdatedAt)

		if rowScanErr != nil {
			rolePermissionsRepository.logger.Error("Failed to fetch the permissions from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching the permissions: " + rowScanErr.Error())
		}

		permissionModels = append(permissionModels, permissionModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		rolePermissionsRepository.logger.Error("Failed to fetch the permissions from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching the permissions: " + rowsErr.Error())
	}

	rolePermissionsRepository.logger.Info("Successfully fetched the permissions from the database",
		zap.Int("count", len(permissionModels)))

	return permissionModels, nil
}

func (rolePermissionsRepository *RolePermissionsRepository) AssignPermissionToRole(rolePermissionsPayload *dtos.AssignPermissionToRolePayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)"
	result, queryExecErr := rolePermissionsRepository.db.Exec(query, rolePermissionsPayload.RoleID, rolePermissionsPayload.PermissionID)

	if queryExecErr != nil {
		rolePermissionsRepository.logger.Error("Failed to assign permission to the role",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to assign permission to the role: " + queryExecErr.Error())
	}

	_, insertErr := result.LastInsertId()

	if insertErr != nil {
		rolePermissionsRepository.logger.Error("Failed to assign permission to the role",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to assign permission to the role: " + insertErr.Error())
	}

	rolePermissionsRepository.logger.Info("Successfully assigned permission to the role",
		zap.Int("role_id", rolePermissionsPayload.RoleID),
		zap.Int("permission_id", rolePermissionsPayload.PermissionID),
	)

	return nil
}

func (rolePermissionsRepository *RolePermissionsRepository) RemovePermissionFromRole(rolePermissionsPayload *dtos.RemovePermissionFromRolePayload) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?"
	result, queryExecErr := rolePermissionsRepository.db.Exec(query, rolePermissionsPayload.RoleID, rolePermissionsPayload.PermissionID)

	if queryExecErr != nil {
		rolePermissionsRepository.logger.Error("Failed to remove permission from the role",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to remove permission from the role: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		rolePermissionsRepository.logger.Error("Failed to remove permission from the role",
			zap.Int("role_id", rolePermissionsPayload.RoleID),
			zap.Int("permission_id", rolePermissionsPayload.PermissionID),
		)

		return utils.InternalServerError("Failed to remove permission from the role: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		rolePermissionsRepository.logger.Error("No permission has been removed from the role",
			zap.Int("role_id", rolePermissionsPayload.RoleID),
			zap.Int("permission_id", rolePermissionsPayload.PermissionID),
		)

		return utils.NotFound("Permission with such id not found")
	}

	rolePermissionsRepository.logger.Info("Successfully removed the permission from the role",
		zap.Int("role_id", rolePermissionsPayload.RoleID),
		zap.Int("permission_id", rolePermissionsPayload.PermissionID),
	)

	return nil
}

func (rolePermissionsRepository *RolePermissionsRepository) CheckRoleHasPermission(rolePermissionsPayload *dtos.CheckRoleHasPermissionPayload) bool {
	// create the dummy instance
	rolePermissionModel := &models.RolePermissionModel{}

	// fetch from the db
	query := "SELECT id FROM role_permissions WHERE role_id = ? AND permission_id = ?"

	queryErr := rolePermissionsRepository.db.QueryRow(query, rolePermissionsPayload.RoleID, rolePermissionsPayload.PermissionID).Scan(&rolePermissionModel.ID)

	if queryErr != nil {
		rolePermissionsRepository.logger.Error("Role does not have the permission",
			zap.Int("role_id", rolePermissionsPayload.RoleID),
			zap.Int("permission_id", rolePermissionsPayload.PermissionID),
		)

		return false
	}

	rolePermissionsRepository.logger.Info("Role has the particular permission",
		zap.Int("role_id", rolePermissionsPayload.RoleID),
		zap.Int("permission_id", rolePermissionsPayload.PermissionID),
	)

	return true
}

func NewRolePermissionsRepository(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) RolePermissionsRepositoryInterface {
	newRolePermissionsRepository := &RolePermissionsRepository{
		db:           db,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newRolePermissionsRepository
}
