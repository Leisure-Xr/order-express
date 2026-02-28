package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"

	"order-express/backend/internal/idgen"
	"order-express/backend/internal/models"

	"gorm.io/gorm"
)

type OrderHandler struct {
	db *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

var validTransitions = map[string][]string{
	"pending":   {"confirmed", "cancelled"},
	"confirmed": {"preparing", "cancelled"},
	"preparing": {"ready", "cancelled"},
	"ready":     {"completed", "delivered"},
	"completed": {},
	"delivered": {},
	"cancelled": {},
}

type orderRow struct {
	ID                    string         `gorm:"column:id"`
	OrderNumber           string         `gorm:"column:order_number"`
	Type                  string         `gorm:"column:type"`
	Status                string         `gorm:"column:status"`
	TableID               sql.NullString `gorm:"column:table_id"`
	ItemsJSON             string         `gorm:"column:items"`
	Subtotal              float64        `gorm:"column:subtotal"`
	DeliveryFee           float64        `gorm:"column:delivery_fee"`
	Discount              float64        `gorm:"column:discount"`
	Total                 float64        `gorm:"column:total"`
	Remarks               sql.NullString `gorm:"column:remarks"`
	PaymentMethod         string         `gorm:"column:payment_method"`
	PaymentStatus         string         `gorm:"column:payment_status"`
	PaymentPaidAt         sql.NullString `gorm:"column:payment_paid_at"`
	DeliveryAddress       sql.NullString `gorm:"column:delivery_address"`
	ContactPhone          sql.NullString `gorm:"column:contact_phone"`
	EstimatedDeliveryTime sql.NullString `gorm:"column:estimated_delivery_time"`
	StatusHistoryJSON     string         `gorm:"column:status_history"`
	CreatedAt             string         `gorm:"column:created_at"`
	UpdatedAt             string         `gorm:"column:updated_at"`
}

func (h *OrderHandler) parseOrder(r orderRow) models.Order {
	var o models.Order
	o.ID = r.ID
	o.OrderNumber = r.OrderNumber
	o.Type = r.Type
	o.Status = r.Status
	o.Subtotal = r.Subtotal
	o.DeliveryFee = r.DeliveryFee
	o.Discount = r.Discount
	o.Total = r.Total
	o.Payment.Method = r.PaymentMethod
	o.Payment.Status = r.PaymentStatus
	o.CreatedAt = r.CreatedAt
	o.UpdatedAt = r.UpdatedAt

	if r.TableID.Valid {
		o.TableID = &r.TableID.String
	}
	if r.Remarks.Valid {
		o.Remarks = &r.Remarks.String
	}
	if r.DeliveryAddress.Valid {
		o.DeliveryAddress = &r.DeliveryAddress.String
	}
	if r.ContactPhone.Valid {
		o.ContactPhone = &r.ContactPhone.String
	}
	if r.EstimatedDeliveryTime.Valid {
		o.EstimatedDeliveryTime = &r.EstimatedDeliveryTime.String
	}
	if r.PaymentPaidAt.Valid {
		o.Payment.PaidAt = &r.PaymentPaidAt.String
	}

	_ = json.Unmarshal([]byte(r.ItemsJSON), &o.Items)
	_ = json.Unmarshal([]byte(r.StatusHistoryJSON), &o.StatusHistory)
	if o.Items == nil {
		o.Items = []models.OrderItem{}
	}
	if o.StatusHistory == nil {
		o.StatusHistory = []models.StatusHistoryEntry{}
	}
	return o
}

const orderSelectCols = `id, order_number, type, status, table_id,
	items, subtotal, delivery_fee, discount, total,
	remarks, payment_method, payment_status, payment_paid_at,
	delivery_address, contact_phone, estimated_delivery_time,
	status_history, created_at, updated_at`

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	status := q.Get("status")
	orderType := q.Get("type")
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	query := "SELECT " + orderSelectCols + " FROM orders WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM orders WHERE 1=1"
	args := []any{}

	if status != "" && status != "all" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
	}
	if orderType != "" {
		query += " AND type = ?"
		countQuery += " AND type = ?"
		args = append(args, orderType)
	}

	var total int
	if err := h.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	offset := (page - 1) * pageSize
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows := []orderRow{}
	tx := h.db.Raw(query, args...).Scan(&rows)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	orders := make([]models.Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, h.parseOrder(row))
	}

	OK(w, PaginatedResult[models.Order]{
		Items:    orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *OrderHandler) resolveTable(tableIdOrNumber string) (tableID string, tableNumber string, err error) {
	tableIdOrNumber = strings.TrimSpace(tableIdOrNumber)
	if tableIdOrNumber == "" {
		return "", "", errors.New("tableId is required")
	}

	type tableRow struct {
		ID             string         `gorm:"column:id"`
		Number         string         `gorm:"column:number"`
		Status         string         `gorm:"column:status"`
		CurrentOrderID sql.NullString `gorm:"column:current_order_id"`
	}

	var row tableRow
	tx := h.db.Raw(`SELECT id, number, status, current_order_id FROM tables WHERE id=? LIMIT 1`, tableIdOrNumber).Scan(&row)
	if tx.Error != nil {
		return "", "", tx.Error
	}
	if tx.RowsAffected == 0 {
		tx = h.db.Raw(`SELECT id, number, status, current_order_id FROM tables WHERE number=? LIMIT 1`, tableIdOrNumber).Scan(&row)
		if tx.Error != nil {
			return "", "", tx.Error
		}
		if tx.RowsAffected == 0 {
			return "", "", errors.New("table not found")
		}
	}

	if row.Status != "available" && row.CurrentOrderID.Valid && row.CurrentOrderID.String != "" {
		return "", "", errors.New("table is occupied")
	}

	return row.ID, row.Number, nil
}

func isValidOrderType(orderType string) bool {
	switch orderType {
	case "dine_in", "takeaway", "pickup":
		return true
	default:
		return false
	}
}

func isValidPaymentMethod(method string) bool {
	switch method {
	case "wechat", "alipay", "cash":
		return true
	default:
		return false
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateOrderPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	if !isValidOrderType(payload.Type) {
		BadRequest(w, "invalid order type")
		return
	}
	if len(payload.Items) == 0 {
		BadRequest(w, "order must have at least one item")
		return
	}
	for _, it := range payload.Items {
		if it.Quantity < 1 {
			BadRequest(w, "quantity must be >= 1")
			return
		}
	}

	switch payload.Type {
	case "dine_in":
		if strings.TrimSpace(payload.TableID) == "" {
			BadRequest(w, "tableId is required for dine_in orders")
			return
		}
	case "takeaway":
		if strings.TrimSpace(payload.DeliveryAddress) == "" {
			BadRequest(w, "deliveryAddress is required for takeaway orders")
			return
		}
		if strings.TrimSpace(payload.ContactPhone) == "" {
			BadRequest(w, "contactPhone is required for takeaway orders")
			return
		}
	case "pickup":
		if strings.TrimSpace(payload.ContactPhone) == "" {
			BadRequest(w, "contactPhone is required for pickup orders")
			return
		}
	}

	paymentMethod := strings.TrimSpace(payload.PaymentMethod)
	if paymentMethod == "" {
		paymentMethod = "cash"
	}
	if !isValidPaymentMethod(paymentMethod) {
		BadRequest(w, "invalid paymentMethod")
		return
	}

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)

	orderID, err := idgen.New("order-", 12)
	if err != nil {
		Fail(w, http.StatusInternalServerError, "internal error")
		return
	}
	orderNumber := fmt.Sprintf("OE%s-%d-%04d", now.Format("20060102"), now.UnixMilli(), rand.IntN(10000))

	// Resolve items
	items := make([]models.OrderItem, 0, len(payload.Items))
	for _, item := range payload.Items {
		var dish struct {
			Price    float64 `gorm:"column:price"`
			NameJSON string  `gorm:"column:name"`
		}
		tx := h.db.Raw(`SELECT price, name FROM dishes WHERE id=?`, item.DishID).Scan(&dish)
		if tx.Error != nil {
			Fail(w, http.StatusInternalServerError, "database error")
			return
		}
		if tx.RowsAffected == 0 {
			BadRequest(w, "dish not found: "+item.DishID)
			return
		}

		var dishName models.I18nString
		_ = json.Unmarshal([]byte(dish.NameJSON), &dishName)

		optionsTotal := 0.0
		for _, opt := range item.SelectedOptions {
			optionsTotal += opt.PriceAdjustment
		}
		unitPrice := dish.Price + optionsTotal
		subtotal := unitPrice * float64(item.Quantity)

		selectedOpts := item.SelectedOptions
		if selectedOpts == nil {
			selectedOpts = []models.SelectedOption{}
		}

		items = append(items, models.OrderItem{
			DishID:          item.DishID,
			DishName:        dishName,
			Quantity:        item.Quantity,
			UnitPrice:       unitPrice,
			SelectedOptions: selectedOpts,
			Subtotal:        subtotal,
		})
	}

	subtotal := 0.0
	for _, item := range items {
		subtotal += item.Subtotal
	}

	deliveryFee := 0.0
	if payload.Type == "takeaway" {
		deliveryFee = 5.0
	}
	total := subtotal + deliveryFee

	statusHistory := []models.StatusHistoryEntry{{
		Status:    "pending",
		Timestamp: nowStr,
		Note:      "Order placed",
	}}

	itemsJSON, _ := json.Marshal(items)
	historyJSON, _ := json.Marshal(statusHistory)

	var tableDBID string
	var tableNumber string
	if payload.Type == "dine_in" {
		var err error
		tableDBID, tableNumber, err = h.resolveTable(payload.TableID)
		if err != nil {
			Fail(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var tableIDArg any
		if payload.Type == "dine_in" {
			tableIDArg = tableNumber
		}

		var remarksArg any
		if strings.TrimSpace(payload.Remarks) != "" {
			remarksArg = strings.TrimSpace(payload.Remarks)
		}

		var deliveryAddrArg any
		if payload.Type == "takeaway" && strings.TrimSpace(payload.DeliveryAddress) != "" {
			deliveryAddrArg = strings.TrimSpace(payload.DeliveryAddress)
		}

		var contactPhoneArg any
		if strings.TrimSpace(payload.ContactPhone) != "" {
			contactPhoneArg = strings.TrimSpace(payload.ContactPhone)
		}

		if err := tx.Exec(
			`INSERT INTO orders (id, order_number, type, status, table_id, items, subtotal, delivery_fee, discount, total, remarks, payment_method, payment_status, delivery_address, contact_phone, status_history, created_at, updated_at)
	         VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			orderID, orderNumber, payload.Type, "pending", tableIDArg,
			string(itemsJSON), subtotal, deliveryFee, 0, total,
			remarksArg, paymentMethod, "unpaid",
			deliveryAddrArg, contactPhoneArg,
			string(historyJSON), nowStr, nowStr,
		).Error; err != nil {
			return err
		}

		if payload.Type == "dine_in" {
			if err := tx.Exec(
				`UPDATE tables SET status='occupied', current_order_id=?, updated_at=? WHERE id=?`,
				orderID, nowStr, tableDBID,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	var row orderRow
	tx := h.db.Raw("SELECT "+orderSelectCols+" FROM orders WHERE id=?", orderID).Scan(&row)
	if tx.Error != nil || tx.RowsAffected == 0 {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	OK(w, h.parseOrder(row))
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var row orderRow
	tx := h.db.Raw("SELECT "+orderSelectCols+" FROM orders WHERE id=?", id).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "order not found")
		return
	}
	OK(w, h.parseOrder(row))
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	type statusRow struct {
		Status            string         `gorm:"column:status"`
		TableID           sql.NullString `gorm:"column:table_id"`
		StatusHistoryJSON string         `gorm:"column:status_history"`
	}
	var current statusRow
	tx := h.db.Raw(`SELECT status, table_id, status_history FROM orders WHERE id=?`, id).Scan(&current)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "order not found")
		return
	}

	allowed := validTransitions[current.Status]
	valid := false
	for _, s := range allowed {
		if s == body.Status {
			valid = true
			break
		}
	}
	if !valid {
		BadRequest(w, fmt.Sprintf("cannot transition from %q to %q", current.Status, body.Status))
		return
	}

	var history []models.StatusHistoryEntry
	_ = json.Unmarshal([]byte(current.StatusHistoryJSON), &history)
	nowStr := time.Now().UTC().Format(time.RFC3339)
	history = append(history, models.StatusHistoryEntry{
		Status:    body.Status,
		Timestamp: nowStr,
		Note:      body.Note,
	})
	updatedHistoryJSON, _ := json.Marshal(history)

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			`UPDATE orders SET status=?, status_history=?, updated_at=? WHERE id=?`,
			body.Status, string(updatedHistoryJSON), nowStr, id,
		).Error; err != nil {
			return err
		}

		if body.Status == "cancelled" && current.TableID.Valid && strings.TrimSpace(current.TableID.String) != "" {
			if err := tx.Exec(
				`UPDATE tables SET status='available', current_order_id=NULL, updated_at=? WHERE current_order_id=?`,
				nowStr, id,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	var row orderRow
	h.db.Raw("SELECT "+orderSelectCols+" FROM orders WHERE id=?", id).Scan(&row)
	OK(w, h.parseOrder(row))
}

func (h *OrderHandler) History(w http.ResponseWriter, r *http.Request) {
	rows := []orderRow{}
	tx := h.db.Raw(
		"SELECT " + orderSelectCols + " FROM orders WHERE status IN ('completed','delivered','cancelled') ORDER BY updated_at DESC",
	).Scan(&rows)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	orders := make([]models.Order, 0, len(rows))
	for _, row := range rows {
		orders = append(orders, h.parseOrder(row))
	}
	OK(w, orders)
}
