package middleware

import (
	"net/http"
	"strings"
)

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Remove "Bearer " prefix
		actualToken := strings.TrimPrefix(token, "Bearer ")
		// Dummy validation for demonstration:
		if actualToken != "validtoken" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Token is valid, proceed
		next.ServeHTTP(w, r)
	})
}
