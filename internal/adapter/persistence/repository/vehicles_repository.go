package repository

import (
	"errors"
	domainrepo "mecanica_xpto/internal/usecase/interfaces"

	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	dto "mecanica_xpto/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) domainrepo.VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) FindAll() ([]entities.Vehicle, error) {
	var vehicleDTOs []dto.VehicleModel
	if err := r.db.Preload("Customer").Find(&vehicleDTOs).Error; err != nil {
		return nil, err
	}

	vehicles := make([]entities.Vehicle, 0, len(vehicleDTOs))
	for _, v := range vehicleDTOs {
		if domain := v.ToDomain(); domain != nil {
			vehicles = append(vehicles, *domain)
		}
	}
	return vehicles, nil
}

func (r *VehicleRepository) FindByID(id uint) (*entities.Vehicle, error) {
	var vehicleDTO dto.VehicleModel
	if err := r.db.Preload("Customer").First(&vehicleDTO, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return vehicleDTO.ToDomain(), nil
}

func (r *VehicleRepository) FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error) {
	var vehicleDTO dto.VehicleModel
	if err := r.db.Preload("Customer").Where("plate = ?", plate.String()).First(&vehicleDTO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return vehicleDTO.ToDomain(), nil
}

func (r *VehicleRepository) FindByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	var vehicleDTOs []dto.VehicleModel
	if err := r.db.Preload("Customer").Where("customer_id = ?", customerID).Find(&vehicleDTOs).Error; err != nil {
		return nil, err
	}

	vehicles := make([]entities.Vehicle, 0, len(vehicleDTOs))
	for _, v := range vehicleDTOs {
		if domain := v.ToDomain(); domain != nil {
			vehicles = append(vehicles, *domain)
		}
	}
	return vehicles, nil
}

func (r *VehicleRepository) Create(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	vehicleDTO := dto.VehicleModel{
		Plate:      vehicle.Plate.String(),
		CustomerID: vehicle.CustomerID,
		Model:      vehicle.Model,
		Year:       vehicle.Year,
		Brand:      vehicle.Brand,
	}

	if err := r.db.Create(&vehicleDTO).Error; err != nil {
		return nil, err
	}

	return vehicleDTO.ToDomain(), nil
}

func (r *VehicleRepository) Update(vehicle entities.Vehicle) error {
	updates := map[string]interface{}{}
	if vehicle.Plate.String() != "" {
		updates["plate"] = vehicle.Plate.String()
	}
	if vehicle.CustomerID != 0 {
		updates["customer_id"] = vehicle.CustomerID
	}
	if vehicle.Model != "" {
		updates["model"] = vehicle.Model
	}
	if vehicle.Year != "" {
		updates["year"] = vehicle.Year
	}
	if vehicle.Brand != "" {
		updates["brand"] = vehicle.Brand
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&dto.VehicleModel{}).Where("id = ?", vehicle.ID).Updates(updates).Error
}

func (r *VehicleRepository) Delete(id uint) error {
	// Using GORM's soft delete functionality which will automatically set the deleted_at timestamp
	if err := r.db.Model(&dto.VehicleModel{}).Where("id = ?", id).Delete(nil).Error; err != nil {
		return err
	}
	return nil
}

// Ensure VehicleRepository implements the domain contract.
var _ domainrepo.VehicleRepository = (*VehicleRepository)(nil)
