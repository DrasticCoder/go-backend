package testutils

import (
	"context"
	"net/http"
)

func WithMockedUser(r *http.Request, role string, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), "role", role)
	ctx = context.WithValue(ctx, "user_id", userID)
	return r.WithContext(ctx)
}
