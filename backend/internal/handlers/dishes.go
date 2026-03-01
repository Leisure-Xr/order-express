package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"order-express/backend/internal/cache"
	"order-express/backend/internal/idgen"
	"order-express/backend/internal/models"

	"gorm.io/gorm"
)

type DishHandler struct {
	db       *gorm.DB
	cache    *cache.Cache
	cacheTTL time.Duration
}

func NewDishHandler(db *gorm.DB, c *cache.Cache, ttl time.Duration) *DishHandler {
	return &DishHandler{db: db, cache: c, cacheTTL: ttl}
}

func (h *DishHandler) scanDish(row dishRow) models.Dish {
	var d models.Dish
	d.ID = row.ID
	d.CategoryID = row.CategoryID
	d.Price = row.Price
	d.Image = row.Image
	d.Status = row.Status
	d.CreatedAt = row.CreatedAt
	d.UpdatedAt = row.UpdatedAt

	_ = json.Unmarshal([]byte(row.NameJSON), &d.Name)
	_ = json.Unmarshal([]byte(row.DescJSON), &d.Description)
	_ = json.Unmarshal([]byte(row.OptionsJSON), &d.Options)
	_ = json.Unmarshal([]byte(row.TagsJSON), &d.Tags)
	_ = json.Unmarshal([]byte(row.ImagesJSON), &d.Images)

	if d.Options == nil {
		d.Options = []models.DishOption{}
	}
	if d.Tags == nil {
		d.Tags = []string{}
	}
	if d.Images == nil {
		d.Images = []string{}
	}
	if row.OriginalPrice.Valid {
		v := row.OriginalPrice.Float64
		d.OriginalPrice = &v
	}
	if row.PrepTime.Valid {
		v := int(row.PrepTime.Int64)
		d.PreparationTime = &v
	}
	return d
}

type dishRow struct {
	ID            string          `gorm:"column:id"`
	CategoryID    string          `gorm:"column:category_id"`
	NameJSON      string          `gorm:"column:name"`
	DescJSON      string          `gorm:"column:description"`
	Price         float64         `gorm:"column:price"`
	OriginalPrice sql.NullFloat64 `gorm:"column:original_price"`
	Image         string          `gorm:"column:image"`
	ImagesJSON    string          `gorm:"column:images"`
	Status        string          `gorm:"column:status"`
	OptionsJSON   string          `gorm:"column:options"`
	TagsJSON      string          `gorm:"column:tags"`
	PrepTime      sql.NullInt64   `gorm:"column:preparation_time"`
	CreatedAt     string          `gorm:"column:created_at"`
	UpdatedAt     string          `gorm:"column:updated_at"`
}

func (h *DishHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	categoryID := strings.TrimSpace(q.Get("categoryId"))
	status := strings.TrimSpace(q.Get("status"))

	cacheKey := h.dishListCacheKey(categoryID, status)
	if h.cache != nil {
		if payload, ok, err := h.cache.GetString(r.Context(), cacheKey); err == nil && ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
			return
		}
	}

	query := `SELECT id, category_id, name, description, price, original_price, image, COALESCE(images,'[]') AS images, status, options, tags, preparation_time, created_at, updated_at FROM dishes WHERE 1=1`
	args := []any{}
	if categoryID != "" {
		query += " AND category_id = ?"
		args = append(args, categoryID)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	query += " ORDER BY created_at"

	rows := []dishRow{}
	tx := h.db.Raw(query, args...).Scan(&rows)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	dishes := make([]models.Dish, 0, len(rows))
	for _, row := range rows {
		dishes = append(dishes, h.scanDish(row))
	}

	resp := ApiResponse{Code: 200, Data: dishes, Message: "success"}
	b, err := json.Marshal(resp)
	if err != nil {
		OK(w, dishes)
		return
	}
	if h.cache != nil {
		_ = h.cache.SetString(r.Context(), cacheKey, string(b), h.cacheTTL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h *DishHandler) dishListCacheKey(categoryID string, status string) string {
	if categoryID == "" && status == "" {
		return "dishes:list"
	}
	return "dishes:list?categoryId=" + url.QueryEscape(categoryID) + "&status=" + url.QueryEscape(status)
}

func (h *DishHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body models.Dish
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	var catCount int64
	if err := h.db.Raw(`SELECT COUNT(*) FROM categories WHERE id=?`, body.CategoryID).Scan(&catCount).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if catCount == 0 {
		BadRequest(w, "category not found")
		return
	}

	if body.Price <= 0 {
		BadRequest(w, "price must be greater than 0")
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id, err := idgen.New("dish-", 8)
	if err != nil {
		Fail(w, http.StatusInternalServerError, "internal error")
		return
	}
	nameJSON, _ := json.Marshal(body.Name)
	descJSON, _ := json.Marshal(body.Description)
	optJSON, _ := json.Marshal(body.Options)
	tagsJSON, _ := json.Marshal(body.Tags)
	imagesJSON, _ := json.Marshal(body.Images)
	if string(imagesJSON) == "null" {
		imagesJSON = []byte("[]")
	}
	if string(optJSON) == "null" {
		optJSON = []byte("[]")
	}
	if string(tagsJSON) == "null" {
		tagsJSON = []byte("[]")
	}

	tx := h.db.Exec(
		`INSERT INTO dishes (id, category_id, name, description, price, original_price, image, images, status, options, tags, preparation_time, created_at, updated_at)
         VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		id, body.CategoryID, string(nameJSON), string(descJSON), body.Price, body.OriginalPrice,
		body.Image, string(imagesJSON), body.Status, string(optJSON), string(tagsJSON),
		body.PreparationTime, now, now,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	if h.cache != nil {
		_ = h.cache.DelPrefix(r.Context(), "dishes:list")
		_ = h.cache.Del(r.Context(), "categories:list")
	}

	h.getAndReturn(w, id)
}

func (h *DishHandler) getAndReturn(w http.ResponseWriter, id string) {
	var row dishRow
	tx := h.db.Raw(`SELECT id, category_id, name, description, price, original_price, image, COALESCE(images,'[]') AS images, status, options, tags, preparation_time, created_at, updated_at FROM dishes WHERE id=?`, id).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "dish not found")
		return
	}

	OK(w, h.scanDish(row))
}

func (h *DishHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.getAndReturn(w, r.PathValue("id"))
}

func (h *DishHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Load existing
	var existingRow dishRow
	tx := h.db.Raw(`SELECT id, category_id, name, description, price, original_price, image, COALESCE(images,'[]') AS images, status, options, tags, preparation_time, created_at, updated_at FROM dishes WHERE id=?`, id).Scan(&existingRow)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "dish not found")
		return
	}
	existing := h.scanDish(existingRow)

	// Decode patch
	var patch map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	if v, ok := patch["categoryId"]; ok {
		json.Unmarshal(v, &existing.CategoryID)
	}
	if v, ok := patch["name"]; ok {
		json.Unmarshal(v, &existing.Name)
	}
	if v, ok := patch["description"]; ok {
		json.Unmarshal(v, &existing.Description)
	}
	if v, ok := patch["price"]; ok {
		json.Unmarshal(v, &existing.Price)
	}
	if v, ok := patch["originalPrice"]; ok {
		json.Unmarshal(v, &existing.OriginalPrice)
	}
	if v, ok := patch["image"]; ok {
		json.Unmarshal(v, &existing.Image)
	}
	if v, ok := patch["images"]; ok {
		json.Unmarshal(v, &existing.Images)
	}
	if v, ok := patch["status"]; ok {
		json.Unmarshal(v, &existing.Status)
	}
	if v, ok := patch["options"]; ok {
		json.Unmarshal(v, &existing.Options)
	}
	if v, ok := patch["tags"]; ok {
		json.Unmarshal(v, &existing.Tags)
	}
	if v, ok := patch["preparationTime"]; ok {
		json.Unmarshal(v, &existing.PreparationTime)
	}

	var catCount int64
	if err := h.db.Raw(`SELECT COUNT(*) FROM categories WHERE id=?`, existing.CategoryID).Scan(&catCount).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if catCount == 0 {
		BadRequest(w, "category not found")
		return
	}

	if existing.Price <= 0 {
		BadRequest(w, "price must be greater than 0")
		return
	}

	nameJSON, _ := json.Marshal(existing.Name)
	descJSON, _ := json.Marshal(existing.Description)
	optJSON, _ := json.Marshal(existing.Options)
	tagsJSON, _ := json.Marshal(existing.Tags)
	imagesJSON, _ := json.Marshal(existing.Images)
	now := time.Now().UTC().Format(time.RFC3339)

	tx = h.db.Exec(
		`UPDATE dishes SET category_id=?, name=?, description=?, price=?, original_price=?, image=?, images=?, status=?, options=?, tags=?, preparation_time=?, updated_at=? WHERE id=?`,
		existing.CategoryID, string(nameJSON), string(descJSON), existing.Price, existing.OriginalPrice,
		existing.Image, string(imagesJSON), existing.Status, string(optJSON), string(tagsJSON),
		existing.PreparationTime, now, id,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	if h.cache != nil {
		_ = h.cache.DelPrefix(r.Context(), "dishes:list")
		_ = h.cache.Del(r.Context(), "categories:list")
	}

	h.getAndReturn(w, id)
}

func (h *DishHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result := h.db.Exec(`DELETE FROM dishes WHERE id=?`, id)
	if result.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if result.RowsAffected == 0 {
		NotFound(w, "dish not found")
		return
	}

	if h.cache != nil {
		_ = h.cache.DelPrefix(r.Context(), "dishes:list")
		_ = h.cache.Del(r.Context(), "categories:list")
	}
	OK(w, nil)
}

func (h *DishHandler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var status string
	tx := h.db.Raw(`SELECT status FROM dishes WHERE id=?`, id).Scan(&status)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "dish not found")
		return
	}

	newStatus := "off_sale"
	if status == "off_sale" || status == "sold_out" {
		newStatus = "on_sale"
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_ = h.db.Exec(`UPDATE dishes SET status=?, updated_at=? WHERE id=?`, newStatus, now, id).Error

	if h.cache != nil {
		_ = h.cache.DelPrefix(r.Context(), "dishes:list")
		_ = h.cache.Del(r.Context(), "categories:list")
	}

	h.getAndReturn(w, id)
}
