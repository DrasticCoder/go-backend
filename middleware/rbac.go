package middleware

import (
	"net/http"
	"strings"
)

func RBAC(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For demo purposes, assume the role is passed in header "X-User-Role"
			role := r.Header.Get("X-User-Role")
			if role == "" {
				http.Error(w, "Forbidden: no role provided", http.StatusForbidden)
				return
			}
			allowed := false
			for _, allowedRole := range allowedRoles {
				if strings.EqualFold(role, allowedRole) {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Forbidden: insufficient privileges", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
