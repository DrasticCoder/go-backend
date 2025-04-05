package routes

import (
	"net/http"

	"go-backend/controllers"
	"go-backend/middleware"

	"github.com/gorilla/mux"
)

// RegisterProtectedRoutes sets up all JWT + RBAC-protected endpoints.
func RegisterProtectedRoutes(router *mux.Router) {
	api := router.PathPrefix("/api/v1").Subrouter()

	// Apply JWT middleware to all below routes
	api.Use(middleware.JWTMiddleware)

	// Example RBAC-based route
	api.Handle("/users/profile", middleware.RBAC("admin", "premium")(http.HandlerFunc(controllers.UserProfile))).Methods(http.MethodGet)

	// Example generic protected route
	api.Handle("/protected", http.HandlerFunc(controllers.ProtectedRoute)).Methods(http.MethodGet)
}
