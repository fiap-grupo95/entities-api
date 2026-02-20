package response

import "time"

// AdditionalRepairServiceResponse represents a service returned to the client.
type AdditionalRepairServiceResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name,omitempty"`
	Price float64 `json:"price,omitempty"`
}

// AdditionalRepairPartsSupplyResponse represents a parts supply returned to the client.
type AdditionalRepairPartsSupplyResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name,omitempty"`
	Price           float64 `json:"price,omitempty"`
	QuantityReserve int     `json:"quantity_reserve"`
	QuantityTotal   int     `json:"quantity_total"`
}

// AdditionalRepairResponse encapsulates the data returned for a single additional repair.
type AdditionalRepairResponse struct {
	ID             uint                                  `json:"id"`
	ServiceOrderID uint                                  `json:"service_order_id"`
	Description    string                                `json:"description"`
	Status         string                                `json:"status"`
	Estimate       float64                               `json:"estimate"`
	CreatedAt      time.Time                             `json:"created_at"`
	UpdatedAt      time.Time                             `json:"updated_at"`
	Services       []AdditionalRepairServiceResponse     `json:"services"`
	PartsSupplies  []AdditionalRepairPartsSupplyResponse `json:"parts_supplies"`
}

// OperationMessageResponse represents a simple operation feedback payload.
type OperationMessageResponse struct {
	Message string `json:"message"`
}
