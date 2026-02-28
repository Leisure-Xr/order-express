package middleware

import (
	"context"
	"net/http"
	"strings"

	"order-express/backend/internal/config"
	"order-express/backend/internal/handlers"
	jwtpkg "order-express/backend/pkg/jwt"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				handlers.Unauthorized(w)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtpkg.ParseToken(cfg.JWTSecret, tokenStr)
			if err != nil {
				handlers.Unauthorized(w)
				return
			}
			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext retrieves JWT claims from the request context.
func ClaimsFromContext(r *http.Request) *jwtpkg.Claims {
	v := r.Context().Value(ClaimsKey)
	if v == nil {
		return nil
	}
	c, _ := v.(*jwtpkg.Claims)
	return c
}

// RequireAuth is a helper that extracts the token from the request directly
// (for handlers that are not wrapped by the Auth middleware).
func RequireAuth(cfg *config.Config, w http.ResponseWriter, r *http.Request) *jwtpkg.Claims {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		handlers.Unauthorized(w)
		return nil
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := jwtpkg.ParseToken(cfg.JWTSecret, tokenStr)
	if err != nil {
		handlers.Unauthorized(w)
		return nil
	}
	return claims
}
