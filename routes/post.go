package routes

import (
	"net/http"

	"go-backend/controllers"
	"go-backend/middleware"

	"github.com/gorilla/mux"
)

func RegisterPostRoutes(router *mux.Router) {

	// Public routes
	router.HandleFunc("", controllers.ListPosts).Methods(http.MethodGet)
	
	router.HandleFunc("/my", controllers.ListMyPosts).Methods(http.MethodGet)

	router.HandleFunc("/{id}", controllers.GetPost).Methods(http.MethodGet)

	// Protected (author)
	router.Use(middleware.JWTMiddleware)
	router.HandleFunc("", controllers.CreatePost).Methods(http.MethodPost)
	router.HandleFunc("/{id}", controllers.UpdatePost).Methods(http.MethodPut)
	router.HandleFunc("/{id}", controllers.DeletePost).Methods(http.MethodDelete)

	// TODO: Add utility routes
	// router.HandleFunc("/{id}/publish", controllers.ForcePublishPost).Methods(http.MethodPost)
	// router.HandleFunc("/upload", controllers.UploadMedia).Methods(http.MethodPost)
}
