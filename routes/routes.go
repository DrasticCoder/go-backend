package routes

import (
	"net/http"

	"go-backend/controllers"
	_ "go-backend/docs" // Swagger docs; generated with swag init
	"go-backend/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Global middlewares
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.RateLimiter)

	// Swagger documentation route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Public endpoint: Health Check
	router.HandleFunc("/api/v1/health", controllers.HealthCheck).Methods(http.MethodGet)

	// Protected endpoints: JWT verification and RBAC enforced
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.JwtVerify)
	// For example, only "admin" and "premium" roles can access the user profile
	api.Handle("/users/profile", middleware.RBAC("admin", "premium")(http.HandlerFunc(controllers.UserProfile))).Methods(http.MethodGet)

	return router
}
