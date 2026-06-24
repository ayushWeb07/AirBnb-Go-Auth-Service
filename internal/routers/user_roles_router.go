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

type UserRolesRouter struct {
	UserRolesController controllers.UserRolesControllerInterface
	logger              *zap.Logger
	serverConfig        *config.ServerConfig
}

func (userRolesRouter *UserRolesRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/user-roles", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateParams[dtos.GetRolesOfUserPayload](
			func(req *http.Request) (*dtos.GetRolesOfUserPayload, *utils.AppError) {
				userId, err := strconv.Atoi(chi.URLParam(req, "user_id"))

				if err != nil {
					return nil, utils.BadRequest("User id must be provided in integer: " + err.Error())
				}

				return &dtos.GetRolesOfUserPayload{
					UserID: userId,
				}, nil
			},
		)).Get("/user/{user_id}", userRolesRouter.UserRolesController.GetRolesOfUser)

		r.With(
			middlewares.AuthMiddleware(userRolesRouter.serverConfig),
			middlewares.RequireUserAllRoles(userRolesRouter.UserRolesController.GetUserRolesService().GetUserRolesRepository(), []string{"admin"}),
			middlewares.DecodeAndValidateRequestBody[dtos.AssignRoleToUserPayload]).Post("/assign", userRolesRouter.UserRolesController.AssignRoleToUser)

		r.With(
			middlewares.AuthMiddleware(userRolesRouter.serverConfig),
			middlewares.RequireUserAllRoles(userRolesRouter.UserRolesController.GetUserRolesService().GetUserRolesRepository(), []string{"admin"}),
			middlewares.DecodeAndValidateRequestBody[dtos.RemoveRoleFromUserPayload]).Post("/remove", userRolesRouter.UserRolesController.RemoveRoleFromUser)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckUserHasRolePayload]).Post("/check-single", userRolesRouter.UserRolesController.CheckUserHasRole)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckUserHasAllRolesPayload]).Post("/check-all", userRolesRouter.UserRolesController.CheckUserHasAllRoles)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CheckUserHasAnyRolesPayload]).Post("/check-any", userRolesRouter.UserRolesController.CheckUserHasAnyRoles)
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
