package model

import "time"

type OrderEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Total     float64   `json:"total"`
	Timestamp time.Time `json:"time"`
	Action    string
}
