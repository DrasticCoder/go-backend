package routes

import (
	_ "go-backend/docs" // Swagger docs

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// RegisterSwaggerRoutes sets up the Swagger UI route.
func RegisterSwaggerRoutes(router *mux.Router) {
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
