package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	if pageSize > 100 {
		pageSize = 100
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
	if len(payload.Items) > 50 {
		BadRequest(w, "too many items in order")
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
		if len(payload.DeliveryAddress) > 500 {
			BadRequest(w, "deliveryAddress too long")
			return
		}
		if strings.TrimSpace(payload.ContactPhone) == "" {
			BadRequest(w, "contactPhone is required for takeaway orders")
			return
		}
		if len(payload.ContactPhone) > 30 {
			BadRequest(w, "contactPhone too long")
			return
		}
	case "pickup":
		if strings.TrimSpace(payload.ContactPhone) == "" {
			BadRequest(w, "contactPhone is required for pickup orders")
			return
		}
		if len(payload.ContactPhone) > 30 {
			BadRequest(w, "contactPhone too long")
			return
		}
	}

	if len(payload.Remarks) > 500 {
		BadRequest(w, "remarks too long")
		return
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
			Price       float64 `gorm:"column:price"`
			NameJSON    string  `gorm:"column:name"`
			OptionsJSON string  `gorm:"column:options"`
		}
		tx := h.db.Raw(`SELECT price, name, options FROM dishes WHERE id=?`, item.DishID).Scan(&dish)
		if tx.Error != nil {
			Fail(w, http.StatusInternalServerError, "database error")
			return
		}
		if tx.RowsAffected == 0 {
			BadRequest(w, "dish not found")
			return
		}

		var dishName models.I18nString
		_ = json.Unmarshal([]byte(dish.NameJSON), &dishName)

		// Parse dish options from database for price validation
		var dishOptions []models.DishOption
		_ = json.Unmarshal([]byte(dish.OptionsJSON), &dishOptions)

		// Build lookup map: various key combos -> serverPriceAdjustment
		optionPriceLookup := make(map[string]float64)
		for _, opt := range dishOptions {
			for _, val := range opt.Values {
				optionPriceLookup[opt.ID+":"+val.ID] = val.PriceAdjustment
				optionPriceLookup[opt.Name.ZH+":"+val.Label.ZH] = val.PriceAdjustment
				if opt.Name.EN != "" && val.Label.EN != "" {
					optionPriceLookup[opt.Name.EN+":"+val.Label.EN] = val.PriceAdjustment
				}
			}
		}

		optionsTotal := 0.0
		validatedOpts := make([]models.SelectedOption, 0, len(item.SelectedOptions))
		for _, opt := range item.SelectedOptions {
			// Look up the real price adjustment from database
			serverPrice, found := optionPriceLookup[opt.OptionName+":"+opt.ValueName]
			if !found {
				// Option not found in dish definition — skip or reject
				BadRequest(w, "invalid dish option")
				return
			}
			optionsTotal += serverPrice
			validatedOpts = append(validatedOpts, models.SelectedOption{
				OptionName:      opt.OptionName,
				ValueName:       opt.ValueName,
				PriceAdjustment: serverPrice,
			})
		}
		unitPrice := dish.Price + optionsTotal
		subtotal := unitPrice * float64(item.Quantity)

		selectedOpts := validatedOpts
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
			log.Printf("resolveTable error: %v", err)
			BadRequest(w, "invalid table")
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
		BadRequest(w, "invalid status transition")
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
		"SELECT " + orderSelectCols + " FROM orders WHERE status IN ('completed','delivered','cancelled') ORDER BY updated_at DESC LIMIT 200",
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

func (h *OrderHandler) Stats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fromStr := q.Get("from")
	toStr := q.Get("to")

	now := time.Now().UTC()
	if toStr == "" {
		toStr = now.Format("2006-01-02")
	}
	if fromStr == "" {
		fromStr = now.AddDate(0, 0, -6).Format("2006-01-02")
	}

	fromDate, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		BadRequest(w, "invalid 'from' date format, expected YYYY-MM-DD")
		return
	}
	toDate, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		BadRequest(w, "invalid 'to' date format, expected YYYY-MM-DD")
		return
	}
	if fromDate.After(toDate) {
		BadRequest(w, "'from' must not be after 'to'")
		return
	}

	// Inclusive range: from <= created_at < to+1day
	fromRFC := fromDate.Format(time.RFC3339)
	toRFC := toDate.AddDate(0, 0, 1).Format(time.RFC3339)

	// 1. Total orders in range
	var totalOrders int
	if err := h.db.Raw(`SELECT COUNT(*) FROM orders WHERE created_at >= ? AND created_at < ?`,
		fromRFC, toRFC).Scan(&totalOrders).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	// 2. Paid revenue
	var paidRevenue float64
	if err := h.db.Raw(`SELECT COALESCE(SUM(total), 0) FROM orders WHERE payment_status = 'paid' AND created_at >= ? AND created_at < ?`,
		fromRFC, toRFC).Scan(&paidRevenue).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	// 3. Pending orders (current state, not range-dependent)
	var pendingOrders int
	if err := h.db.Raw(`SELECT COUNT(*) FROM orders WHERE status = 'pending'`).Scan(&pendingOrders).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	// 4. Daily revenue
	type dailyRow struct {
		Date   string  `gorm:"column:date"`
		Amount float64 `gorm:"column:amount"`
		Orders int     `gorm:"column:orders"`
	}
	var dailyRows []dailyRow
	if err := h.db.Raw(`SELECT LEFT(created_at, 10) AS date,
		COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN total ELSE 0 END), 0) AS amount,
		COUNT(*) AS orders
		FROM orders
		WHERE created_at >= ? AND created_at < ?
		GROUP BY LEFT(created_at, 10)
		ORDER BY date`, fromRFC, toRFC).Scan(&dailyRows).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	dailyMap := make(map[string]dailyRow)
	for _, dr := range dailyRows {
		dailyMap[dr.Date] = dr
	}
	dailyRevenue := make([]models.DailyRevenue, 0)
	for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		if dr, ok := dailyMap[key]; ok {
			dailyRevenue = append(dailyRevenue, models.DailyRevenue{Date: key, Amount: dr.Amount, Orders: dr.Orders})
		} else {
			dailyRevenue = append(dailyRevenue, models.DailyRevenue{Date: key, Amount: 0, Orders: 0})
		}
	}

	// 5. Popular items
	type popRow struct {
		DishID   string `gorm:"column:dish_id"`
		DishName string `gorm:"column:dish_name"`
		Count    int    `gorm:"column:count"`
	}
	var popRows []popRow
	if err := h.db.Raw(`SELECT elem->>'dishId' AS dish_id,
		elem->>'dishName' AS dish_name,
		SUM((elem->>'quantity')::int) AS count
		FROM orders,
		jsonb_array_elements(items::jsonb) AS elem
		WHERE created_at >= ? AND created_at < ?
		GROUP BY dish_id, dish_name
		ORDER BY count DESC
		LIMIT 10`, fromRFC, toRFC).Scan(&popRows).Error; err != nil {
		log.Printf("popular items query error: %v", err)
		popRows = nil
	}

	popularItems := make([]models.PopularItem, 0, len(popRows))
	for _, pr := range popRows {
		var name models.I18nString
		_ = json.Unmarshal([]byte(pr.DishName), &name)
		popularItems = append(popularItems, models.PopularItem{DishID: pr.DishID, DishName: name, Count: pr.Count})
	}

	// 6. Orders by type
	type typeRow struct {
		Type  string `gorm:"column:type"`
		Count int    `gorm:"column:count"`
	}
	var typeRows []typeRow
	if err := h.db.Raw(`SELECT type, COUNT(*) AS count FROM orders
		WHERE created_at >= ? AND created_at < ?
		GROUP BY type`, fromRFC, toRFC).Scan(&typeRows).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	ordersByType := make(map[string]int)
	for _, tr := range typeRows {
		ordersByType[tr.Type] = tr.Count
	}

	OK(w, models.OrderStats{
		TotalOrders:   totalOrders,
		PaidRevenue:   paidRevenue,
		PendingOrders: pendingOrders,
		DailyRevenue:  dailyRevenue,
		PopularItems:  popularItems,
		OrdersByType:  ordersByType,
	})
}
