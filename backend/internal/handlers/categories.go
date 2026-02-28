package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"order-express/backend/internal/cache"
	"order-express/backend/internal/idgen"
	"order-express/backend/internal/models"

	"gorm.io/gorm"
)

type CategoryHandler struct {
	db       *gorm.DB
	cache    *cache.Cache
	cacheTTL time.Duration
}

func NewCategoryHandler(db *gorm.DB, c *cache.Cache, ttl time.Duration) *CategoryHandler {
	return &CategoryHandler{db: db, cache: c, cacheTTL: ttl}
}

func (h *CategoryHandler) listCategories() ([]models.Category, error) {
	type categoryRow struct {
		ID        string `gorm:"column:id"`
		NameJSON  string `gorm:"column:name"`
		Icon      string `gorm:"column:icon"`
		Image     string `gorm:"column:image"`
		SortOrder int    `gorm:"column:sort_order"`
		Status    string `gorm:"column:status"`
		DishCount int    `gorm:"column:dish_count"`
		CreatedAt string `gorm:"column:created_at"`
		UpdatedAt string `gorm:"column:updated_at"`
	}

	rows := []categoryRow{}
	tx := h.db.Raw(`
		SELECT c.id, c.name, COALESCE(c.icon,'') AS icon, COALESCE(c.image,'') AS image, c.sort_order, c.status,
		       COUNT(d.id) AS dish_count, c.created_at, c.updated_at
		FROM categories c
		LEFT JOIN dishes d ON d.category_id = c.id
		GROUP BY c.id
		ORDER BY c.sort_order
	`).Scan(&rows)
	if tx.Error != nil {
		return nil, tx.Error
	}

	cats := make([]models.Category, 0, len(rows))
	for _, r := range rows {
		var name models.I18nString
		_ = json.Unmarshal([]byte(r.NameJSON), &name)
		cats = append(cats, models.Category{
			ID:        r.ID,
			Name:      name,
			Icon:      r.Icon,
			Image:     r.Image,
			SortOrder: r.SortOrder,
			Status:    r.Status,
			DishCount: r.DishCount,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return cats, nil
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	if h.cache != nil {
		if payload, ok, err := h.cache.GetString(r.Context(), "categories:list"); err == nil && ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
			return
		}
	}

	cats, err := h.listCategories()
	if err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	resp := ApiResponse{Code: 200, Data: cats, Message: "success"}
	b, err := json.Marshal(resp)
	if err != nil {
		OK(w, cats)
		return
	}
	if h.cache != nil {
		_ = h.cache.SetString(r.Context(), "categories:list", string(b), h.cacheTTL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body models.Category
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id, err := idgen.New("cat-", 8)
	if err != nil {
		Fail(w, http.StatusInternalServerError, "internal error")
		return
	}
	nameJSON, _ := json.Marshal(body.Name)
	status := body.Status
	if status == "" {
		status = "active"
	}

	tx := h.db.Exec(
		`INSERT INTO categories (id, name, icon, image, sort_order, status, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)`,
		id, string(nameJSON), body.Icon, body.Image, body.SortOrder, status, now, now,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	type categoryRow struct {
		ID        string `gorm:"column:id"`
		NameJSON  string `gorm:"column:name"`
		Icon      string `gorm:"column:icon"`
		Image     string `gorm:"column:image"`
		SortOrder int    `gorm:"column:sort_order"`
		Status    string `gorm:"column:status"`
		DishCount int    `gorm:"column:count"`
		CreatedAt string `gorm:"column:created_at"`
		UpdatedAt string `gorm:"column:updated_at"`
	}
	var row categoryRow
	tx = h.db.Raw(`SELECT c.id, c.name, COALESCE(c.icon,'') AS icon, COALESCE(c.image,'') AS image, c.sort_order, c.status, COUNT(d.id) AS count, c.created_at, c.updated_at FROM categories c LEFT JOIN dishes d ON d.category_id=c.id WHERE c.id=? GROUP BY c.id`, id).Scan(&row)
	if tx.Error != nil || tx.RowsAffected == 0 {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	var name models.I18nString
	_ = json.Unmarshal([]byte(row.NameJSON), &name)
	cat := models.Category{
		ID:        row.ID,
		Name:      name,
		Icon:      row.Icon,
		Image:     row.Image,
		SortOrder: row.SortOrder,
		Status:    row.Status,
		DishCount: row.DishCount,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "categories:list")
	}
	OK(w, cat)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	type categoryRow struct {
		ID        string `gorm:"column:id"`
		NameJSON  string `gorm:"column:name"`
		Icon      string `gorm:"column:icon"`
		Image     string `gorm:"column:image"`
		SortOrder int    `gorm:"column:sort_order"`
		Status    string `gorm:"column:status"`
		DishCount int    `gorm:"column:count"`
		CreatedAt string `gorm:"column:created_at"`
		UpdatedAt string `gorm:"column:updated_at"`
	}
	var row categoryRow
	tx := h.db.Raw(`SELECT c.id, c.name, COALESCE(c.icon,'') AS icon, COALESCE(c.image,'') AS image, c.sort_order, c.status, COUNT(d.id) AS count, c.created_at, c.updated_at FROM categories c LEFT JOIN dishes d ON d.category_id=c.id WHERE c.id=? GROUP BY c.id`, id).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "category not found")
		return
	}
	var name models.I18nString
	_ = json.Unmarshal([]byte(row.NameJSON), &name)
	cat := models.Category{
		ID:        row.ID,
		Name:      name,
		Icon:      row.Icon,
		Image:     row.Image,
		SortOrder: row.SortOrder,
		Status:    row.Status,
		DishCount: row.DishCount,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	OK(w, cat)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))

	// Load existing
	var existing models.Category
	var nameJSON string
	var row struct {
		ID        string `gorm:"column:id"`
		NameJSON  string `gorm:"column:name"`
		Icon      string `gorm:"column:icon"`
		Image     string `gorm:"column:image"`
		SortOrder int    `gorm:"column:sort_order"`
		Status    string `gorm:"column:status"`
	}
	tx := h.db.Raw(`SELECT id, name, COALESCE(icon,'') AS icon, COALESCE(image,'') AS image, sort_order, status FROM categories WHERE id=?`, id).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "category not found")
		return
	}
	existing.ID = row.ID
	nameJSON = row.NameJSON
	existing.Icon = row.Icon
	existing.Image = row.Image
	existing.SortOrder = row.SortOrder
	existing.Status = row.Status
	_ = json.Unmarshal([]byte(nameJSON), &existing.Name)

	// Decode patch
	var patch map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	// Apply patch
	if v, ok := patch["name"]; ok {
		json.Unmarshal(v, &existing.Name)
	}
	if v, ok := patch["icon"]; ok {
		json.Unmarshal(v, &existing.Icon)
	}
	if v, ok := patch["image"]; ok {
		json.Unmarshal(v, &existing.Image)
	}
	if v, ok := patch["sortOrder"]; ok {
		json.Unmarshal(v, &existing.SortOrder)
	}
	if v, ok := patch["status"]; ok {
		json.Unmarshal(v, &existing.Status)
	}

	updatedNameJSON, _ := json.Marshal(existing.Name)
	now := time.Now().UTC().Format(time.RFC3339)

	tx = h.db.Exec(
		`UPDATE categories SET name=?, icon=?, image=?, sort_order=?, status=?, updated_at=? WHERE id=?`,
		string(updatedNameJSON), existing.Icon, existing.Image, existing.SortOrder, existing.Status, now, id,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	type categoryRow struct {
		ID        string `gorm:"column:id"`
		NameJSON  string `gorm:"column:name"`
		Icon      string `gorm:"column:icon"`
		Image     string `gorm:"column:image"`
		SortOrder int    `gorm:"column:sort_order"`
		Status    string `gorm:"column:status"`
		DishCount int    `gorm:"column:count"`
		CreatedAt string `gorm:"column:created_at"`
		UpdatedAt string `gorm:"column:updated_at"`
	}
	var out categoryRow
	h.db.Raw(`SELECT c.id, c.name, COALESCE(c.icon,'') AS icon, COALESCE(c.image,'') AS image, c.sort_order, c.status, COUNT(d.id) AS count, c.created_at, c.updated_at FROM categories c LEFT JOIN dishes d ON d.category_id=c.id WHERE c.id=? GROUP BY c.id`, id).Scan(&out)
	var name models.I18nString
	_ = json.Unmarshal([]byte(out.NameJSON), &name)
	cat := models.Category{
		ID:        out.ID,
		Name:      name,
		Icon:      out.Icon,
		Image:     out.Image,
		SortOrder: out.SortOrder,
		Status:    out.Status,
		DishCount: out.DishCount,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "categories:list")
	}
	OK(w, cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))

	var dishCount int64
	if err := h.db.Raw(`SELECT COUNT(*) FROM dishes WHERE category_id=?`, id).Scan(&dishCount).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if dishCount > 0 {
		Fail(w, http.StatusConflict, "category has dishes; cannot delete")
		return
	}

	result := h.db.Exec(`DELETE FROM categories WHERE id=?`, id)
	if result.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if result.RowsAffected == 0 {
		NotFound(w, "category not found")
		return
	}

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "categories:list")
	}
	OK(w, nil)
}

func (h *CategoryHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	var body struct {
		OrderedIDs []string `json:"orderedIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	for i, id := range body.OrderedIDs {
		_ = h.db.Exec(`UPDATE categories SET sort_order=?, updated_at=? WHERE id=?`, i+1, now, id).Error
	}

	cats, err := h.listCategories()
	if err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "categories:list")
	}
	OK(w, cats)
}
