package controllers

import (
	"encoding/json"
	"go-backend/config"
	"go-backend/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// UserProfile godoc
// @Summary Get User Profile
// @Description Returns the profile of the logged-in user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/users/profile [get]
func UserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// In a real-world app, you would extract user data from JWT claims here
	profile := map[string]string{
		"user": "John Doe",
		"role": "premium",
	}
	json.NewEncoder(w).Encode(profile)
}

// UserList godoc
// @Summary List all users
// @Description Returns a list of all users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {string} string "Unauthorized"
// @Router /api/v1/users [get]
func UserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// In a real-world app, you would fetch users from a database here
	users := []models.User{}
	cursor, err := config.UserCollection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	if err := cursor.All(r.Context(), &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

