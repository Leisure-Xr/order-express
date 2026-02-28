package models

type Category struct {
	ID        string     `json:"id"`
	Name      I18nString `json:"name"`
	Icon      string     `json:"icon,omitempty"`
	Image     string     `json:"image,omitempty"`
	SortOrder int        `json:"sortOrder"`
	Status    string     `json:"status"`
	DishCount int        `json:"dishCount,omitempty"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}
