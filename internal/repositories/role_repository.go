package repositories

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
)

type RoleRepositoryInterface interface {
	GetRoleById(userPayload *dtos.GetUserByIdParams) (*models.UserModel, *utils.AppError)
}
