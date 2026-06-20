package routers

import (
	"database/sql"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/middlewares"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
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
	router.Use(middlewares.RateLimiter(serverConfig))

	// register health router
	//SetupHealthRouter(router)

	// register user router
	userRepository := repositories.NewUserRepository(logger, db, serverConfig)
	userService := services.NewUserService(userRepository, logger, serverConfig)
	userController := controllers.NewUserController(userService, logger, serverConfig)
	userRouter := NewUserRouter(userController, logger, serverConfig)

	userRouter.Register(router)

	// register roles router
	roleRepository := repositories.NewRoleRepository(logger, db, serverConfig)
	roleService := services.NewRoleService(roleRepository, logger, serverConfig)
	roleController := controllers.NewRoleController(roleService, logger, serverConfig)
	roleRouter := NewRoleRouter(roleController, logger, serverConfig)

	roleRouter.Register(router)

	// register permissions router
	permissionRepository := repositories.NewPermissionRepository(logger, db, serverConfig)
	permissionService := services.NewPermissionService(permissionRepository, logger, serverConfig)
	permissionController := controllers.NewPermissionController(permissionService, logger, serverConfig)
	permissionRouter := NewPermissionRouter(permissionController, logger, serverConfig)

	permissionRouter.Register(router)

	// register user roles router
	userRolesRepository := repositories.NewUserRolesRepository(logger, db, serverConfig)
	userRepository_2 := repositories.NewUserRepository(logger, db, serverConfig)
	roleRepository_2 := repositories.NewRoleRepository(logger, db, serverConfig)

	userRolesService := services.NewUserRolesService(userRolesRepository, userRepository_2, roleRepository_2, logger, serverConfig)
	userRolesController := controllers.NewUserRolesController(userRolesService, logger, serverConfig)
	userRolesRouter := NewUserRolesRouter(userRolesController, logger, serverConfig)

	userRolesRouter.Register(router)

	// register the reverse proxy servers
	hotelService := utils.ProxyToService("http://localhost:3000")
	bookingService := utils.ProxyToService("http://localhost:3010")

	router.Handle("/api/v1/hotels/*", hotelService)
	router.Handle("/api/v1/bookings/*", bookingService)

	return router
}
