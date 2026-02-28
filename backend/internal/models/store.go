package models

type StoreInfo struct {
	Name        I18nString `json:"name"`
	Address     I18nString `json:"address"`
	Phone       string     `json:"phone"`
	Logo        string     `json:"logo"`
	Description I18nString `json:"description"`
}

type BusinessHours struct {
	DayOfWeek int    `json:"dayOfWeek"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	IsClosed  bool   `json:"isClosed"`
}

type DeliverySettings struct {
	Enabled               bool    `json:"enabled"`
	MinimumOrder          float64 `json:"minimumOrder"`
	DeliveryFee           float64 `json:"deliveryFee"`
	FreeDeliveryThreshold float64 `json:"freeDeliveryThreshold"`
	EstimatedMinutes      int     `json:"estimatedMinutes"`
	DeliveryRadius        float64 `json:"deliveryRadius"`
}
