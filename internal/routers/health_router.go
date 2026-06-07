package routers

import (
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/controllers"
	"github.com/go-chi/chi/v5"
)

func SetupHealthRouter(router *chi.Mux) {
	router.Route("/health", func(router chi.Router) {
		router.Get("/", controllers.CheckHealthStatus)
	})
}
