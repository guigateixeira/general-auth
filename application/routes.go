package application

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/guigateixeira/general-auth/handler"
	"github.com/guigateixeira/general-auth/middlewares"
)

func loadRoutes(userHandler *handler.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middlewares.SanitizeInputMiddleware)

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

	v1Router.Route("/users", func(router chi.Router) {
		loadUserRoutes(router, userHandler)
	})

	router.Mount("/v1", v1Router)

	return router
}

func loadUserRoutes(router chi.Router, userHandler *handler.UserHandler) {
	router.Post("/signup", middlewares.EmailValidatorMiddleware(http.HandlerFunc(userHandler.CreateUser)))
	router.Post("/signin", middlewares.EmailValidatorMiddleware(http.HandlerFunc(userHandler.SignIn)))
}
