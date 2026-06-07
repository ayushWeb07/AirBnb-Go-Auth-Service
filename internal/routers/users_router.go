package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserRouter struct {
	UserController controllers.UserControllerInterface
	logger         *zap.Logger
}

func (ur *UserRouter) Register(r *chi.Mux) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/create", ur.UserController.CreateUser)
		r.Get("/id", ur.UserController.GetUserById)
		r.Get("/", ur.UserController.GetAllUsers)
	})
}

func NewUserRouter(controller controllers.UserControllerInterface, logger *zap.Logger) RouterInterface {
	newUserRouter := &UserRouter{
		UserController: controller,
		logger:         logger,
	}

	return newUserRouter
}
