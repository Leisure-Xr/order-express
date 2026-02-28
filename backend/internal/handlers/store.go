package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"order-express/backend/internal/cache"
	"order-express/backend/internal/models"

	"gorm.io/gorm"
)

type StoreHandler struct {
	db       *gorm.DB
	cache    *cache.Cache
	cacheTTL time.Duration
}

func NewStoreHandler(db *gorm.DB, c *cache.Cache, ttl time.Duration) *StoreHandler {
	return &StoreHandler{db: db, cache: c, cacheTTL: ttl}
}

func (h *StoreHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	if h.cache != nil {
		if payload, ok, err := h.cache.GetString(r.Context(), "store:info"); err == nil && ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
			return
		}
	}

	var info struct {
		NameJSON string `gorm:"column:name"`
		AddrJSON string `gorm:"column:address"`
		Phone    string `gorm:"column:phone"`
		Logo     string `gorm:"column:logo"`
		DescJSON string `gorm:"column:description"`
	}
	tx := h.db.Raw(`SELECT name, address, phone, logo, description FROM store_info WHERE id=1`).Scan(&info)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "store info not found")
		return
	}

	var out models.StoreInfo
	_ = json.Unmarshal([]byte(info.NameJSON), &out.Name)
	_ = json.Unmarshal([]byte(info.AddrJSON), &out.Address)
	_ = json.Unmarshal([]byte(info.DescJSON), &out.Description)
	out.Phone = info.Phone
	out.Logo = info.Logo

	resp := ApiResponse{Code: 200, Data: out, Message: "success"}
	b, err := json.Marshal(resp)
	if err != nil {
		OK(w, out)
		return
	}
	if h.cache != nil {
		_ = h.cache.SetString(r.Context(), "store:info", string(b), h.cacheTTL)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h *StoreHandler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	var existing struct {
		NameJSON string `gorm:"column:name"`
		AddrJSON string `gorm:"column:address"`
		Phone    string `gorm:"column:phone"`
		Logo     string `gorm:"column:logo"`
		DescJSON string `gorm:"column:description"`
	}
	tx := h.db.Raw(`SELECT name, address, phone, logo, description FROM store_info WHERE id=1`).Scan(&existing)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "store info not found")
		return
	}

	var info models.StoreInfo
	_ = json.Unmarshal([]byte(existing.NameJSON), &info.Name)
	_ = json.Unmarshal([]byte(existing.AddrJSON), &info.Address)
	_ = json.Unmarshal([]byte(existing.DescJSON), &info.Description)
	info.Phone = existing.Phone
	info.Logo = existing.Logo

	var patch map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	if v, ok := patch["name"]; ok {
		_ = json.Unmarshal(v, &info.Name)
	}
	if v, ok := patch["address"]; ok {
		_ = json.Unmarshal(v, &info.Address)
	}
	if v, ok := patch["phone"]; ok {
		_ = json.Unmarshal(v, &info.Phone)
	}
	if v, ok := patch["logo"]; ok {
		_ = json.Unmarshal(v, &info.Logo)
	}
	if v, ok := patch["description"]; ok {
		_ = json.Unmarshal(v, &info.Description)
	}

	updNameJSON, _ := json.Marshal(info.Name)
	updAddrJSON, _ := json.Marshal(info.Address)
	updDescJSON, _ := json.Marshal(info.Description)
	now := time.Now().UTC().Format(time.RFC3339)

	tx = h.db.Exec(
		`UPDATE store_info SET name=?, address=?, phone=?, logo=?, description=?, updated_at=? WHERE id=1`,
		string(updNameJSON), string(updAddrJSON), info.Phone, info.Logo, string(updDescJSON), now,
	)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "store:info")
	}
	OK(w, info)
}

func (h *StoreHandler) GetBusinessHours(w http.ResponseWriter, r *http.Request) {
	if h.cache != nil {
		if payload, ok, err := h.cache.GetString(r.Context(), "store:business_hours"); err == nil && ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
			return
		}
	}

	type row struct {
		DayOfWeek int    `gorm:"column:day_of_week"`
		OpenTime  string `gorm:"column:open_time"`
		CloseTime string `gorm:"column:close_time"`
		IsClosed  bool   `gorm:"column:is_closed"`
	}
	rows := []row{}
	if err := h.db.Raw(`SELECT day_of_week, open_time, close_time, is_closed FROM business_hours ORDER BY day_of_week`).Scan(&rows).Error; err != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}

	hours := make([]models.BusinessHours, 0, len(rows))
	for _, r := range rows {
		hours = append(hours, models.BusinessHours{
			DayOfWeek: r.DayOfWeek,
			OpenTime:  r.OpenTime,
			CloseTime: r.CloseTime,
			IsClosed:  r.IsClosed,
		})
	}

	resp := ApiResponse{Code: 200, Data: hours, Message: "success"}
	b, err := json.Marshal(resp)
	if err != nil {
		OK(w, hours)
		return
	}
	if h.cache != nil {
		_ = h.cache.SetString(r.Context(), "store:business_hours", string(b), h.cacheTTL)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h *StoreHandler) UpdateBusinessHours(w http.ResponseWriter, r *http.Request) {
	var hours []models.BusinessHours
	if err := json.NewDecoder(r.Body).Decode(&hours); err != nil {
		BadRequest(w, "invalid request body")
		return
	}

	for _, bh := range hours {
		_ = h.db.Exec(
			`INSERT INTO business_hours (day_of_week, open_time, close_time, is_closed) VALUES (?,?,?,?)
             ON CONFLICT(day_of_week) DO UPDATE SET open_time=excluded.open_time, close_time=excluded.close_time, is_closed=excluded.is_closed`,
			bh.DayOfWeek, bh.OpenTime, bh.CloseTime, bh.IsClosed,
		).Error
	}

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "store:business_hours")
	}
	h.GetBusinessHours(w, r)
}

func (h *StoreHandler) GetDeliverySettings(w http.ResponseWriter, r *http.Request) {
	if h.cache != nil {
		if payload, ok, err := h.cache.GetString(r.Context(), "store:delivery_settings"); err == nil && ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(payload))
			return
		}
	}

	var row struct {
		Enabled               bool    `gorm:"column:enabled"`
		MinimumOrder          float64 `gorm:"column:minimum_order"`
		DeliveryFee           float64 `gorm:"column:delivery_fee"`
		FreeDeliveryThreshold float64 `gorm:"column:free_delivery_threshold"`
		EstimatedMinutes      int     `gorm:"column:estimated_minutes"`
		DeliveryRadius        float64 `gorm:"column:delivery_radius"`
	}
	tx := h.db.Raw(`SELECT enabled, minimum_order, delivery_fee, free_delivery_threshold, estimated_minutes, delivery_radius FROM delivery_settings WHERE id=1`).Scan(&row)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "delivery settings not found")
		return
	}

	ds := models.DeliverySettings{
		Enabled:               row.Enabled,
		MinimumOrder:          row.MinimumOrder,
		DeliveryFee:           row.DeliveryFee,
		FreeDeliveryThreshold: row.FreeDeliveryThreshold,
		EstimatedMinutes:      row.EstimatedMinutes,
		DeliveryRadius:        row.DeliveryRadius,
	}

	resp := ApiResponse{Code: 200, Data: ds, Message: "success"}
	b, err := json.Marshal(resp)
	if err != nil {
		OK(w, ds)
		return
	}
	if h.cache != nil {
		_ = h.cache.SetString(r.Context(), "store:delivery_settings", string(b), h.cacheTTL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h *StoreHandler) UpdateDeliverySettings(w http.ResponseWriter, r *http.Request) {
	var existing models.DeliverySettings
	tx := h.db.Raw(`SELECT enabled, minimum_order, delivery_fee, free_delivery_threshold, estimated_minutes, delivery_radius FROM delivery_settings WHERE id=1`).Scan(&existing)
	if tx.Error != nil {
		Fail(w, http.StatusInternalServerError, "database error")
		return
	}
	if tx.RowsAffected == 0 {
		NotFound(w, "delivery settings not found")
		return
	}

	var patch map[string]json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		BadRequest(w, "invalid request body")
		return
	}
	if v, ok := patch["enabled"]; ok {
		_ = json.Unmarshal(v, &existing.Enabled)
	}
	if v, ok := patch["minimumOrder"]; ok {
		_ = json.Unmarshal(v, &existing.MinimumOrder)
	}
	if v, ok := patch["deliveryFee"]; ok {
		_ = json.Unmarshal(v, &existing.DeliveryFee)
	}
	if v, ok := patch["freeDeliveryThreshold"]; ok {
		_ = json.Unmarshal(v, &existing.FreeDeliveryThreshold)
	}
	if v, ok := patch["estimatedMinutes"]; ok {
		_ = json.Unmarshal(v, &existing.EstimatedMinutes)
	}
	if v, ok := patch["deliveryRadius"]; ok {
		_ = json.Unmarshal(v, &existing.DeliveryRadius)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_ = h.db.Exec(
		`UPDATE delivery_settings SET enabled=?, minimum_order=?, delivery_fee=?, free_delivery_threshold=?, estimated_minutes=?, delivery_radius=?, updated_at=? WHERE id=1`,
		existing.Enabled, existing.MinimumOrder, existing.DeliveryFee, existing.FreeDeliveryThreshold, existing.EstimatedMinutes, existing.DeliveryRadius, now,
	).Error

	if h.cache != nil {
		_ = h.cache.Del(r.Context(), "store:delivery_settings")
	}
	OK(w, existing)
}
