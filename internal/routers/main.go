package routers

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type RouterInterface interface {
	Register(r *chi.Mux)
}

func RegisterRouters(logger *zap.Logger, db *sql.DB) *chi.Mux {
	// create the router instance
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// register health router
	//SetupHealthRouter(router)

	// register user router
	newUserRouter := NewUserRouter(
		controllers.NewUserController(
			services.NewUserService(
				repositories.NewUserRepository(logger, db), logger,
			), logger,
		), logger,
	)

	newUserRouter.Register(router)

	return router
}
