package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"order-express/backend/internal/cache"
	"order-express/backend/internal/config"
	"order-express/backend/internal/db"
	"order-express/backend/internal/middleware"
	"order-express/backend/internal/router"

	"github.com/redis/go-redis/v9"
)

func main() {
	cmd := "serve"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	cfg := config.Load()

	if cfg.CORSAllowCredentials {
		for _, origin := range cfg.CORSAllowedOrigins {
			if strings.TrimSpace(origin) == "*" {
				log.Fatalf("invalid CORS config: CORS_ALLOW_CREDENTIALS=true cannot be used with CORS_ALLOWED_ORIGINS=*")
			}
		}
	}

	if cfg.AppEnv != "local" {
		if cfg.JWTSecret == config.DefaultJWTSecret || strings.TrimSpace(cfg.JWTSecret) == "" {
			log.Fatalf("refusing to start: JWT_SECRET must be set for non-local environments")
		}
		if cfg.DB.Password == config.DefaultDBPassword || strings.TrimSpace(cfg.DB.Password) == "" {
			log.Fatalf("refusing to start: DB_PASSWORD must be set for non-local environments")
		}
		if cmd == "serve" && (cfg.DB.AutoMigrate || cfg.DB.AutoSeed) {
			log.Fatalf("refusing to start: DB_AUTO_MIGRATE/DB_AUTO_SEED must be disabled for non-local environments")
		}
	}

	gdb, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	sqlDB, err := gdb.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %v", err)
	}
	defer sqlDB.Close()

	var rdb *redis.Client
	if cfg.Redis.Addr != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("redis disabled: %v", err)
			_ = rdb.Close()
			rdb = nil
		}
	}
	if rdb != nil {
		defer rdb.Close()
	}

	c := cache.New(rdb, "oe:v1:")
	rl := middleware.NewRateLimiter(rdb, "oe:v1:", cfg.RateLimit.TrustProxy)

	switch cmd {
	case "migrate":
		if err := db.Migrate(gdb); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Printf("migration completed")
		return
	case "seed":
		if err := db.Seed(gdb); err != nil {
			log.Fatalf("seed failed: %v", err)
		}
		log.Printf("seed completed")
		return
	case "serve":
		if cfg.DB.AutoMigrate {
			if err := db.Migrate(gdb); err != nil {
				log.Fatalf("migration failed: %v", err)
			}
		}
		if cfg.DB.AutoSeed {
			if err := db.Seed(gdb); err != nil {
				log.Fatalf("seed failed: %v", err)
			}
		}
	default:
		log.Fatalf("unknown command: %s (expected: serve|migrate|seed)", cmd)
	}

	mux := router.New(gdb, cfg, c, rl)

	log.Printf("Order Express backend listening on %s", cfg.ListenAddr)
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
