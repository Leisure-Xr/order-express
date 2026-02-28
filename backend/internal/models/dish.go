package models

type DishOptionValue struct {
	ID              string     `json:"id"`
	Label           I18nString `json:"label"`
	PriceAdjustment float64    `json:"priceAdjustment"`
}

type DishOption struct {
	ID          string            `json:"id"`
	Name        I18nString        `json:"name"`
	Values      []DishOptionValue `json:"values"`
	Required    bool              `json:"required"`
	MultiSelect bool              `json:"multiSelect"`
}

type Dish struct {
	ID              string       `json:"id"`
	CategoryID      string       `json:"categoryId"`
	Name            I18nString   `json:"name"`
	Description     I18nString   `json:"description"`
	Price           float64      `json:"price"`
	OriginalPrice   *float64     `json:"originalPrice,omitempty"`
	Image           string       `json:"image"`
	Images          []string     `json:"images,omitempty"`
	Status          string       `json:"status"`
	Options         []DishOption `json:"options"`
	Tags            []string     `json:"tags"`
	PreparationTime *int         `json:"preparationTime,omitempty"`
	CreatedAt       string       `json:"createdAt"`
	UpdatedAt       string       `json:"updatedAt"`
}
