package routes

import (
	"net/http"

	"go-backend/controllers"

	"github.com/gorilla/mux"
)

// RegisterHealthRoutes sets up the /health endpoint.
func RegisterHealthRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/health", controllers.HealthCheck).Methods(http.MethodGet)
}
