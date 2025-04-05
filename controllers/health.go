package controllers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck godoc
// @Summary Health Check
// @Description Returns OK if the server is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "OK"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
