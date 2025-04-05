package routes

import (
	"net/http"

	"go-backend/controllers"
	_ "go-backend/docs" // Swagger docs; generated with swag init
	"go-backend/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// InitRoutes configures all API routes including public, auth and protected routes.
func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Global middlewares
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.RateLimiter)

	// Swagger docs route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Health Check
	router.HandleFunc("/api/v1/health", controllers.HealthCheck).Methods(http.MethodGet)

	// Auth endpoints (public)
	router.HandleFunc("/api/v1/auth/register", controllers.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/auth/login", controllers.LoginUser).Methods(http.MethodPost)

	// Protected endpoints (JWT & RBAC required)
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.JwtVerify)
	api.Handle("/users/profile", middleware.RBAC("admin", "premium")(http.HandlerFunc(controllers.UserProfile))).Methods(http.MethodGet)

	return router
}
