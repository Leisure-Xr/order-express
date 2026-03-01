package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"order-express/backend/internal/idgen"
	"order-express/backend/internal/models"

	"gorm.io/gorm"
)

type PaymentHandler struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

type paymentRow struct {
	PaymentID string         `gorm:"column:payment_id"`
	OrderID   string         `gorm:"column:order_id"`
	Method    string         `gorm:"column:method"`
	Status    string         `gorm:"column:status"`
	Amount    float64        `gorm:"column:amount"`
	CreatedAt string         `gorm:"column:created_at"`
	PaidAt    sql.NullString `gorm:"column:paid_at"`
}

func (h *PaymentHandler) toPayment(r paymentRow) models.Payment {
	p := models.Payment{
		PaymentID: r.PaymentID,
		OrderID:   r.OrderID,
		Method:    r.Method,
		Status:    r.Status,
		Amount:    r.Amount,
		CreatedAt: r.CreatedAt,
	}
	if r.PaidAt.Valid {
		p.PaidAt = &r.PaidAt.String
	}
	return p
}

func (h *PaymentHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		OrderID string  `json:"orderId"`
		Method  string  `json:"method"`
		Amount  float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	body.OrderID = strings.TrimSpace(body.OrderID)
	body.Method = strings.TrimSpace(body.Method)
	if body.OrderID == "" {
		BadRequest(w, "orderId is required")
		return
	}
	if body.Amount <= 0 {
		BadRequest(w, "amount must be > 0")
		return
	}
	if body.Method == "" {
		body.Method = "cash"
	}

	now := time.Now().UTC().Format(time.RFC3339)
	paymentID, err := idgen.New("pay-", 12)
	if err != nil {
		Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	var payment paymentRow

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var order struct {
			ID    string  `gorm:"column:id"`
			Total float64 `gorm:"column:total"`
		}
		orderTx := tx.Raw(`SELECT id, total FROM orders WHERE id=?`, body.OrderID).Scan(&order)
		if orderTx.Error != nil {
			return orderTx.Error
		}
		if orderTx.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if order.Total != body.Amount {
			return fmt.Errorf("amount mismatch: expected %.2f, got %.2f", order.Total, body.Amount)
		}

		// Check for existing paid payment (idempotency)
		var existingCount int64
		if err := tx.Raw(`SELECT COUNT(*) FROM payments WHERE order_id=? AND status='paid'`, body.OrderID).Scan(&existingCount).Error; err != nil {
			return err
		}
		if existingCount > 0 {
			return fmt.Errorf("duplicate payment")
		}

		// Simulate immediate payment success
		if err := tx.Exec(
			`INSERT INTO payments (payment_id, order_id, method, status, amount, created_at, paid_at) VALUES (?,?,?,?,?,?,?)`,
			paymentID, body.OrderID, body.Method, "paid", body.Amount, now, now,
		).Error; err != nil {
			return err
		}

		if err := tx.Exec(
			`UPDATE orders SET payment_method=?, payment_status='paid', payment_paid_at=?, updated_at=? WHERE id=?`,
			body.Method, now, now, body.OrderID,
		).Error; err != nil {
			return err
		}

		getTx := tx.Raw(`SELECT payment_id, order_id, method, status, amount, created_at, paid_at FROM payments WHERE payment_id=?`, paymentID).Scan(&payment)
		if getTx.Error != nil {
			return getTx.Error
		}
		if getTx.RowsAffected == 0 {
			return fmt.Errorf("payment not found after insert")
		}

		return nil
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(w, "order not found")
			return
		}
		if strings.Contains(err.Error(), "amount mismatch") {
			log.Printf("payment amount mismatch: %v", err)
			BadRequest(w, "payment amount does not match order total")
			return
		}
		if strings.Contains(err.Error(), "duplicate payment") {
			BadRequest(w, "order already paid")
			return
		}
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	OK(w, h.toPayment(payment))
}

func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentId")
	var row paymentRow
	tx := h.db.Raw(`SELECT payment_id, order_id, method, status, amount, created_at, paid_at FROM payments WHERE payment_id=?`, paymentID).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "payment not found")
		return
	}
	OK(w, h.toPayment(row))
}

func (h *PaymentHandler) ByOrderID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderId")
	rows := []paymentRow{}
	tx := h.db.Raw(`SELECT payment_id, order_id, method, status, amount, created_at, paid_at FROM payments WHERE order_id=? ORDER BY created_at`, orderID).Scan(&rows)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	payments := make([]models.Payment, 0, len(rows))
	for _, row := range rows {
		payments = append(payments, h.toPayment(row))
	}
	OK(w, payments)
}

func (h *PaymentHandler) Refund(w http.ResponseWriter, r *http.Request) {
	paymentID := r.PathValue("paymentId")
	now := time.Now().UTC().Format(time.RFC3339)

	var payment paymentRow

	err := h.db.Transaction(func(tx *gorm.DB) error {
		upd := tx.Exec(`UPDATE payments SET status='refunded', paid_at=? WHERE payment_id=?`, now, paymentID)
		if upd.Error != nil {
			return upd.Error
		}
		if upd.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Exec(
			`UPDATE orders SET payment_status='refunded', updated_at=? WHERE id=(SELECT order_id FROM payments WHERE payment_id=?)`,
			now, paymentID,
		).Error; err != nil {
			return err
		}

		getTx := tx.Raw(`SELECT payment_id, order_id, method, status, amount, created_at, paid_at FROM payments WHERE payment_id=?`, paymentID).Scan(&payment)
		if getTx.Error != nil {
			return getTx.Error
		}
		if getTx.RowsAffected == 0 {
			return fmt.Errorf("payment not found after update")
		}

		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(w, "payment not found")
			return
		}
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	OK(w, h.toPayment(payment))
}
