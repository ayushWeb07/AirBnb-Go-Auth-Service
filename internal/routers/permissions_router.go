package routers

import (
	"net/http"
	"strconv"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/middlewares"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type PermissionRouter struct {
	PermissionController controllers.PermissionControllerInterface
	logger               *zap.Logger
	serverConfig         *config.ServerConfig
}

func (permissionRouter *PermissionRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/permissions", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CreatePermissionPayload]).Post("/", permissionRouter.PermissionController.CreatePermission)

		r.Get("/", permissionRouter.PermissionController.GetAllPermissions)

		r.With(middlewares.DecodeAndValidateParams[dtos.GetPermissionByIdParams](
			func(req *http.Request) (*dtos.GetPermissionByIdParams, *utils.AppError) {
				permissionId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("Permission id must be provided in integer: " + err.Error())
				}

				return &dtos.GetPermissionByIdParams{
					ID: permissionId,
				}, nil
			},
		)).Get("/{id}", permissionRouter.PermissionController.GetPermissionById)

		r.With(middlewares.DecodeAndValidateParams[dtos.UpdatePermissionByIdParams](
			func(req *http.Request) (*dtos.UpdatePermissionByIdParams, *utils.AppError) {
				permissionId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("Permission id must be provided in integer: " + err.Error())
				}

				return &dtos.UpdatePermissionByIdParams{
					ID: permissionId,
				}, nil
			},
		)).With(middlewares.DecodeAndValidateRequestBody[dtos.UpdatePermissionByIdPayload]).Put("/{id}", permissionRouter.PermissionController.UpdatePermissionById)

		r.With(middlewares.DecodeAndValidateParams[dtos.DeletePermissionByIdParams](
			func(req *http.Request) (*dtos.DeletePermissionByIdParams, *utils.AppError) {
				permissionId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("Permission id must be provided in integer: " + err.Error())
				}

				return &dtos.DeletePermissionByIdParams{
					ID: permissionId,
				}, nil
			},
		)).Delete("/{id}", permissionRouter.PermissionController.DeletePermissionById)
	})
}

func NewPermissionRouter(controller controllers.PermissionControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newPermissionRouter := &PermissionRouter{
		PermissionController: controller,
		logger:               logger,
		serverConfig:         serverConfig,
	}

	return newPermissionRouter
}
