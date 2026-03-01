package router

import (
	"net/http"
	"time"

	"order-express/backend/internal/cache"
	"order-express/backend/internal/config"
	"order-express/backend/internal/handlers"
	"order-express/backend/internal/middleware"

	"gorm.io/gorm"
)

func New(db *gorm.DB, cfg *config.Config, c *cache.Cache, rl *middleware.RateLimiter) http.Handler {
	mux := http.NewServeMux()

	auth := middleware.Auth(cfg)
	requireAdmin := middleware.RequireAdmin()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	authH := handlers.NewAuthHandler(db, cfg)
	cacheTTL := time.Duration(cfg.CacheTTLSeconds) * time.Second
	catH := handlers.NewCategoryHandler(db, c, cacheTTL)
	dishH := handlers.NewDishHandler(db, c, cacheTTL)
	orderH := handlers.NewOrderHandler(db)
	payH := handlers.NewPaymentHandler(db)
	tableH := handlers.NewTableHandler(db)
	storeH := handlers.NewStoreHandler(db, c, cacheTTL)

	var loginHandler http.Handler = http.HandlerFunc(authH.Login)
	var orderCreateHandler http.Handler = http.HandlerFunc(orderH.Create)
	var paymentInitHandler http.Handler = http.HandlerFunc(payH.Initiate)
	var orderGetHandler http.Handler = http.HandlerFunc(orderH.GetByID)
	var paymentGetHandler http.Handler = http.HandlerFunc(payH.GetByID)
	if cfg.RateLimit.Enabled && rl != nil {
		loginHandler = rl.Middleware("login", cfg.RateLimit.LoginPerMin, time.Minute)(loginHandler)
		orderCreateHandler = rl.Middleware("order_create", cfg.RateLimit.OrderCreatePerMin, time.Minute)(orderCreateHandler)
		paymentInitHandler = rl.Middleware("payment_initiate", cfg.RateLimit.PaymentPerMin, time.Minute)(paymentInitHandler)
		orderGetHandler = rl.Middleware("order_get", cfg.RateLimit.PublicReadPerMin, time.Minute)(orderGetHandler)
		paymentGetHandler = rl.Middleware("payment_get", cfg.RateLimit.PublicReadPerMin, time.Minute)(paymentGetHandler)
	}

	// Auth
	mux.Handle("POST /api/auth/login", loginHandler)
	mux.HandleFunc("POST /api/auth/logout", authH.Logout)
	mux.Handle("GET /api/auth/me", auth(http.HandlerFunc(authH.Me)))

	// Categories
	mux.HandleFunc("GET /api/categories", catH.List)
	mux.Handle("POST /api/categories/reorder", auth(requireAdmin(http.HandlerFunc(catH.Reorder))))
	mux.Handle("POST /api/categories", auth(requireAdmin(http.HandlerFunc(catH.Create))))
	mux.HandleFunc("GET /api/categories/{id}", catH.GetByID)
	mux.Handle("PATCH /api/categories/{id}", auth(requireAdmin(http.HandlerFunc(catH.Update))))
	mux.Handle("DELETE /api/categories/{id}", auth(requireAdmin(http.HandlerFunc(catH.Delete))))

	// Dishes
	mux.HandleFunc("GET /api/dishes", dishH.List)
	mux.Handle("POST /api/dishes", auth(requireAdmin(http.HandlerFunc(dishH.Create))))
	mux.HandleFunc("GET /api/dishes/{id}", dishH.GetByID)
	mux.Handle("PATCH /api/dishes/{id}", auth(requireAdmin(http.HandlerFunc(dishH.Update))))
	mux.Handle("DELETE /api/dishes/{id}", auth(requireAdmin(http.HandlerFunc(dishH.Delete))))
	mux.Handle("POST /api/dishes/{id}/toggle-status", auth(requireAdmin(http.HandlerFunc(dishH.ToggleStatus))))

	// Orders — note: fixed paths must come before wildcard paths
	mux.Handle("GET /api/orders/history", auth(requireAdmin(http.HandlerFunc(orderH.History))))
	mux.Handle("GET /api/orders/stats", auth(requireAdmin(http.HandlerFunc(orderH.Stats))))
	mux.Handle("GET /api/orders", auth(requireAdmin(http.HandlerFunc(orderH.List))))
	mux.Handle("POST /api/orders", orderCreateHandler)
	mux.Handle("GET /api/orders/{id}", orderGetHandler)
	mux.Handle("PATCH /api/orders/{id}/status", auth(requireAdmin(http.HandlerFunc(orderH.UpdateStatus))))
	mux.Handle("GET /api/orders/{orderId}/payments", auth(requireAdmin(http.HandlerFunc(payH.ByOrderID))))

	// Payments
	mux.Handle("POST /api/payments/initiate", paymentInitHandler)
	mux.Handle("GET /api/payments/{paymentId}", paymentGetHandler)
	mux.Handle("POST /api/payments/{paymentId}/refund", auth(requireAdmin(http.HandlerFunc(payH.Refund))))

	// Public tables (customer table selection)
	mux.HandleFunc("GET /api/public/tables", tableH.ListPublic)

	// Tables (admin) — fixed paths before wildcard
	mux.Handle("GET /api/tables/areas", auth(requireAdmin(http.HandlerFunc(tableH.Areas))))
	mux.Handle("GET /api/tables", auth(requireAdmin(http.HandlerFunc(tableH.List))))
	mux.Handle("POST /api/tables", auth(requireAdmin(http.HandlerFunc(tableH.Create))))
	mux.Handle("GET /api/tables/{id}", auth(requireAdmin(http.HandlerFunc(tableH.GetByID))))
	mux.Handle("PATCH /api/tables/{id}", auth(requireAdmin(http.HandlerFunc(tableH.Update))))
	mux.Handle("DELETE /api/tables/{id}", auth(requireAdmin(http.HandlerFunc(tableH.Delete))))
	mux.Handle("PATCH /api/tables/{id}/status", auth(requireAdmin(http.HandlerFunc(tableH.UpdateStatus))))

	// Store
	mux.HandleFunc("GET /api/store/business-hours", storeH.GetBusinessHours)
	mux.Handle("PUT /api/store/business-hours", auth(requireAdmin(http.HandlerFunc(storeH.UpdateBusinessHours))))
	mux.HandleFunc("GET /api/store/delivery-settings", storeH.GetDeliverySettings)
	mux.Handle("PATCH /api/store/delivery-settings", auth(requireAdmin(http.HandlerFunc(storeH.UpdateDeliverySettings))))
	mux.HandleFunc("GET /api/store", storeH.GetInfo)
	mux.Handle("PATCH /api/store", auth(requireAdmin(http.HandlerFunc(storeH.UpdateInfo))))

	return middleware.MaxBodySize(1 << 20)(middleware.CORS(cfg)(mux)) // 1 MiB body limit
}
