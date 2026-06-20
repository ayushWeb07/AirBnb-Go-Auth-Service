package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserRolesRouter struct {
	UserRolesController controllers.UserRolesControllerInterface
	logger              *zap.Logger
	serverConfig        *config.ServerConfig
}

func (userRolesRouter *UserRolesRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/user-roles", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.GetRolesOfUserPayload]).Get("/", userRolesRouter.UserRolesController.GetRolesOfUser)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.AssignRoleToUserPayload]).Post("/assign", userRolesRouter.UserRolesController.AssignRoleToUser)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.RemoveRoleFromUserPayload]).Post("/remove", userRolesRouter.UserRolesController.RemoveRoleFromUser)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckUserHasRolePayload]).Get("/check", userRolesRouter.UserRolesController.CheckUserHasRole)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckUserHasAllRolesPayload]).Get("/check-all", userRolesRouter.UserRolesController.CheckUserHasAllRoles)
	})
}

func NewUserRolesRouter(controller controllers.UserRolesControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newUserRolesRouter := &UserRolesRouter{
		UserRolesController: controller,
		logger:              logger,
		serverConfig:        serverConfig,
	}

	return newUserRolesRouter
}
