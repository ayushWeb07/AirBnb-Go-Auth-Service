package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	UserController controllers.UserControllerInterface
}

func (ur *UserRouter) Register(r *chi.Mux) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/create", ur.UserController.CreateUser)
	})
}

func NewUserRouter(controller controllers.UserControllerInterface) RouterInterface {
	newUserRouter := &UserRouter{
		UserController: controller,
	}

	return newUserRouter
}
