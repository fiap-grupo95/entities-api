package entities

import (
	"time"
)

type Payment struct {
	ID             uint          `json:"id"`
	ServiceOrderID uint          `json:"service_order_id"`
	ServiceOrder   *ServiceOrder `json:"service_order,omitempty"`
	PaymentDate    time.Time     `json:"payment_date"`
	Amount         float64       `json:"amount"`
}
