package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-backend/config"
	"go-backend/internal/testutils"
	"go-backend/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMain(m *testing.M) {
	os.Getenv("MONGO_URI")
	os.Setenv("JWT_SECRET", "test-secret")

	config.InitDB()
	os.Exit(m.Run())
}

func seedAdmin() {
	config.UserCollection.InsertOne(context.TODO(), models.User{
		ID:        primitive.NewObjectID(),
		Name:  "adminuser",
		Email:     "admin@crm.com",
		Password:  "fakehashed",
		Role:      "admin",
		CreatedAt: time.Now(),
	})
}

func TestListUsers_AdminAccess(t *testing.T) {
	seedAdmin()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req = testutils.WithMockedUser(req, "admin", "admin-id")
	rec := httptest.NewRecorder()

	ListUsers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "adminuser")
}
