package entities

import (
	"mecanica_xpto/internal/domain/valueobject"
	"time"
)

type Vehicle struct {
	ID         uint              `json:"id"`
	Plate      valueobject.Plate `json:"plate"`
	CustomerID uint              `json:"customer_id"`
	Customer   *Customer         `json:"customer,omitempty"`
	Model      string            `json:"model"`
	Year       string            `json:"year"`
	Brand      string            `json:"brand"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  *time.Time        `json:"updated_at,"`
	DeletedAt  *time.Time        `json:"deleted_at,omitempty"`
	//ServiceOrders []ServiceOrder    `json:"service_orders,omitempty"`
}
