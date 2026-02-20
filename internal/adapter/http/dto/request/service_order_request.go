package request

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
)

type ServiceOrderCreateRequest struct {
	CustomerID uint `json:"customer_id" binding:"required"`
	VehicleID  uint `json:"vehicle_id" binding:"required"`
}

type ServiceOrderDiagnosisUpdateRequest struct {
	ServiceOrderStatus string                        `json:"service_order_status" binding:"required"`
	Services           []ServiceOrderServiceItem     `json:"services" binding:"omitempty,dive"`
	PartsSupplies      []ServiceOrderPartsSupplyItem `json:"parts_supplies" binding:"omitempty,dive"`
}

type ServiceOrderEstimateUpdateRequest struct {
	ServiceOrderStatus string                        `json:"service_order_status" binding:"required"`
	Services           []ServiceOrderServiceItem     `json:"services,omitempty"`
	PartsSupplies      []ServiceOrderPartsSupplyItem `json:"parts_supplies,omitempty"`
}

type ServiceOrderExecutionUpdateRequest struct {
	ServiceOrderStatus string `json:"service_order_status" binding:"required"`
}

type ServiceOrderDeliveryUpdateRequest struct {
	ServiceOrderStatus string `json:"service_order_status" binding:"required"`
}

type ServiceOrderServiceItem struct {
	ID uint `json:"id" binding:"required"`
}

type ServiceOrderPartsSupplyItem struct {
	ID              uint `json:"id" binding:"required"`
	QuantityReserve int  `json:"quantity_reserve,omitempty"`
	QuantityTotal   int  `json:"quantity_total,omitempty"`
}

func (r ServiceOrderCreateRequest) ToEntity() entities.ServiceOrder {
	return entities.ServiceOrder{
		CustomerID: r.CustomerID,
		VehicleID:  r.VehicleID,
	}
}

func (r ServiceOrderDiagnosisUpdateRequest) ToEntity(id uint) entities.ServiceOrder {
	return entities.ServiceOrder{
		ID:                 id,
		ServiceOrderStatus: valueobject.ParseServiceOrderStatus(r.ServiceOrderStatus),
		Services:           mapServiceItemsToEntities(r.Services),
		PartsSupplies:      mapPartsSupplyItemsToEntities(r.PartsSupplies),
	}
}

func (r ServiceOrderEstimateUpdateRequest) ToEntity(id uint) entities.ServiceOrder {
	return entities.ServiceOrder{
		ID:                 id,
		ServiceOrderStatus: valueobject.ParseServiceOrderStatus(r.ServiceOrderStatus),
		Services:           mapServiceItemsToEntities(r.Services),
		PartsSupplies:      mapPartsSupplyItemsToEntities(r.PartsSupplies),
	}
}

func (r ServiceOrderExecutionUpdateRequest) ToEntity(id uint) entities.ServiceOrder {
	return entities.ServiceOrder{
		ID:                 id,
		ServiceOrderStatus: valueobject.ParseServiceOrderStatus(r.ServiceOrderStatus),
	}
}

func (r ServiceOrderDeliveryUpdateRequest) ToEntity(id uint) entities.ServiceOrder {
	return entities.ServiceOrder{
		ID:                 id,
		ServiceOrderStatus: valueobject.ParseServiceOrderStatus(r.ServiceOrderStatus),
	}
}

func mapServiceItemsToEntities(items []ServiceOrderServiceItem) []entities.Service {
	if len(items) == 0 {
		return nil
	}

	result := make([]entities.Service, 0, len(items))
	for _, item := range items {
		result = append(result, entities.Service{ID: item.ID})
	}
	return result
}

func mapPartsSupplyItemsToEntities(items []ServiceOrderPartsSupplyItem) []entities.PartsSupply {
	if len(items) == 0 {
		return nil
	}

	result := make([]entities.PartsSupply, 0, len(items))
	for _, item := range items {
		result = append(result, entities.PartsSupply{
			ID:              item.ID,
			QuantityReserve: item.QuantityReserve,
			QuantityTotal:   item.QuantityTotal,
		})
	}
	return result
}
