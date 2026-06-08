package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
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
		r.Get("/create", ur.UserController.CreateUser)
		r.Get("/id", ur.UserController.GetUserById)
		r.Get("/", ur.UserController.GetAllUsers)
		r.Get("/delete", ur.UserController.DeleteUserById)
		r.Get("/login", ur.UserController.LoginUser)
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
