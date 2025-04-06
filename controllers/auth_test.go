package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-backend/config"
	"go-backend/models"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	// Cleanup test user
	config.UserCollection.DeleteMany(nil, map[string]string{"email": "test@example.com"})
	payload := models.User{
		Name:     "testuser",
		Email:    "test@example.com",
		Password: "Pa$$w0rd",
		Role:     "admin",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	RegisterUser(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "test@example.com")
}

func TestLoginUser(t *testing.T) {
	payload := map[string]string{
		"email":    "test@example.com",
		"password": "Pa$$w0rd",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	LoginUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "token")
}
