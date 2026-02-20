package mocks

import (
	"mecanica_xpto/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockVehicleService struct {
	mock.Mock
}

func (m *MockVehicleService) GetAllVehicles() ([]entities.Vehicle, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Vehicle), args.Error(1)
}

func (m *MockVehicleService) GetVehicleByID(id uint) (*entities.Vehicle, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleService) GetVehiclesByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Vehicle), args.Error(1)
}

func (m *MockVehicleService) CreateVehicle(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	args := m.Called(vehicle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Vehicle), args.Error(1)
}

func (m *MockVehicleService) UpdateVehicle(vehicle entities.Vehicle) (string, error) {
	args := m.Called(vehicle)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}

func (m *MockVehicleService) DeleteVehicle(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVehicleService) GetVehicleByPlate(plate string) (*entities.Vehicle, error) {
	args := m.Called(plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Vehicle), args.Error(1)
}
func (m *MockVehicleService) UpdateVehiclePartial(id uint, updates map[string]interface{}) (string, error) {
	args := m.Called(id, updates)
	return args.String(0), args.Error(1)
}
