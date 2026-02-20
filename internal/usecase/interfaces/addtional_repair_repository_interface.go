package interfaces

import (
	"context"
	"mecanica_xpto/internal/domain/entities"
)

type IAdditionalRepairRepository interface {
	Create(ctx context.Context, additionalRepair entities.AdditionalRepair) (entities.AdditionalRepair, error)
	GetByID(ctx context.Context, id uint) (entities.AdditionalRepair, error)
	AddPartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error
	ReplacePartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error
	GetByServiceOrder(ctx context.Context, serviceOrderId uint) ([]entities.AdditionalRepair, error)
	UpdateStatus(ctx context.Context, id uint, status entities.AdditionalRepairStatusDTO) error
	GetPartsSupplyQuantity(ctx context.Context, partsSupplyID uint, additionalRepairID uint) (int, error)
}
