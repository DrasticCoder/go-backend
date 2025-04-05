package main

import (
	"log"
	"net/http"

	"go-backend/config"
	"go-backend/routes"
)

func main() {
	// Initialize logger and cache (Redis)
	config.InitLogger()
	config.InitCache()

	router := routes.InitRoutes()
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
