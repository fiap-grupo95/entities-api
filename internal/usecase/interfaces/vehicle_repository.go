package interfaces

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
)

// VehicleRepository defines the contracts required by the domain/use case layer
// to persist and retrieve vehicle aggregates. Implementations live in the
// adapter layer.
type VehicleRepository interface {
	FindAll() ([]entities.Vehicle, error)
	FindByID(id uint) (*entities.Vehicle, error)
	FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error)
	FindByCustomerID(customerID uint) ([]entities.Vehicle, error)
	Create(vehicle entities.Vehicle) (*entities.Vehicle, error)
	Update(vehicle entities.Vehicle) error
	Delete(id uint) error
}
