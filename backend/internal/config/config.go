package config

import (
	"os"
	"strconv"
	"strings"
)

type DBConfig struct {
	Host        string
	Port        int
	Name        string
	User        string
	Password    string
	SSLMode     string
	AutoMigrate bool
	AutoSeed    bool
}

const DefaultJWTSecret = "order-express-secret-change-in-production"
const DefaultDBPassword = "order_express_password"

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type RateLimitConfig struct {
	Enabled           bool
	LoginPerMin       int
	OrderCreatePerMin int
	PaymentPerMin     int
}

type Config struct {
	AppEnv     string
	ListenAddr string

	JWTSecret      string
	JWTExpiryHours int

	CORSAllowedOrigins   []string
	CORSAllowCredentials bool

	DB              DBConfig
	Redis           RedisConfig
	CacheTTLSeconds int
	RateLimit       RateLimitConfig
}

func Load() *Config {
	appEnv := getEnv("APP_ENV", "local")
	defaultCORSOrigins := []string{}
	if appEnv == "local" {
		defaultCORSOrigins = []string{"http://localhost:5173"}
	}

	return &Config{
		AppEnv:               appEnv,
		ListenAddr:           getEnv("LISTEN_ADDR", ":3000"),
		JWTSecret:            getEnv("JWT_SECRET", DefaultJWTSecret),
		JWTExpiryHours:       getEnvInt("JWT_EXPIRY_HOURS", 72),
		CORSAllowedOrigins:   getEnvCSV("CORS_ALLOWED_ORIGINS", defaultCORSOrigins),
		CORSAllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", false),
		DB: DBConfig{
			Host:        getEnv("DB_HOST", "127.0.0.1"),
			Port:        getEnvInt("DB_PORT", 5432),
			Name:        getEnv("DB_NAME", "order_express"),
			User:        getEnv("DB_USER", "order_express"),
			Password:    getEnv("DB_PASSWORD", DefaultDBPassword),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			AutoMigrate: getEnvBool("DB_AUTO_MIGRATE", false),
			AutoSeed:    getEnvBool("DB_AUTO_SEED", false),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		CacheTTLSeconds: getEnvInt("CACHE_TTL_SECONDS", 60),
		RateLimit: RateLimitConfig{
			Enabled:           getEnvBool("RL_ENABLED", false),
			LoginPerMin:       getEnvInt("RL_LOGIN_PER_MIN", 20),
			OrderCreatePerMin: getEnvInt("RL_ORDER_CREATE_PER_MIN", 30),
			PaymentPerMin:     getEnvInt("RL_PAYMENT_PER_MIN", 30),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvCSV(key string, fallback []string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		s := strings.TrimSpace(part)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func getEnvInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}
