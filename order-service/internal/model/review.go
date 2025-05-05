package model

import "time"

type Review struct {
	ID        string
	ProductID string
	UserID    string
	Rating    float64
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
