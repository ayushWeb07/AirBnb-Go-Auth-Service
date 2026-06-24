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

type RoleRouter struct {
	RoleController controllers.RoleControllerInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (roleRouter *RoleRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/roles", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CreateRolePayload]).Post("/", roleRouter.RoleController.CreateRole)

		r.Get("/", roleRouter.RoleController.GetAllRoles)

		r.With(middlewares.DecodeAndValidateParams[dtos.GetRoleByIdParams](
			func(req *http.Request) (*dtos.GetRoleByIdParams, *utils.AppError) {
				roleId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("Role id must be provided in integer: " + err.Error())
				}

				return &dtos.GetRoleByIdParams{
					ID: roleId,
				}, nil
			},
		)).Get("/{id}", roleRouter.RoleController.GetRoleById)

		r.With(middlewares.DecodeAndValidateParams[dtos.UpdateRoleByIdParams](
			func(req *http.Request) (*dtos.UpdateRoleByIdParams, *utils.AppError) {
				roleId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("Role id must be provided in integer: " + err.Error())
				}

				return &dtos.UpdateRoleByIdParams{
					ID: roleId,
				}, nil
			},
		)).With(middlewares.DecodeAndValidateRequestBody[dtos.UpdateRoleByIdPayload]).Put("/{id}", roleRouter.RoleController.UpdateRoleById)

		r.With(middlewares.DecodeAndValidateParams[dtos.DeleteRoleByIdParams](
			func(req *http.Request) (*dtos.DeleteRoleByIdParams, *utils.AppError) {
				roleId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("Role id must be provided in integer: " + err.Error())
				}

				return &dtos.DeleteRoleByIdParams{
					ID: roleId,
				}, nil
			},
		)).Delete("/{id}", roleRouter.RoleController.DeleteRoleById)
	})
}

func NewRoleRouter(controller controllers.RoleControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newRoleRouter := &RoleRouter{
		RoleController: controller,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newRoleRouter
}
