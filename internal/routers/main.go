package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterInterface interface {
	Register(r *chi.Mux)
}

func RegisterRouters() *chi.Mux {
	// create the router instance
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// register health router
	//SetupHealthRouter(router)

	// register user router
	newUserRouter := NewUserRouter(controllers.NewUserController(services.NewUserService(repositories.NewUserRepository())))
	newUserRouter.Register(router)

	return router
}
