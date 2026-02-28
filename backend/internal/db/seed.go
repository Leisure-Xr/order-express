package db

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	now := time.Now().UTC().Format(time.RFC3339)

	// Users (admin)
	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := db.Exec(
		`INSERT INTO users (id, username, name, password, role, created_at, updated_at)
         VALUES (?,?,?,?,?,?,?)
         ON CONFLICT (username) DO NOTHING`,
		"user-admin", "admin", "管理员", string(adminHash), "admin", now, now,
	).Error; err != nil {
		return err
	}

	// Categories
	type catSeed struct {
		id        string
		nameZH    string
		nameEN    string
		icon      string
		image     string
		sortOrder int
		status    string
	}
	cats := []catSeed{
		{"cat-01", "热菜", "Hot Dishes", "🍲", "", 1, "active"},
		{"cat-02", "凉菜", "Cold Dishes", "🥗", "", 2, "active"},
		{"cat-03", "饮品", "Drinks", "🥤", "", 3, "active"},
	}
	for _, c := range cats {
		nameJSON := jsonStr(map[string]string{"zh": c.nameZH, "en": c.nameEN})
		if err := db.Exec(
			`INSERT INTO categories (id, name, icon, image, sort_order, status, created_at, updated_at)
             VALUES (?,?,?,?,?,?,?,?)
             ON CONFLICT (id) DO NOTHING`,
			c.id, nameJSON, c.icon, c.image, c.sortOrder, c.status, now, now,
		).Error; err != nil {
			return fmt.Errorf("insert category %s: %w", c.id, err)
		}
	}

	// Dishes
	type dishSeed struct {
		id            string
		catID         string
		nameZH        string
		nameEN        string
		descZH        string
		descEN        string
		price         float64
		originalPrice *float64
		image         string
		status        string
		optionsJSON   string
		tags          []string
		prepTime      *int
	}
	emptyOpts := "[]"
	dishes := []dishSeed{
		{"dish-001", "cat-01", "宫保鸡丁", "Kung Pao Chicken", "经典川味，香辣微甜", "Classic Sichuan stir-fry with peanuts", 32, nil, "https://placehold.co/400x300/e74c3c/white?text=KungPao", "on_sale", emptyOpts, []string{"popular", "spicy"}, ptrI(12)},
		{"dish-002", "cat-01", "麻婆豆腐", "Mapo Tofu", "麻辣鲜香，豆腐嫩滑", "Spicy tofu with minced meat", 26, nil, "https://placehold.co/400x300/e74c3c/white?text=MapoTofu", "on_sale", emptyOpts, []string{"spicy"}, ptrI(10)},
		{"dish-003", "cat-02", "拍黄瓜", "Smashed Cucumber", "清爽开胃", "Refreshing cucumber salad", 14, nil, "https://placehold.co/400x300/2ecc71/white?text=Cucumber", "on_sale", emptyOpts, []string{"vegetarian"}, ptrI(5)},
		{"dish-004", "cat-03", "可乐", "Cola", "冰爽畅饮", "Classic soda", 8, nil, "https://placehold.co/400x300/3498db/white?text=Cola", "on_sale", emptyOpts, []string{}, nil},
	}

	for _, d := range dishes {
		var descJSON string
		nameJSON := jsonStr(map[string]string{"zh": d.nameZH, "en": d.nameEN})
		descJSON = jsonStr(map[string]string{"zh": d.descZH, "en": d.descEN})
		tagsJSON := jsonStr(d.tags)
		imagesJSON := "[]"

		if err := db.Exec(
			`INSERT INTO dishes (id, category_id, name, description, price, original_price, image, images, status, options, tags, preparation_time, created_at, updated_at)
             VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
             ON CONFLICT (id) DO NOTHING`,
			d.id, d.catID, nameJSON, descJSON, d.price, d.originalPrice, d.image, imagesJSON, d.status, d.optionsJSON, tagsJSON, d.prepTime, now, now,
		).Error; err != nil {
			return fmt.Errorf("insert dish %s: %w", d.id, err)
		}
	}

	// Tables
	type tableSeed struct {
		id     string
		number string
		area   string
		seats  int
		status string
	}
	tables := []tableSeed{
		{"table-01", "T01", "Main Hall", 2, "available"},
		{"table-02", "T02", "Main Hall", 2, "available"},
		{"table-03", "T03", "Main Hall", 2, "available"},
		{"table-04", "T04", "Main Hall", 2, "available"},
		{"table-05", "T05", "Main Hall", 4, "available"},
		{"table-06", "T06", "Main Hall", 4, "available"},
		{"table-07", "T07", "Main Hall", 4, "available"},
		{"table-08", "T08", "Main Hall", 4, "available"},
		{"table-09", "T09", "Main Hall", 6, "available"},
		{"table-10", "T10", "Main Hall", 6, "available"},
		{"table-11", "T11", "Private Room", 8, "available"},
		{"table-12", "T12", "Private Room", 8, "available"},
	}
	for _, t := range tables {
		qrURL := "/qr/" + t.id
		if err := db.Exec(
			`INSERT INTO tables (id, number, seats, status, current_order_id, qr_code_url, area, created_at, updated_at)
             VALUES (?,?,?,?,NULL,?,?,?,?)
             ON CONFLICT (id) DO NOTHING`,
			t.id, t.number, t.seats, t.status, qrURL, t.area, now, now,
		).Error; err != nil {
			return fmt.Errorf("insert table %s: %w", t.id, err)
		}
	}

	// Store info
	storeNameJSON := jsonStr(map[string]string{"zh": "食光小馆", "en": "Order Express"})
	storeAddrJSON := jsonStr(map[string]string{"zh": "北京市朝阳区建国路88号", "en": "88 Jianguo Road, Chaoyang, Beijing"})
	storeDescJSON := jsonStr(map[string]string{"zh": "精选时令食材，匠心烹饪美味", "en": "Fresh seasonal ingredients, crafted with care"})
	if err := db.Exec(
		`INSERT INTO store_info (id, name, address, phone, logo, description, updated_at)
         VALUES (1,?,?,?,?,?,?)
         ON CONFLICT (id) DO NOTHING`,
		storeNameJSON, storeAddrJSON, "010-88886666", "/images/logo.png", storeDescJSON, now,
	).Error; err != nil {
		return err
	}

	// Business hours
	type bhSeed struct {
		day      int
		open     string
		close    string
		isClosed bool
	}
	bhs := []bhSeed{
		{0, "", "", true},
		{1, "10:00", "22:00", false},
		{2, "10:00", "22:00", false},
		{3, "10:00", "22:00", false},
		{4, "10:00", "22:00", false},
		{5, "10:00", "23:00", false},
		{6, "10:00", "23:00", false},
	}
	for _, bh := range bhs {
		if err := db.Exec(
			`INSERT INTO business_hours (day_of_week, open_time, close_time, is_closed)
             VALUES (?,?,?,?)
             ON CONFLICT (day_of_week) DO NOTHING`,
			bh.day, bh.open, bh.close, bh.isClosed,
		).Error; err != nil {
			return err
		}
	}

	// Delivery settings
	if err := db.Exec(
		`INSERT INTO delivery_settings (id, enabled, minimum_order, delivery_fee, free_delivery_threshold, estimated_minutes, delivery_radius, updated_at)
         VALUES (1,TRUE,30,5,80,30,5,?)
         ON CONFLICT (id) DO NOTHING`,
		now,
	).Error; err != nil {
		return err
	}

	return nil
}

func ptrI(v int) *int { return &v }

func jsonStr(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
