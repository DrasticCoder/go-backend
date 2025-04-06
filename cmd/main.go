// @title CRM API
// @version 1.0
// @description API Backend in Go with PostgreSQL, JWT, and Swagger
// @host localhost:8080
// @BasePath /

package main

import (
	"log"
	"net/http"

	_ "go-backend/docs"

	"go-backend/config"
	"go-backend/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}
}


func main() {
	// Initialize Logger, Database and Cache in proper order
	config.InitLogger()
	config.InitDB()
	config.InitCache()

	// Auto-migrate models
	// if err := config.DB.AutoMigrate(&models.User{}); err != nil {
	// 	config.Logger.Fatalf("AutoMigrate failed: %v", err)
	// }

	// Initialize all routes
	router := routes.InitRoutes()

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
