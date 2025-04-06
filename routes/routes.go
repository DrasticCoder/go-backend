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
	adminRouter := router.PathPrefix("/api/v1/admin").Subrouter()
	adminRouter.Use(middleware.JWTMiddleware)
	adminRouter.Use(middleware.RBAC("admin"))
	RegisterAdminRoutes(adminRouter)
	RegisterProtectedRoutes(router)

	return router
}
