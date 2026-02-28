package models

type SelectedOption struct {
	OptionName      string  `json:"optionName"`
	ValueName       string  `json:"valueName"`
	PriceAdjustment float64 `json:"priceAdjustment"`
}

type OrderItem struct {
	DishID          string           `json:"dishId"`
	DishName        I18nString       `json:"dishName"`
	Quantity        int              `json:"quantity"`
	UnitPrice       float64          `json:"unitPrice"`
	SelectedOptions []SelectedOption `json:"selectedOptions"`
	Subtotal        float64          `json:"subtotal"`
}

type StatusHistoryEntry struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Note      string `json:"note,omitempty"`
}

type PaymentInfo struct {
	Method string  `json:"method"`
	Status string  `json:"status"`
	PaidAt *string `json:"paidAt,omitempty"`
}

type Order struct {
	ID                    string               `json:"id"`
	OrderNumber           string               `json:"orderNumber"`
	Type                  string               `json:"type"`
	Status                string               `json:"status"`
	TableID               *string              `json:"tableId,omitempty"`
	Items                 []OrderItem          `json:"items"`
	Subtotal              float64              `json:"subtotal"`
	DeliveryFee           float64              `json:"deliveryFee"`
	Discount              float64              `json:"discount"`
	Total                 float64              `json:"total"`
	Remarks               *string              `json:"remarks,omitempty"`
	Payment               PaymentInfo          `json:"payment"`
	DeliveryAddress       *string              `json:"deliveryAddress,omitempty"`
	ContactPhone          *string              `json:"contactPhone,omitempty"`
	EstimatedDeliveryTime *string              `json:"estimatedDeliveryTime,omitempty"`
	StatusHistory         []StatusHistoryEntry `json:"statusHistory"`
	CreatedAt             string               `json:"createdAt"`
	UpdatedAt             string               `json:"updatedAt"`
}

type CreateOrderPayload struct {
	Type            string            `json:"type"`
	TableID         string            `json:"tableId"`
	Items           []CreateOrderItem `json:"items"`
	Remarks         string            `json:"remarks"`
	DeliveryAddress string            `json:"deliveryAddress"`
	ContactPhone    string            `json:"contactPhone"`
	PaymentMethod   string            `json:"paymentMethod"`
}

type CreateOrderItem struct {
	DishID          string           `json:"dishId"`
	Quantity        int              `json:"quantity"`
	SelectedOptions []SelectedOption `json:"selectedOptions"`
}
