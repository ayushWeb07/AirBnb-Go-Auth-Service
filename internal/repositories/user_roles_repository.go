package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"go.uber.org/zap"
)

type UserRolesRepositoryInterface interface {
	GetRolesOfUser(userRolesPayload *dtos.GetRolesOfUserPayload) ([]*models.RoleModel, *utils.AppError)
	AssignRoleToUser(userRolesPayload *dtos.AssignRoleToUserPayload) *utils.AppError
	RemoveRoleFromUser(userRolesPayload *dtos.RemoveRoleFromUserPayload) *utils.AppError
	CheckUserHasRole(userRolesPayload *dtos.CheckUserHasRolePayload) bool

	CheckUserHasAllRoles(userRolesPayload *dtos.CheckUserHasAllRolesPayload) bool
}

type UserRolesRepository struct {
	db           *sql.DB
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (userRolesRepository *UserRolesRepository) GetRolesOfUser(userRolesPayload *dtos.GetRolesOfUserPayload) ([]*models.RoleModel, *utils.AppError) {
	// create the dummy instance
	var roleModels []*models.RoleModel

	// load the rows
	query := "SELECT r.* FROM user_roles ur INNER JOIN roles r ON r.id = ur.role_id WHERE ur.user_id = ?"
	rows, queryErr := userRolesRepository.db.Query(query, userRolesPayload.UserID)

	if queryErr != nil {
		userRolesRepository.logger.Error("Something went wrong while fetching the roles",
			zap.String("error", queryErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching the roles: " + queryErr.Error())
	}

	defer rows.Close()

	// loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		roleModel := &models.RoleModel{}

		rowScanErr := rows.Scan(&roleModel.ID, &roleModel.Name, &roleModel.Description, &roleModel.CreatedAt, &roleModel.UpdatedAt)

		if rowScanErr != nil {
			userRolesRepository.logger.Error("Failed to fetch the roles from the database",
				zap.String("error", rowScanErr.Error()))

			return nil, utils.InternalServerError("Something went wrong while fetching the roles: " + rowScanErr.Error())
		}

		roleModels = append(roleModels, roleModel)
	}

	rowsErr := rows.Err()

	if rowsErr != nil {
		userRolesRepository.logger.Error("Failed to fetch the roles from the database",
			zap.String("error", rowsErr.Error()))

		return nil, utils.InternalServerError("Something went wrong while fetching the roles: " + rowsErr.Error())
	}

	userRolesRepository.logger.Info("Successfully fetched the roles from the database",
		zap.Int("count", len(roleModels)))

	return roleModels, nil
}

func (userRolesRepository *UserRolesRepository) AssignRoleToUser(userRolesPayload *dtos.AssignRoleToUserPayload) *utils.AppError {
	// insert into the db
	query := "INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)"
	result, queryExecErr := userRolesRepository.db.Exec(query, userRolesPayload.UserID, userRolesPayload.RoleID)

	if queryExecErr != nil {
		userRolesRepository.logger.Error("Failed to assign role to the user",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to assign role to the user: " + queryExecErr.Error())
	}

	_, insertErr := result.LastInsertId()

	if insertErr != nil {
		userRolesRepository.logger.Error("Failed to assign role to the user",
			zap.String("error", insertErr.Error()))

		return utils.InternalServerError("Failed to assign role to the user: " + insertErr.Error())
	}

	userRolesRepository.logger.Info("Successfully assigned role to the user",
		zap.Int("user_id", userRolesPayload.UserID),
		zap.Int("role_id", userRolesPayload.RoleID),
	)

	return nil
}

func (userRolesRepository *UserRolesRepository) RemoveRoleFromUser(userRolesPayload *dtos.RemoveRoleFromUserPayload) *utils.AppError {
	// prepare and execute the query
	query := "DELETE FROM user_roles WHERE user_id = ? AND role_id = ?"
	result, queryExecErr := userRolesRepository.db.Exec(query, userRolesPayload.UserID, userRolesPayload.RoleID)

	if queryExecErr != nil {
		userRolesRepository.logger.Error("Failed to remove role from the user",
			zap.String("error", queryExecErr.Error()))

		return utils.InternalServerError("Failed to remove role from the user: " + queryExecErr.Error())
	}

	// check if any rows got affected
	rowsAffected, rowsErr := result.RowsAffected()

	if rowsErr != nil {
		userRolesRepository.logger.Error("Failed to remove role from the user",
			zap.Int("user_id", userRolesPayload.UserID),
			zap.Int("role_id", userRolesPayload.RoleID),
		)

		return utils.InternalServerError("Failed to remove role from the user: " + rowsErr.Error())
	}

	if rowsAffected == 0 {
		userRolesRepository.logger.Error("No role has been removed from the user",
			zap.Int("user_id", userRolesPayload.UserID),
			zap.Int("role_id", userRolesPayload.RoleID),
		)

		return utils.NotFound("Role with such id not found")
	}

	userRolesRepository.logger.Info("Successfully removed the role from the user",
		zap.Int("user_id", userRolesPayload.UserID),
		zap.Int("role_id", userRolesPayload.RoleID),
	)

	return nil
}

func (userRolesRepository *UserRolesRepository) CheckUserHasRole(userRolesPayload *dtos.CheckUserHasRolePayload) bool {
	// create the dummy instance
	userRoleModel := &models.UserRoleModel{}

	// fetch from the db
	query := "SELECT id FROM user_roles WHERE user_id = ? AND role_id = ?"

	queryErr := userRolesRepository.db.QueryRow(query, userRolesPayload.UserID, userRolesPayload.RoleID).Scan(&userRoleModel.ID)

	if queryErr != nil {
		userRolesRepository.logger.Error("User does not have the role",
			zap.Int("user_id", userRolesPayload.UserID),
			zap.Int("role_id", userRolesPayload.RoleID),
		)

		return false
	}

	userRolesRepository.logger.Info("User has the particular role",
		zap.Int("user_id", userRolesPayload.UserID),
		zap.Int("role_id", userRolesPayload.RoleID),
	)

	return true
}

func (userRolesRepository *UserRolesRepository) CheckUserHasAllRoles(userRolesPayload *dtos.CheckUserHasAllRolesPayload) bool {
	count := 0

	// build the query
	placeholders := make([]string, len(userRolesPayload.RoleNames))
	args := make([]any, 0, len(userRolesPayload.RoleNames)+1)

	args = append(args, userRolesPayload.UserID)

	for i, role := range userRolesPayload.RoleNames {
		placeholders[i] = "?"
		args = append(args, role)
	}

	roleNamesStr := strings.Join(placeholders, ",")

	query := fmt.Sprintf(
		`
		SELECT COUNT(*)
		FROM user_roles ur
		INNER JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ?
		AND r.name IN (%s)
		`, roleNamesStr)

	queryErr := userRolesRepository.db.QueryRow(query, args...).Scan(&count)

	if queryErr != nil {
		userRolesRepository.logger.Error("User does not have all the required roles",
			zap.Int("user_id", userRolesPayload.UserID),
			zap.String("role_names", roleNamesStr),
		)

		return false
	}

	hasAllRoles := count == len(userRolesPayload.RoleNames)

	return hasAllRoles
}

func NewUserRolesRepository(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) UserRolesRepositoryInterface {
	newUserRolesRepository := &UserRolesRepository{
		db:           db,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newUserRolesRepository
}
