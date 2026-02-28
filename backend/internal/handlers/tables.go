package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"order-express/backend/internal/idgen"
	"order-express/backend/internal/models"

	"gorm.io/gorm"
)

type TableHandler struct {
	db *gorm.DB
}

func NewTableHandler(db *gorm.DB) *TableHandler {
	return &TableHandler{db: db}
}

type tableRow struct {
	ID             string  `gorm:"column:id"`
	Number         string  `gorm:"column:number"`
	Seats          int     `gorm:"column:seats"`
	Status         string  `gorm:"column:status"`
	CurrentOrderID *string `gorm:"column:current_order_id"`
	QRCodeURL      *string `gorm:"column:qr_code_url"`
	Area           *string `gorm:"column:area"`
	CreatedAt      string  `gorm:"column:created_at"`
	UpdatedAt      string  `gorm:"column:updated_at"`
}

func (h *TableHandler) toTable(r tableRow) models.Table {
	t := models.Table{
		ID:        r.ID,
		Number:    r.Number,
		Seats:     r.Seats,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
	if r.CurrentOrderID != nil && strings.TrimSpace(*r.CurrentOrderID) != "" {
		t.CurrentOrderID = r.CurrentOrderID
	}
	if r.Area != nil {
		t.Area = *r.Area
	}
	if r.QRCodeURL != nil {
		t.QRCodeURL = *r.QRCodeURL
	}
	return t
}

const tableSelectCols = `id, number, seats, status, current_order_id, qr_code_url, area, created_at, updated_at`

func (h *TableHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	area := strings.TrimSpace(q.Get("area"))
	status := strings.TrimSpace(q.Get("status"))

	query := "SELECT " + tableSelectCols + " FROM tables WHERE 1=1"
	args := []any{}
	if area != "" {
		query += " AND area = ?"
		args = append(args, area)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	query += " ORDER BY number"

	rows := []tableRow{}
	tx := h.db.Raw(query, args...).Scan(&rows)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	tables := make([]models.Table, 0, len(rows))
	for _, row := range rows {
		tables = append(tables, h.toTable(row))
	}
	OK(w, tables)
}

func (h *TableHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	area := strings.TrimSpace(q.Get("area"))
	status := strings.TrimSpace(q.Get("status"))

	query := "SELECT id, number, seats, status, area FROM tables WHERE 1=1"
	args := []any{}
	if area != "" {
		query += " AND area = ?"
		args = append(args, area)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	query += " ORDER BY number"

	rows := []tableRow{}
	tx := h.db.Raw(query, args...).Scan(&rows)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	tables := make([]models.Table, 0, len(rows))
	for _, row := range rows {
		tables = append(tables, h.toTable(row))
	}
	OK(w, tables)
}

func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body models.Table
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	body.Number = strings.TrimSpace(body.Number)
	if body.Number == "" {
		BadRequest(w, "number is required")
		return
	}
	if body.Seats <= 0 {
		body.Seats = 4
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id, err := idgen.New("table-", 8)
	if err != nil {
		Fail(w, http.StatusInternalServerError, "internal error")
		return
	}
	qrURL := "/qr/" + id
	status := strings.TrimSpace(body.Status)
	if status == "" {
		status = "available"
	}

	tx := h.db.Exec(
		`INSERT INTO tables (id, number, seats, status, qr_code_url, area, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)`,
		id, body.Number, body.Seats, status, qrURL, strings.TrimSpace(body.Area), now, now,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	var row tableRow
	tx = h.db.Raw("SELECT "+tableSelectCols+" FROM tables WHERE id=?", id).Scan(&row)
	if tx.Error != nil || tx.RowsAffected == 0 {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	OK(w, h.toTable(row))
}

func (h *TableHandler) Areas(w http.ResponseWriter, r *http.Request) {
	areas := []string{}
	if err := h.db.Raw(`SELECT DISTINCT area FROM tables WHERE area IS NOT NULL AND area != '' ORDER BY area`).Scan(&areas).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	OK(w, areas)
}

func (h *TableHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var row tableRow
	tx := h.db.Raw("SELECT "+tableSelectCols+" FROM tables WHERE id=?", id).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "table not found")
		return
	}
	OK(w, h.toTable(row))
}

func (h *TableHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var existing tableRow
	tx := h.db.Raw("SELECT "+tableSelectCols+" FROM tables WHERE id=?", id).Scan(&existing)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "table not found")
		return
	}

	var patch map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	updated := h.toTable(existing)

	if v, ok := patch["number"]; ok {
		json.Unmarshal(v, &updated.Number)
		updated.Number = strings.TrimSpace(updated.Number)
	}
	if v, ok := patch["seats"]; ok {
		json.Unmarshal(v, &updated.Seats)
	}
	if v, ok := patch["status"]; ok {
		json.Unmarshal(v, &updated.Status)
	}
	if v, ok := patch["area"]; ok {
		json.Unmarshal(v, &updated.Area)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	tx = h.db.Exec(
		`UPDATE tables SET number=?, seats=?, status=?, area=?, updated_at=? WHERE id=?`,
		updated.Number, updated.Seats, updated.Status, updated.Area, now, id,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	var row tableRow
	h.db.Raw("SELECT "+tableSelectCols+" FROM tables WHERE id=?", id).Scan(&row)
	OK(w, h.toTable(row))
}

func (h *TableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tx := h.db.Exec(`DELETE FROM tables WHERE id=?`, id)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "table not found")
		return
	}
	OK(w, nil)
}

func (h *TableHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Status         string  `json:"status"`
		CurrentOrderID *string `json:"currentOrderId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	var currentOrderID any
	if body.CurrentOrderID != nil {
		currentOrderID = strings.TrimSpace(*body.CurrentOrderID)
	}

	tx := h.db.Exec(
		`UPDATE tables SET status=?, current_order_id=?, updated_at=? WHERE id=?`,
		strings.TrimSpace(body.Status), currentOrderID, now, id,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "table not found")
		return
	}

	var row tableRow
	h.db.Raw("SELECT "+tableSelectCols+" FROM tables WHERE id=?", id).Scan(&row)
	OK(w, h.toTable(row))
}
