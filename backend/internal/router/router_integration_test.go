package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"order-express/backend/internal/config"
	"order-express/backend/internal/db"
	"order-express/backend/internal/models"
	"order-express/backend/internal/router"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type apiResponse[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_DSN to run integration tests")
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := db.Seed(gdb); err != nil {
		t.Fatalf("seed: %v", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_ = gdb.Exec(`DELETE FROM payments`).Error
	_ = gdb.Exec(`DELETE FROM orders`).Error
	_ = gdb.Exec(`UPDATE tables SET status='available', current_order_id=NULL, updated_at=?`, now).Error

	return gdb
}

func postJSON(t *testing.T, baseURL string, path string, body any, out any) int {
	t.Helper()

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	res, err := http.Post(baseURL+path, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			t.Fatalf("decode: %v", err)
		}
	}
	return res.StatusCode
}

func patchJSON(t *testing.T, baseURL string, path string, body any, out any) int {
	t.Helper()

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, baseURL+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	defer res.Body.Close()

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			t.Fatalf("decode: %v", err)
		}
	}
	return res.StatusCode
}

func TestOrderFlow_DineInTableNumber_OccupyAndRelease(t *testing.T) {
	gdb := openTestDB(t)

	cfg := &config.Config{
		ListenAddr:      ":0",
		JWTSecret:       "test",
		JWTExpiryHours:  72,
		CacheTTLSeconds: 1,
	}

	srv := httptest.NewServer(router.New(gdb, cfg, nil, nil))
	defer srv.Close()

	var created apiResponse[*models.Order]
	status := postJSON(t, srv.URL, "/api/orders", map[string]any{
		"type":    "dine_in",
		"tableId": "T05",
		"items": []map[string]any{
			{
				"dishId":          "dish-001",
				"quantity":        1,
				"selectedOptions": []any{},
			},
		},
		"paymentMethod": "cash",
	}, &created)
	if status != http.StatusOK || created.Code != 200 || created.Data == nil {
		t.Fatalf("create order: status=%d code=%d msg=%q", status, created.Code, created.Message)
	}

	var table struct {
		Status string  `gorm:"column:status"`
		OID    *string `gorm:"column:current_order_id"`
	}
	if err := gdb.Raw(`SELECT status, current_order_id FROM tables WHERE id='table-05'`).Scan(&table).Error; err != nil {
		t.Fatalf("query table: %v", err)
	}
	if table.Status != "occupied" || table.OID == nil || *table.OID != created.Data.ID {
		t.Fatalf("expected table occupied by %s, got status=%s oid=%v", created.Data.ID, table.Status, table.OID)
	}

	// Invalid transition: pending -> ready
	var invalid apiResponse[any]
	status = patchJSON(t, srv.URL, "/api/orders/"+created.Data.ID+"/status", map[string]any{
		"status": "ready",
		"note":   "skip",
	}, &invalid)
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid transition, got %d", status)
	}

	// Cancel order: should release table
	var cancelled apiResponse[*models.Order]
	status = patchJSON(t, srv.URL, "/api/orders/"+created.Data.ID+"/status", map[string]any{
		"status": "cancelled",
		"note":   "test cancel",
	}, &cancelled)
	if status != http.StatusOK || cancelled.Code != 200 || cancelled.Data == nil {
		t.Fatalf("cancel order: status=%d code=%d msg=%q", status, cancelled.Code, cancelled.Message)
	}

	var tableAfter struct {
		Status string  `gorm:"column:status"`
		OID    *string `gorm:"column:current_order_id"`
	}
	if err := gdb.Raw(`SELECT status, current_order_id FROM tables WHERE id='table-05'`).Scan(&tableAfter).Error; err != nil {
		t.Fatalf("query table after: %v", err)
	}
	if tableAfter.Status != "available" || tableAfter.OID != nil {
		t.Fatalf("expected table available and no order id, got status=%s oid=%v", tableAfter.Status, tableAfter.OID)
	}
}

func TestPaymentFlow_AmountMismatchAndPaid(t *testing.T) {
	gdb := openTestDB(t)

	cfg := &config.Config{
		ListenAddr:      ":0",
		JWTSecret:       "test",
		JWTExpiryHours:  72,
		CacheTTLSeconds: 1,
	}
	srv := httptest.NewServer(router.New(gdb, cfg, nil, nil))
	defer srv.Close()

	var created apiResponse[*models.Order]
	status := postJSON(t, srv.URL, "/api/orders", map[string]any{
		"type": "pickup",
		"items": []map[string]any{
			{
				"dishId":          "dish-001",
				"quantity":        1,
				"selectedOptions": []any{},
			},
		},
		"contactPhone":  "13800138000",
		"paymentMethod": "wechat",
	}, &created)
	if status != http.StatusOK || created.Code != 200 || created.Data == nil {
		t.Fatalf("create order: status=%d code=%d msg=%q", status, created.Code, created.Message)
	}

	// Amount mismatch
	var mismatch apiResponse[any]
	status = postJSON(t, srv.URL, "/api/payments/initiate", map[string]any{
		"orderId": created.Data.ID,
		"method":  "wechat",
		"amount":  created.Data.Total + 1,
	}, &mismatch)
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400 for amount mismatch, got %d", status)
	}

	// Correct payment
	var paid apiResponse[*models.Payment]
	status = postJSON(t, srv.URL, "/api/payments/initiate", map[string]any{
		"orderId": created.Data.ID,
		"method":  "wechat",
		"amount":  created.Data.Total,
	}, &paid)
	if status != http.StatusOK || paid.Code != 200 || paid.Data == nil {
		t.Fatalf("init payment: status=%d code=%d msg=%q", status, paid.Code, paid.Message)
	}
	if paid.Data.Status != "paid" {
		t.Fatalf("expected payment paid, got %s", paid.Data.Status)
	}

	var got apiResponse[*models.Order]
	res, err := http.Get(srv.URL + "/api/orders/" + created.Data.ID)
	if err != nil {
		t.Fatalf("get order: %v", err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode order: %v", err)
	}
	if got.Data == nil || got.Data.Payment.Status != "paid" {
		t.Fatalf("expected order payment status paid, got %+v", got.Data)
	}
}
