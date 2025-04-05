package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func RBAC(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role")
			fmt.Println("🔐 role from context:", role)
			if role == nil {
				http.Error(w, "🔐 Forbidden: no role found in context", http.StatusForbidden)
				return
			}

			roleStr, ok := role.(string)
			if !ok {
				http.Error(w, "🔐 Forbidden: invalid role type", http.StatusForbidden)
				return
			}

			allowed := false
			for _, allowedRole := range allowedRoles {
				if strings.EqualFold(roleStr, allowedRole) {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "🔐 Forbidden: insufficient privileges", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
