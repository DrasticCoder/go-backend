package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-backend/internal/testutils"

	"github.com/stretchr/testify/assert"
)

func TestRBAC_AllowsAdmin(t *testing.T) {
	handler := RBAC("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = testutils.WithMockedUser(req, "admin", "123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRBAC_DeniesUser(t *testing.T) {
	handler := RBAC("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = testutils.WithMockedUser(req, "free", "123")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
