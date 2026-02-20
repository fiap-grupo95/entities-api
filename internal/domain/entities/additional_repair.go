package entities

import (
	"mecanica_xpto/internal/domain/valueobject"
	"time"
)

type AdditionalRepair struct {
	ID             uint                               `json:"id"`
	Description    string                             `json:"description"`
	ServiceOrderID uint                               `json:"service_order_id"`
	ServiceOrder   *ServiceOrder                      `json:"service_order,omitempty"`
	ARStatus       valueobject.AdditionalRepairStatus `json:"ar_status,omitempty"`
	Estimate       float64                            `json:"estimate"`
	CreatedAt      time.Time                          `json:"created_at"`
	UpdatedAt      time.Time                          `json:"updated_at"`
	PartsSupplies  []PartsSupply                      `json:"parts_supplies,omitempty"`
	Services       []Service                          `json:"services,omitempty"`
}

type AdditionalRepairStatusDTO struct {
	ApprovalStatus string `json:"approval_status" binding:"required,oneof=APPROVED DENIED"`
}
