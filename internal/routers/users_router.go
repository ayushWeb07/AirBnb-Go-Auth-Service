package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserRouter struct {
	UserController controllers.UserControllerInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (ur *UserRouter) Register(r *chi.Mux) {
	r.Route("/users", func(r chi.Router) {
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.CreateUser]).Post("/register", ur.UserController.CreateUser)
		r.With(middlewares.DecodeAndValidateRequestBody[dtos.LoginUser]).Post("/login", ur.UserController.LoginUser)
		r.Get("/", ur.UserController.GetAllUsers)
		r.With(middlewares.DecodeAndValidateParams[dtos.GetUserById]).Get("/{id}", ur.UserController.GetUserById)
		r.With(middlewares.DecodeAndValidateParams[dtos.DeleteUserById]).Delete("/{id}", ur.UserController.DeleteUserById)
	})
}

func NewUserRouter(controller controllers.UserControllerInterface, logger *zap.Logger, serverConfig *config.ServerConfig) RouterInterface {
	newUserRouter := &UserRouter{
		UserController: controller,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newUserRouter
}
