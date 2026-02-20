package entities

import (
	"mecanica_xpto/internal/domain/valueobject"
	"time"
)

type ServiceOrder struct {
	ID                       uint                           `json:"id"`
	CustomerID               uint                           `json:"customer_id"`
	Customer                 *Customer                      `json:"customer,omitempty"`
	VehicleID                uint                           `json:"vehicle_id"`
	Vehicle                  *Vehicle                       `json:"vehicle,omitempty"`
	ServiceOrderStatus       valueobject.ServiceOrderStatus `json:"service_order_status"`
	Estimate                 float64                        `json:"estimate,omitempty"`
	StartedExecutionDate     *time.Time                     `json:"started_execution_date,omitempty"`
	FinalExecutionDate       *time.Time                     `json:"final_execution_date,omitempty"`
	ExecutionDurationInHours float64                        `json:"execution_duration_in_hours,omitempty"`
	CreatedAt                *time.Time                     `json:"created_at,omitempty"`
	UpdatedAt                *time.Time                     `json:"updated_at,omitempty"`
	Payment                  *Payment                       `json:"payment,omitempty"`
	AdditionalRepairs        []AdditionalRepair             `json:"additional_repairs,omitempty"`
	PartsSupplies            []PartsSupply                  `json:"parts_supplies,omitempty"`
	Services                 []Service                      `json:"services,omitempty"`
}
