package models

type Payment struct {
	PaymentID string  `json:"paymentId"`
	OrderID   string  `json:"orderId"`
	Method    string  `json:"method"`
	Status    string  `json:"status"`
	Amount    float64 `json:"amount"`
	CreatedAt string  `json:"createdAt"`
	PaidAt    *string `json:"paidAt,omitempty"`
}
