package interfaces

import "mecanica_xpto/internal/domain/entities"

type IServiceOrderRepository interface {
	Create(serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error)
	GetByID(id uint) (*entities.ServiceOrder, error)
	Update(serviceOrder *entities.ServiceOrder) error
	List() ([]*entities.ServiceOrder, error)
	UpdateEstimate(id uint, estimate float64) error
	GetPartsSupplyServiceOrder(partsSupplyID uint, serviceOrderID uint) (*entities.ServiceOrderPartsSupply, error)
}
