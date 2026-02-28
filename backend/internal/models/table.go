package models

type Table struct {
	ID             string  `json:"id"`
	Number         string  `json:"number"`
	Seats          int     `json:"seats"`
	Status         string  `json:"status"`
	CurrentOrderID *string `json:"currentOrderId,omitempty"`
	QRCodeURL      string  `json:"qrCodeUrl,omitempty"`
	Area           string  `json:"area,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}
