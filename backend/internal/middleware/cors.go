package middleware

import (
	"net/http"
	"strings"

	"order-express/backend/internal/config"
)

// CORS allows the frontend (typically the Vue dev server) to call the API across origins.
// In production, prefer same-origin deployments (e.g. Nginx reverse proxy) and keep
// the allowlist minimal.
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	allowedOrigins := make([]string, 0, len(cfg.CORSAllowedOrigins))
	allowAllOrigins := false
	for _, origin := range cfg.CORSAllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		if origin == "*" {
			allowAllOrigins = true
			continue
		}
		allowedOrigins = append(allowedOrigins, origin)
	}

	isAllowed := func(origin string) (allowed bool, wildcard bool) {
		if origin == "" {
			return false, false
		}
		if allowAllOrigins {
			return true, true
		}
		for _, o := range allowedOrigins {
			if o == origin {
				return true, false
			}
		}
		return false, false
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := strings.TrimSpace(r.Header.Get("Origin"))
			allowed, wildcard := isAllowed(origin)

			if allowed {
				if wildcard {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")
				}

				if cfg.CORSAllowCredentials && !wildcard {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
