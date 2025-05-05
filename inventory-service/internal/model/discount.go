package model

import "time"

type Discount struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DiscountPercentage float64   `json:"discount_percentage"`
	ApplicableProducts []string  `json:"applicable_products"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	IsActive           bool      `json:"is_active"`
}
