package response

import "time"

// PaymentResponse represents the payload returned for payment operations.
type PaymentResponse struct {
	ID             uint      `json:"id"`
	ServiceOrderID uint      `json:"service_order_id"`
	PaymentDate    time.Time `json:"payment_date"`
	Amount         float64   `json:"amount"`
}
