package model

import "time"

type OrderEvent struct {
	OrderID   string
	UserID    string
	Total     float64
	Status    string
	Timestamp time.Time
	Action    string
}
