package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	role := r.Context().Value("role")

	fmt.Println("👤 user_id from context:", userID)
	fmt.Println("🔐 role from context:", role)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "You accessed a protected route!",
		"user_id": userID,
		"role":    role,
	})
}
