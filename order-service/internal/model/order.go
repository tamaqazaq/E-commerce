package model

type Order struct {
	ID     string      `json:"id"`
	UserID string      `json:"user_id"`
	Items  []OrderItem `json:"items"`
	Total  float64     `json:"total"`
	Status string      `json:"status"`
}
