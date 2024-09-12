package application

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/guigateixeira/general-auth/handler"
)

func loadRoutes() *chi.Mux {
	// Creates a new router for the microservice
	router := chi.NewRouter()

	// Use logger middleware
	router.Use(middleware.Logger)

	// Use CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	// Starndard health check endpoint
	v1Router.Get("/health", handler.HandlerReadiness)

	v1Router.Route("/orders", loadOrderRoutes)

	router.Mount("/v1", v1Router)

	return router
}

func loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteById)
}
