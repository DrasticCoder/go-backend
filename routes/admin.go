package routes

import (
	"net/http"

	"go-backend/controllers"

	"github.com/gorilla/mux"
)

// RegisterAdminRoutes sets up the admin endpoints for user management.
func RegisterAdminRoutes(router *mux.Router) {
	router.HandleFunc("/users", controllers.ListUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", controllers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users", controllers.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods(http.MethodDelete)
}
