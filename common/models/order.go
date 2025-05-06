package models

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Products  []OrderItem `json:"products"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
