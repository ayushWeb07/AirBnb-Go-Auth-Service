package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type RolePermissionsRouter struct {
	RolePermissionsController controllers.RolePermissionsControllerInterface
	logger                    *zap.Logger
	serverConfig              *config.ServerConfig
}

func (rolePermissionsRouter *RolePermissionsRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/role-permissions", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.GetPermissionsOfRolePayload]).Get("/", rolePermissionsRouter.RolePermissionsController.GetPermissionsOfRole)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.AssignPermissionToRolePayload]).Post("/assign", rolePermissionsRouter.RolePermissionsController.AssignPermissionToRole)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.RemovePermissionFromRolePayload]).Post("/remove", rolePermissionsRouter.RolePermissionsController.RemovePermissionFromRole)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckRoleHasPermissionPayload]).Get("/check", rolePermissionsRouter.RolePermissionsController.CheckRoleHasPermission)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.GetPermissionsOfUserPayload]).Get("/user", rolePermissionsRouter.RolePermissionsController.GetPermissionsOfUser)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckUserHasPermissionPayload]).Get("/user/check", rolePermissionsRouter.RolePermissionsController.CheckUserHasPermission)
	})
}

func NewRolePermissionsRouter(controller controllers.RolePermissionsControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newRolePermissionsRouter := &RolePermissionsRouter{
		RolePermissionsController: controller,
		logger:                    logger,
		serverConfig:              serverConfig,
	}

	return newRolePermissionsRouter
}
