package routes

import (
	"go-backend/middleware"

	"github.com/gorilla/mux"
)

// InitRoutes initializes all API routes and middlewares.
func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Global middlewares
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.RateLimiter)

	// Register route modules
	RegisterHealthRoutes(router)
	RegisterSwaggerRoutes(router)
	RegisterAuthRoutes(router)
	RegisterProtectedRoutes(router)

	return router
}
