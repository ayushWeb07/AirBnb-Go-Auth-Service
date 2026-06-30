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

type UserRouter struct {
	UserController      controllers.UserControllerInterface
	UserRolesController controllers.UserRolesControllerInterface
	logger              *zap.Logger
	serverConfig        *config.ServerConfig
}

func (ur *UserRouter) Register(r *chi.Mux) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CreateUserPayload]).Post("/register", ur.UserController.CreateUser)
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.LoginUserPayload]).Post("/login", ur.UserController.LoginUser)

		r.Get("/", ur.UserController.GetAllUsers)

		r.With(
			middlewares.AuthMiddleware(ur.serverConfig),
			middlewares.RequireUserAnyRoles(ur.UserRolesController.GetUserRolesService().GetUserRolesRepository(), []string{"admin", "user", "moderator"}),
		).Get("/profile", ur.UserController.GetProfile)

		r.With(middlewares.DecodeAndValidateParams[dtos.GetUserByIdParams](
			func(req *http.Request) (*dtos.GetUserByIdParams, *utils.AppError) {
				userId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("User id must be provided in integer: " + err.Error())
				}

				return &dtos.GetUserByIdParams{
					ID: userId,
				}, nil
			},
		)).Get("/{id}", ur.UserController.GetUserById)

		r.With(middlewares.DecodeAndValidateParams[dtos.DeleteUserByIdParams](
			func(req *http.Request) (*dtos.DeleteUserByIdParams, *utils.AppError) {
				userId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.BadRequest("User id must be provided in integer: " + err.Error())
				}

				return &dtos.DeleteUserByIdParams{
					ID: userId,
				}, nil
			},
		)).Delete("/{id}", ur.UserController.DeleteUserById)

		r.With(middlewares.DecodeAndValidateParams[dtos.UpdateUserByIdParams](
			func(req *http.Request) (*dtos.UpdateUserByIdParams, *utils.AppError) {
				userId, err := strconv.Atoi(chi.URLParam(req, "id"))

				if err != nil {
					return nil, utils.InternalServerError("User id must be provided in integer: " + err.Error())
				}

				return &dtos.UpdateUserByIdParams{
					ID: userId,
				}, nil
			},
		)).With(middlewares.DecodeAndValidateRequestBody[dtos.UpdateUserByIdPayload]).Put("/{id}", ur.UserController.UpdateUserById)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CreateOtpServicePayload]).Post("/send-otp-for-verification", ur.UserController.SendOtpForVerification)

		r.With(middlewares.DecodeAndValidateRequestBody[dtos.VerifyOtpPayload]).Post("/verify-otp", ur.UserController.VerifyOtp)

		r.Post("/refresh-access-token", ur.UserController.VerifyOtp)
	})
}

func NewUserRouter(controller controllers.UserControllerInterface, userRolesController controllers.UserRolesControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newUserRouter := &UserRouter{
		UserController:      controller,
		UserRolesController: userRolesController,
		logger:              logger,
		serverConfig:        serverConfig,
	}

	return newUserRouter
}
