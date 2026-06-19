package repositories

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type PermissionRepositoryInterface interface {
	CreatePermission(permissionPayload *dtos.CreatePermissionPayload) *utils.AppError
	GetAllPermissions() ([]*models.PermissionModel, *utils.AppError)
	GetPermissionById(permissionParams *dtos.GetPermissionByIdParams) (*models.PermissionModel, *utils.AppError)
	UpdatePermissionById(permissionParams *dtos.UpdatePermissionByIdParams, permissionPayload *dtos.UpdatePermissionByIdPayload) *utils.AppError
	DeletePermissionById(permissionParams *dtos.DeletePermissionByIdParams) *utils.AppError
}

type PermissionRepository struct {
	db           *sql.DB
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (permissionRepository *PermissionRepository) CreatePermission(permissionPayload *dtos.CreatePermissionPayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO permissions (name, description, resource, action) VALUES (?, ?, ?, ?)"
	result, queryExecErr := permissionRepository.db.Exec(query, permissionPayload.Name, permissionPayload.Description, permissionPayload.Resource, permissionPayload.Action)

	if queryExecErr != nil {
		permissionRepository.logger.Error("Failed to insert permission into the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to insert permission into the database: " + queryExecErr.Error())
	}

	id, insertErr := result.LastInsertId()

	if insertErr != nil {
		permissionRepository.logger.Error("Failed to insert permission into the database",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to insert permission into the database: " + insertErr.Error())
	}

	permissionRepository.logger.Info("Successfully inserted permission into the database",
		zap.Int64("permission_id", id))

	return nil
}

func (permissionRepository *PermissionRepository) GetAllPermissions() ([]*models.PermissionModel, *utils.AppError) {
	// create the dummy instance
	var permissionModels []*models.PermissionModel

	// load the rows
	query := "SELECT * FROM permissions"
	rows, queryErr := permissionRepository.db.Query(query)

	if queryErr != nil {
		permissionRepository.logger.Error("Something went wrong while fetching all the permissions",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the permissions: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		permissionModel := &models.PermissionModel{}

		rowScanErr := rows.Scan(&permissionModel.ID, &permissionModel.Name, &permissionModel.Description, &permissionModel.Resource, &permissionModel.Action, &permissionModel.CreatedAt, &permissionModel.UpdatedAt)

		if rowScanErr != nil {
			permissionRepository.logger.Error("Failed to fetch all the permissions from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching all the permissions: " + rowScanErr.Error())
		}

		permissionModels = append(permissionModels, permissionModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		permissionRepository.logger.Error("Failed to fetch all the permissions from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the permissions: " + rowsErr.Error())
	}

	permissionRepository.logger.Info("Successfully fetched all the permissions from the database",
		zap.Int("count", len(permissionModels)))

	return permissionModels, nil
}

func (permissionRepository *PermissionRepository) GetPermissionById(permissionParams *dtos.GetPermissionByIdParams) (*models.PermissionModel, *utils.AppError) {
	// create the dummy instance
	permissionModel := &models.PermissionModel{}

	// fetch from the db
	query := "SELECT * FROM permissions WHERE id = ?"

	queryErr := permissionRepository.db.QueryRow(query, permissionParams.ID).Scan(&permissionModel.ID, &permissionModel.Name, &permissionModel.Description, &permissionModel.Resource, &permissionModel.Action, &permissionModel.CreatedAt, &permissionModel.UpdatedAt)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			permissionRepository.logger.Error("Such permission not found",
				zap.Int("permission_id", permissionParams.ID))

			return nil, utils.NotFound("Permission with such id not found")
		}

		permissionRepository.logger.Error("Failed to fetch the permission from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the permission from the database: " + queryErr.Error())
	}

	permissionRepository.logger.Info("Successfully fetched the permission from the database",
		zap.Int("permission_id", permissionModel.ID),
	)

	return permissionModel, nil
}

func (permissionRepository *PermissionRepository) UpdatePermissionById(permissionParams *dtos.UpdatePermissionByIdParams, permissionPayload *dtos.UpdatePermissionByIdPayload) *utils.AppError {
	// prepare and execute the query
	query := "UPDATE permissions SET name = ?, description = ?, resource = ?, action = ? WHERE id = ?"
	result, queryExecErr := permissionRepository.db.Exec(query, permissionPayload.Name, permissionPayload.Description, permissionPayload.Resource, permissionPayload.Action, permissionParams.ID)

	if queryExecErr != nil {
		permissionRepository.logger.Error("Failed to update permission from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to update permission from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		permissionRepository.logger.Error("Failed to update permission from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to update permission from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		permissionRepository.logger.Error("No permission has been updated from the database",
			zap.Int("permission_id", permissionParams.ID))

		return utils.NotFound("Permission with such id not found")
	}

	permissionRepository.logger.Info("Successfully updated the permission from the database",
		zap.Int("permission_id", permissionParams.ID))

	return nil
}

func (permissionRepository *PermissionRepository) DeletePermissionById(permissionParams *dtos.DeletePermissionByIdParams) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM permissions WHERE id = ?"
	result, queryExecErr := permissionRepository.db.Exec(query, permissionParams.ID)

	if queryExecErr != nil {
		permissionRepository.logger.Error("Failed to delete permission from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to delete permission from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		permissionRepository.logger.Error("Failed to delete permission from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to delete permission from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		permissionRepository.logger.Error("No permission has been deleted from the database",
			zap.Int("permission_id", permissionParams.ID))

		return utils.NotFound("Permission with such id not found")
	}

	permissionRepository.logger.Info("Successfully deleted the permission from the database",
		zap.Int("permission_id", permissionParams.ID))

	return nil
}

func NewPermissionRepository(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) PermissionRepositoryInterface {
	newPermissionRepository := &PermissionRepository{
		db:           db,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newPermissionRepository
}
