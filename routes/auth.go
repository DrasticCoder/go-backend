package routes

import (
	"net/http"

	"go-backend/controllers"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes sets up public auth routes like register and login.
func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/auth/register", controllers.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/auth/login", controllers.LoginUser).Methods(http.MethodPost)
}
