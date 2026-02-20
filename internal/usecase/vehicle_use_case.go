package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	domainrepo "mecanica_xpto/internal/usecase/interfaces"
)

var (
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrInvalidPlateFormat   = errors.New("invalid plate format")
	ErrVehicleAlreadyExists = errors.New("vehicle already exists")
	ErrInvalidID            = errors.New("invalid vehicle ID")
)

var (
	MessageVehicleUpdatedSuccessfully = "Vehicle updated successfully"
	MessageVehicleNotFound            = "Vehicle not found"
	MessageInvalidPlateFormat         = "Invalid plate format"
	MessageErrorUpdatingVehicle       = "Error updating the vehicle"
	MessageErrorSearch                = "Error searching existing vehicles"
)

type VehicleServiceInterface interface {
	GetAllVehicles() ([]entities.Vehicle, error)
	GetVehicleByID(id uint) (*entities.Vehicle, error)
	GetVehicleByPlate(plate string) (*entities.Vehicle, error)
	GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error)
	CreateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error)
	UpdateVehicle(vehicle entities.Vehicle) (string, error)
	UpdateVehiclePartial(id uint, updates map[string]interface{}) (string, error)
	DeleteVehicle(id uint) error
}

type VehicleService struct {
	repo domainrepo.VehicleRepository
}

func NewVehicleService(repo domainrepo.VehicleRepository) VehicleServiceInterface {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetAllVehicles() ([]entities.Vehicle, error) {
	vehicles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}
func (s *VehicleService) GetVehicleByID(id uint) (*entities.Vehicle, error) {
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if vehicle == nil || vehicle.ID == 0 {
		return nil, ErrVehicleNotFound
	}
	return vehicle, nil
}
func (s *VehicleService) GetVehicleByPlate(plate string) (*entities.Vehicle, error) {
	voPlate := valueobject.ParsePlate(plate)
	if !voPlate.IsValidFormat() {
		return nil, ErrInvalidPlateFormat
	}
	vehicle, err := s.repo.FindByPlate(voPlate)
	if err != nil {
		return nil, err
	}
	if vehicle == nil || vehicle.ID == 0 {
		return nil, ErrVehicleNotFound
	}
	return vehicle, nil
}
func (s *VehicleService) GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	result, err := s.repo.FindByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return []entities.Vehicle{}, nil
	}
	return result, nil
}
func (s *VehicleService) CreateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	if !vehicle.Plate.IsValidFormat() {
		return nil, ErrInvalidPlateFormat
	}

	existingVehicle, err := s.repo.FindByPlate(vehicle.Plate)
	if err != nil {
		return nil, err
	}
	if existingVehicle != nil && existingVehicle.ID != 0 {
		return nil, ErrVehicleAlreadyExists
	}
	created, err := s.repo.Create(vehicle)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *VehicleService) UpdateVehicle(vehicle entities.Vehicle) (string, error) {
	if !vehicle.Plate.IsValidFormat() {
		return MessageInvalidPlateFormat, ErrInvalidPlateFormat
	}
	existingVehicle, err := s.repo.FindByID(vehicle.ID)
	if err != nil {
		return MessageErrorSearch, err
	}
	if existingVehicle == nil || existingVehicle.ID == 0 {
		return MessageVehicleNotFound, ErrVehicleNotFound
	}
	err = s.repo.Update(vehicle)
	if err != nil {
		return MessageErrorUpdatingVehicle, err
	}
	return MessageVehicleUpdatedSuccessfully, nil
}

func (s *VehicleService) UpdateVehiclePartial(id uint, updates map[string]interface{}) (string, error) {
	// First get the existing vehicle
	existingVehicle, err := s.repo.FindByID(id)
	if err != nil {
		return "", err
	}
	if existingVehicle == nil {
		return MessageVehicleNotFound, ErrVehicleNotFound
	}

	// Update only the fields that were provided
	if plate, ok := updates["plate"].(string); ok {
		parsedPlate := valueobject.ParsePlate(plate)
		if !parsedPlate.IsValidFormat() {
			return MessageInvalidPlateFormat, ErrInvalidPlateFormat
		}
		existingVehicle.Plate = parsedPlate
	}
	if model, ok := updates["model"].(string); ok {
		existingVehicle.Model = model
	}
	if year, ok := updates["year"].(string); ok {
		existingVehicle.Year = year
	}
	if brand, ok := updates["brand"].(string); ok {
		existingVehicle.Brand = brand
	}
	if customerId, ok := updates["customer_id"].(float64); ok {
		existingVehicle.CustomerID = uint(customerId)
	}

	// Convert DTO to domain entity and update
	if err := s.repo.Update(*existingVehicle); err != nil {
		return "", err
	}

	return MessageVehicleUpdatedSuccessfully, nil
}

func (s *VehicleService) DeleteVehicle(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	vehicle, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if vehicle == nil || vehicle.ID == 0 {
		return ErrVehicleNotFound
	}

	err = s.repo.Delete(id)
	return err
}
