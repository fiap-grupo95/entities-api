package mocks

import (
	"mecanica_xpto/internal/domain/entities"
	dto "mecanica_xpto/internal/infrastructure/database/model"

	"github.com/stretchr/testify/mock"
)

// Mock ServiceOrder Repository
type MockServiceOrderRepository struct {
	mock.Mock
}

func (m *MockServiceOrderRepository) Create(serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	args := m.Called(serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
}

func (m *MockServiceOrderRepository) Update(serviceOrder *entities.ServiceOrder) error {
	args := m.Called(serviceOrder)
	return args.Error(0)
}

func (m *MockServiceOrderRepository) GetByID(id uint) (*entities.ServiceOrder, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	switch value := args.Get(0).(type) {
	case *entities.ServiceOrder:
		return value, args.Error(1)
	case entities.ServiceOrder:
		return &value, args.Error(1)
	case *dto.ServiceOrderModel:
		return value.ToDomain(), args.Error(1)
	case dto.ServiceOrderModel:
		copy := value
		return copy.ToDomain(), args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockServiceOrderRepository) List() ([]*entities.ServiceOrder, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	switch value := args.Get(0).(type) {
	case []*entities.ServiceOrder:
		return value, args.Error(1)
	case []entities.ServiceOrder:
		result := make([]*entities.ServiceOrder, 0, len(value))
		for _, item := range value {
			itemCopy := item
			result = append(result, &itemCopy)
		}
		return result, args.Error(1)
	case []dto.ServiceOrderModel:
		result := make([]*entities.ServiceOrder, 0, len(value))
		for _, item := range value {
			itemCopy := item
			result = append(result, itemCopy.ToDomain())
		}
		return result, args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockServiceOrderRepository) GetPartsSupplyServiceOrder(partsSupplyID uint, serviceOrderID uint) (*entities.ServiceOrderPartsSupply, error) {
	args := m.Called(partsSupplyID, serviceOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrderPartsSupply), args.Error(1)
}

func (m *MockServiceOrderRepository) UpdateEstimate(id uint, estimate float64) error {
	args := m.Called(id, estimate)
	return args.Error(0)
}
