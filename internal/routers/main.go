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

	// create all repositories
	userRepository := repositories.NewUserRepository(logger, db, serverConfig)
	roleRepository := repositories.NewRoleRepository(logger, db, serverConfig)
	permissionRepository := repositories.NewPermissionRepository(logger, db, serverConfig)

	userRolesRepository := repositories.NewUserRolesRepository(logger, db, serverConfig)
	rolePermissionsRepository := repositories.NewRolePermissionsRepository(logger, db, serverConfig)

	// create all services
	userService := services.NewUserService(userRepository, logger, serverConfig)
	roleService := services.NewRoleService(roleRepository, logger, serverConfig)
	permissionService := services.NewPermissionService(permissionRepository, logger, serverConfig)

	userRolesService := services.NewUserRolesService(userRolesRepository, userRepository, roleRepository, logger, serverConfig)
	rolePermissionsService := services.NewRolePermissionsService(rolePermissionsRepository, roleRepository, permissionRepository, userRepository, logger, serverConfig)

	// create all controllers
	userController := controllers.NewUserController(userService, logger, serverConfig)
	roleController := controllers.NewRoleController(roleService, logger, serverConfig)
	permissionController := controllers.NewPermissionController(permissionService, logger, serverConfig)

	userRolesController := controllers.NewUserRolesController(userRolesService, logger, serverConfig)
	rolePermissionsController := controllers.NewRolePermissionsController(rolePermissionsService, logger, serverConfig)

	// create all routers
	userRouter := NewUserRouter(userController, userRolesController, logger, serverConfig)
	roleRouter := NewRoleRouter(roleController, logger, serverConfig)
	permissionRouter := NewPermissionRouter(permissionController, logger, serverConfig)

	userRolesRouter := NewUserRolesRouter(userRolesController, logger, serverConfig)
	rolePermissionsRouter := NewRolePermissionsRouter(rolePermissionsController, logger, serverConfig)

	// register all routers
	userRouter.Register(router)
	roleRouter.Register(router)
	permissionRouter.Register(router)

	userRolesRouter.Register(router)
	rolePermissionsRouter.Register(router)

	// register the reverse proxy servers
	hotelService := utils.ProxyToService("http://localhost:3000")
	bookingService := utils.ProxyToService("http://localhost:3010")

	router.Handle("/api/v1/hotels/*", hotelService)
	router.Handle("/api/v1/bookings/*", bookingService)

	return router
}
