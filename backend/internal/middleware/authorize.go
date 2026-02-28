package middleware

import (
	"net/http"

	"order-express/backend/internal/handlers"
)

func RequireAdmin() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r)
			if claims == nil {
				handlers.Unauthorized(w)
				return
			}
			if claims.Role != "admin" {
				handlers.Forbidden(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
