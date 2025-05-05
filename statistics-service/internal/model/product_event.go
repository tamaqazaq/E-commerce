package model

import "time"

type ProductEvent struct {
	ProductID string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Timestamp time.Time `json:"time"`
	Action    string    // created, updated, deleted
}
