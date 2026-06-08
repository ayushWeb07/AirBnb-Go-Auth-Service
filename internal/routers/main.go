package routers

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
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

func RegisterRouters(logger *zap.Logger, db *sql.DB, serverConfig *config.ServerConfig) *chi.Mux {
	// create the router instance
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// register health router
	//SetupHealthRouter(router)

	// register user router
	userRepository := repositories.NewUserRepository(logger, db, serverConfig)
	userService := services.NewUserService(userRepository, logger, serverConfig)
	userController := controllers.NewUserController(userService, logger, serverConfig)
	userRouter := NewUserRouter(userController, logger, serverConfig)

	userRouter.Register(router)

	return router
}
