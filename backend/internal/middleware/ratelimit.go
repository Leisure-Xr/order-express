package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"order-express/backend/internal/handlers"
)

type RateLimiter struct {
	client *redis.Client
	prefix string
}

func NewRateLimiter(client *redis.Client, prefix string) *RateLimiter {
	if client == nil {
		return nil
	}
	return &RateLimiter{client: client, prefix: prefix}
}

func (rl *RateLimiter) key(name string, r *http.Request) string {
	ip := clientIP(r)
	if ip == "" {
		ip = "unknown"
	}
	return rl.prefix + "rl:" + name + ":" + ip
}

func (rl *RateLimiter) Middleware(name string, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl == nil || limit <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
			defer cancel()

			key := rl.key(name, r)
			pipe := rl.client.TxPipeline()
			incr := pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, window)
			_, err := pipe.Exec(ctx)
			if err != nil {
				// Fail open: do not block requests if Redis is unavailable.
				next.ServeHTTP(w, r)
				return
			}

			if incr.Val() > int64(limit) {
				handlers.Fail(w, http.StatusTooManyRequests, "too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func ParseLimitPerMinute(value int) (limit int, window time.Duration) {
	if value <= 0 {
		return 0, time.Minute
	}
	return value, time.Minute
}

func ParseLimitFromHeader(r *http.Request, header string) int {
	v := strings.TrimSpace(r.Header.Get(header))
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}
