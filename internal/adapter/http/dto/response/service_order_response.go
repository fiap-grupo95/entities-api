package response

import (
	"time"

	"mecanica_xpto/internal/domain/entities"
)

type ServiceOrderResponse struct {
	ID                       uint                                   `json:"id"`
	CustomerID               uint                                   `json:"customer_id"`
	VehicleID                uint                                   `json:"vehicle_id"`
	Status                   string                                 `json:"status"`
	Estimate                 float64                                `json:"estimate,omitempty"`
	StartedExecutionDate     *time.Time                             `json:"started_execution_date,omitempty"`
	FinalExecutionDate       *time.Time                             `json:"final_execution_date,omitempty"`
	ExecutionDurationInHours float64                                `json:"execution_duration_in_hours,omitempty"`
	CreatedAt                *time.Time                             `json:"created_at,omitempty"`
	UpdatedAt                *time.Time                             `json:"updated_at,omitempty"`
	Services                 []ServiceOrderServiceResponse          `json:"services"`
	PartsSupplies            []ServiceOrderPartsSupplyResponse      `json:"parts_supplies"`
	AdditionalRepairs        []ServiceOrderAdditionalRepairResponse `json:"additional_repairs"`
	Payment                  *ServiceOrderPaymentResponse           `json:"payment,omitempty"`
}

type ServiceOrderServiceResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ServiceOrderPartsSupplyResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	QuantityReserve int     `json:"quantity_reserve"`
	QuantityTotal   int     `json:"quantity_total"`
}

type ServiceOrderAdditionalRepairResponse struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Estimate    float64 `json:"estimate"`
}

type ServiceOrderPaymentResponse struct {
	ID          uint      `json:"id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
}

func NewServiceOrderResponse(entity *entities.ServiceOrder) ServiceOrderResponse {
	if entity == nil {
		return ServiceOrderResponse{}
	}

	return ServiceOrderResponse{
		ID:                       entity.ID,
		CustomerID:               entity.CustomerID,
		VehicleID:                entity.VehicleID,
		Status:                   entity.ServiceOrderStatus.String(),
		Estimate:                 entity.Estimate,
		StartedExecutionDate:     entity.StartedExecutionDate,
		FinalExecutionDate:       entity.FinalExecutionDate,
		ExecutionDurationInHours: entity.ExecutionDurationInHours,
		CreatedAt:                entity.CreatedAt,
		UpdatedAt:                entity.UpdatedAt,
		Services:                 mapServiceResponses(entity.Services),
		PartsSupplies:            mapPartsSupplyResponses(entity.PartsSupplies),
		AdditionalRepairs:        mapAdditionalRepairResponses(entity.AdditionalRepairs),
		Payment:                  mapPaymentResponse(entity.Payment),
	}
}

func NewServiceOrderListResponse(orders []*entities.ServiceOrder) []ServiceOrderResponse {
	if len(orders) == 0 {
		return []ServiceOrderResponse{}
	}

	result := make([]ServiceOrderResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, NewServiceOrderResponse(order))
	}
	return result
}

func mapServiceResponses(services []entities.Service) []ServiceOrderServiceResponse {
	if len(services) == 0 {
		return []ServiceOrderServiceResponse{}
	}

	result := make([]ServiceOrderServiceResponse, 0, len(services))
	for _, service := range services {
		result = append(result, ServiceOrderServiceResponse{
			ID:    service.ID,
			Name:  service.Name,
			Price: service.Price,
		})
	}
	return result
}

func mapPartsSupplyResponses(partsSupplies []entities.PartsSupply) []ServiceOrderPartsSupplyResponse {
	if len(partsSupplies) == 0 {
		return []ServiceOrderPartsSupplyResponse{}
	}

	result := make([]ServiceOrderPartsSupplyResponse, 0, len(partsSupplies))
	for _, ps := range partsSupplies {
		result = append(result, ServiceOrderPartsSupplyResponse{
			ID:              ps.ID,
			Name:            ps.Name,
			Price:           ps.Price,
			QuantityReserve: ps.QuantityReserve,
			QuantityTotal:   ps.QuantityTotal,
		})
	}
	return result
}

func mapAdditionalRepairResponses(repairs []entities.AdditionalRepair) []ServiceOrderAdditionalRepairResponse {
	if len(repairs) == 0 {
		return []ServiceOrderAdditionalRepairResponse{}
	}

	result := make([]ServiceOrderAdditionalRepairResponse, 0, len(repairs))
	for _, repair := range repairs {
		result = append(result, ServiceOrderAdditionalRepairResponse{
			ID:          repair.ID,
			Description: repair.Description,
			Status:      repair.ARStatus.String(),
			Estimate:    repair.Estimate,
		})
	}
	return result
}

func mapPaymentResponse(payment *entities.Payment) *ServiceOrderPaymentResponse {
	if payment == nil {
		return nil
	}

	return &ServiceOrderPaymentResponse{
		ID:          payment.ID,
		Amount:      payment.Amount,
		PaymentDate: payment.PaymentDate,
	}
}
