package controllers

import (
	"encoding/json"
	"net/http"
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
	// For demo, we return static profile information
	profile := map[string]string{
		"user": "John Doe",
		"role": "premium",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
