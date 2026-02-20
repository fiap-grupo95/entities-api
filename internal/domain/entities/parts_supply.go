package entities

import (
	"time"
)

type PartsSupply struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Price           float64    `json:"price"`
	QuantityTotal   int        `json:"quantity_total"`
	QuantityReserve int        `json:"quantity_reserve"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
