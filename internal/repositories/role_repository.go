package repositories

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type RoleRepositoryInterface interface {
	CreateRole(rolePayload *dtos.CreateRolePayload) *utils.AppError
	GetAllRoles() ([]*models.RoleModel, *utils.AppError)
	GetRoleById(roleParams *dtos.GetRoleByIdParams) (*models.RoleModel, *utils.AppError)
	UpdateRoleById(roleParams *dtos.UpdateRoleByIdParams, rolePayload *dtos.UpdateRoleByIdPayload) *utils.AppError
	DeleteRoleById(roleParams *dtos.DeleteRoleByIdParams) *utils.AppError
}

type RoleRepository struct {
	db           *sql.DB
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (roleRepository *RoleRepository) CreateRole(rolePayload *dtos.CreateRolePayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO roles (name, description) VALUES (?, ?)"
	result, queryExecErr := roleRepository.db.Exec(query, rolePayload.Name, rolePayload.Description)

	if queryExecErr != nil {
		roleRepository.logger.Error("Failed to insert role into the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to insert role into the database: " + queryExecErr.Error())
	}

	id, insertErr := result.LastInsertId()

	if insertErr != nil {
		roleRepository.logger.Error("Failed to insert role into the database",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to insert role into the database: " + insertErr.Error())
	}

	roleRepository.logger.Info("Successfully inserted role into the database",
		zap.Int64("role_id", id))

	return nil
}

func (roleRepository *RoleRepository) GetAllRoles() ([]*models.RoleModel, *utils.AppError) {
	// create the dummy instance
	var roleModels []*models.RoleModel

	// load the rows
	query := "SELECT * FROM roles"
	rows, queryErr := roleRepository.db.Query(query)

	if queryErr != nil {
		roleRepository.logger.Error("Something went wrong while fetching all the roles",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the roles: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		roleModel := &models.RoleModel{}

		rowScanErr := rows.Scan(&roleModel.ID, &roleModel.Name, &roleModel.Description, &roleModel.CreatedAt, &roleModel.UpdatedAt)

		if rowScanErr != nil {
			roleRepository.logger.Error("Failed to fetch all the roles from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching all the roles: " + rowScanErr.Error())
		}

		roleModels = append(roleModels, roleModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		roleRepository.logger.Error("Failed to fetch all the roles from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching all the roles: " + rowsErr.Error())
	}

	roleRepository.logger.Info("Successfully fetched all the roles from the database",
		zap.Int("count", len(roleModels)))

	return roleModels, nil
}

func (roleRepository *RoleRepository) GetRoleById(roleParams *dtos.GetRoleByIdParams) (*models.RoleModel, *utils.AppError) {
	// create the dummy instance
	roleModel := &models.RoleModel{}

	// fetch from the db
	query := "SELECT * FROM roles WHERE id = ?"

	queryErr := roleRepository.db.QueryRow(query, roleParams.ID).Scan(&roleModel.ID, &roleModel.Name, &roleModel.Description, &roleModel.CreatedAt, &roleModel.UpdatedAt)

	if queryErr != nil {
		if queryErr == sql.ErrNoRows {
			roleRepository.logger.Error("Such role not found",
				zap.Int("role_id", roleParams.ID))

			return nil, utils.NotFound("Role with such id not found")
		}

		roleRepository.logger.Error("Failed to fetch the role from the database",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Failed to fetch the role from the database: " + queryErr.Error())
	}

	roleRepository.logger.Info("Successfully fetched the role from the database",
		zap.Int("role_id", roleModel.ID),
	)

	return roleModel, nil
}

func (roleRepository *RoleRepository) UpdateRoleById(roleParams *dtos.UpdateRoleByIdParams, rolePayload *dtos.UpdateRoleByIdPayload) *utils.AppError {
	// prepare and execute the query
	query := "UPDATE roles SET name = ?, description = ? WHERE id = ?"
	result, queryExecErr := roleRepository.db.Exec(query, rolePayload.Name, rolePayload.Description, roleParams.ID)

	if queryExecErr != nil {
		roleRepository.logger.Error("Failed to update role from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to update role from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		roleRepository.logger.Error("Failed to update role from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to update role from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		roleRepository.logger.Error("No role has been updated from the database",
			zap.Int("role_id", roleParams.ID))

		return utils.NotFound("Role with such id not found")
	}

	roleRepository.logger.Info("Successfully updated the role from the database",
		zap.Int("role_id", roleParams.ID))

	return nil
}

func (roleRepository *RoleRepository) DeleteRoleById(roleParams *dtos.DeleteRoleByIdParams) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM roles WHERE id = ?"
	result, queryExecErr := roleRepository.db.Exec(query, roleParams.ID)

	if queryExecErr != nil {
		roleRepository.logger.Error("Failed to delete role from the database",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to delete role from the database: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		roleRepository.logger.Error("Failed to delete role from the database",
			zap.String("error", rowsErr.Error()))

		return utils.InternalServerError("Failed to delete role from the database: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		roleRepository.logger.Error("No role has been deleted from the database",
			zap.Int("role_id", roleParams.ID))

		return utils.NotFound("Role with such id not found")
	}

	roleRepository.logger.Info("Successfully deleted the role from the database",
		zap.Int("role_id", roleParams.ID))

	return nil
}

func NewRoleRepository(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) RoleRepositoryInterface {
	newRoleRepository := &RoleRepository{
		db:           db,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newRoleRepository
}
