//go:build legacy
// +build legacy

package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	dto "mecanica_xpto/internal/infrastructure/database/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func toServiceOrderEntity(model *dto.ServiceOrderModel) *entities.ServiceOrder {
	if model == nil {
		return nil
	}
	return model.ToDomain()
}

const (
	StatusRecebida            = string(valueobject.StatusRecebida)
	StatusEmDiagnostico       = string(valueobject.StatusEmDiagnostico)
	StatusAguardandoAprovacao = string(valueobject.StatusAguardandoAprovacao)
	StatusAprovada            = string(valueobject.StatusAprovada)
	StatusRejeitada           = string(valueobject.StatusRejeitada)
	StatusEmExecucao          = string(valueobject.StatusEmExecucao)
	StatusFinalizada          = string(valueobject.StatusFinalizada)
	StatusEntregue            = string(valueobject.StatusEntregue)
	StatusCancelada           = string(valueobject.StatusCancelada)
)

type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) Update(ctx context.Context, so *entities.Service) error {
	args := m.Called(ctx, so)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockServiceRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockServiceRepository) List(ctx context.Context) ([]entities.Service, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceRepository) Create(ctx context.Context, so *entities.Service) (entities.Service, error) {
	args := m.Called(ctx, so)
	if args.Get(0) == nil {
		return entities.Service{}, args.Error(1)
	}
	return args.Get(0).(entities.Service), args.Error(1)
}
func (m *MockServiceRepository) GetByID(ctx context.Context, id uint) (entities.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return entities.Service{}, args.Error(1)
	}
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceRepository) GetByName(ctx context.Context, name string) (entities.Service, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return entities.Service{}, args.Error(1)
	}
	return args.Get(0).(entities.Service), args.Error(1)
}

// Mock Vehicle Repository - "github.com/stretchr/testify/mock"
type MockVehicleRepository struct {
	mock.Mock
}

func (m *MockVehicleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVehicleRepository) FindAll() ([]entities.Vehicle, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	switch v := args.Get(0).(type) {
	case []entities.Vehicle:
		return v, args.Error(1)
	case []dto.VehicleModel:
		result := make([]entities.Vehicle, 0, len(v))
		for _, item := range v {
			if domain := item.ToDomain(); domain != nil {
				result = append(result, *domain)
			}
		}
		return result, args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockVehicleRepository) FindByID(id uint) (*entities.Vehicle, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	switch v := args.Get(0).(type) {
	case *entities.Vehicle:
		return v, args.Error(1)
	case entities.Vehicle:
		return &v, args.Error(1)
	case *dto.VehicleModel:
		if v == nil {
			return nil, args.Error(1)
		}
		return v.ToDomain(), args.Error(1)
	case dto.VehicleModel:
		copy := v
		return copy.ToDomain(), args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockVehicleRepository) FindByPlate(plate valueobject.Plate) (*entities.Vehicle, error) {
	args := m.Called(plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	switch v := args.Get(0).(type) {
	case *entities.Vehicle:
		return v, args.Error(1)
	case entities.Vehicle:
		return &v, args.Error(1)
	case *dto.VehicleModel:
		if v == nil {
			return nil, args.Error(1)
		}
		return v.ToDomain(), args.Error(1)
	case dto.VehicleModel:
		copy := v
		return copy.ToDomain(), args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockVehicleRepository) FindByCustomerID(customerID uint) ([]entities.Vehicle, error) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	switch v := args.Get(0).(type) {
	case []entities.Vehicle:
		return v, args.Error(1)
	case []dto.VehicleModel:
		result := make([]entities.Vehicle, 0, len(v))
		for _, item := range v {
			if domain := item.ToDomain(); domain != nil {
				result = append(result, *domain)
			}
		}
		return result, args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockVehicleRepository) Create(vehicle entities.Vehicle) (*entities.Vehicle, error) {
	args := m.Called(vehicle)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	switch v := args.Get(0).(type) {
	case *entities.Vehicle:
		return v, args.Error(1)
	case entities.Vehicle:
		return &v, args.Error(1)
	case *dto.VehicleModel:
		if v == nil {
			return nil, args.Error(1)
		}
		return v.ToDomain(), args.Error(1)
	case dto.VehicleModel:
		copy := v
		return copy.ToDomain(), args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

func (m *MockVehicleRepository) Update(vehicle entities.Vehicle) error {
	args := m.Called(vehicle)
	return args.Error(0)
}

// Mock Customer Repository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetByID(id uint) (*entities.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CustomerModel).ToDomain(), args.Error(1)
}

func (m *MockCustomerRepository) GetByDocument(CpfCnpj string) (*entities.Customer, error) {
	args := m.Called(CpfCnpj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CustomerModel).ToDomain(), args.Error(1)
}
func (m *MockCustomerRepository) Create(customer *entities.Customer) error {
	args := m.Called(customer)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}
func (m *MockCustomerRepository) Update(customer *entities.Customer) error {
	customerDto := dto.FromDomainCustomer(customer)
	args := m.Called(customerDto)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(id uint) error {
	args := m.Called(id)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockCustomerRepository) List() ([]entities.Customer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Customer), args.Error(1)
}

// Mock Service Order Repository
type MockServiceOrderRepository struct {
	mock.Mock
}

func (m *MockServiceOrderRepository) UpdateEstimate(id uint, estimate float64) error {
	args := m.Called(id, estimate)
	return args.Error(0)
}

func (m *MockServiceOrderRepository) GetByName(ctx context.Context, name string) (entities.Service, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return entities.Service{}, args.Error(1)
	}
	return args.Get(0).(entities.Service), args.Error(1)
}
func (m *MockServiceOrderRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockServiceOrderRepository) Create(serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	args := m.Called(serviceOrder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.ServiceOrder), args.Error(1)
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

func (m *MockServiceOrderRepository) Update(serviceOrder *entities.ServiceOrder) error {
	args := m.Called(serviceOrder)
	return args.Error(0)
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
	switch value := args.Get(0).(type) {
	case *entities.ServiceOrderPartsSupply:
		return value, args.Error(1)
	case entities.ServiceOrderPartsSupply:
		return &value, args.Error(1)
	case *dto.PartsSupplyServiceOrder:
		return &entities.ServiceOrderPartsSupply{
			PartsSupplyID:  value.PartsSupplyID,
			ServiceOrderID: value.ServiceOrderID,
			Quantity:       value.Quantity,
		}, args.Error(1)
	case dto.PartsSupplyServiceOrder:
		return &entities.ServiceOrderPartsSupply{
			PartsSupplyID:  value.PartsSupplyID,
			ServiceOrderID: value.ServiceOrderID,
			Quantity:       value.Quantity,
		}, args.Error(1)
	default:
		return nil, args.Error(1)
	}
}

// Mock Parts Supply Repository
type MockPartsSupplyRepository struct {
	mock.Mock
}

func (m *MockPartsSupplyRepository) GetByName(ctx context.Context, name string) (entities.PartsSupply, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return entities.PartsSupply{}, args.Error(1)
	}
	return args.Get(0).(entities.PartsSupply), args.Error(1)
}

func (m *MockPartsSupplyRepository) Create(ctx context.Context, ps *entities.PartsSupply) (entities.PartsSupply, error) {
	args := m.Called(ctx, ps)
	if args.Get(0) == nil {
		return entities.PartsSupply{}, args.Error(1)
	}
	return args.Get(0).(entities.PartsSupply), args.Error(1)
}

func (m *MockPartsSupplyRepository) GetByID(ctx context.Context, id uint) (entities.PartsSupply, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return entities.PartsSupply{}, args.Error(1)
	}
	return args.Get(0).(entities.PartsSupply), args.Error(1)
}

func (m *MockPartsSupplyRepository) Update(ctx context.Context, partsSupply *entities.PartsSupply) error {
	args := m.Called(ctx, partsSupply)
	return args.Error(0)
}

func (m *MockPartsSupplyRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPartsSupplyRepository) List(ctx context.Context) ([]entities.PartsSupply, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.PartsSupply), args.Error(1)
}

func (m *MockPartsSupplyRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error) {
	args := m.Called(ctx, serviceOrderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.PartsSupply), args.Error(1)
}

func TestCreateServiceOrder(t *testing.T) {
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceOrderRepo := new(MockServiceOrderRepository)
	serviceRepo := new(MockServiceRepository)
	partsSupplyRepo := new(MockPartsSupplyRepository)

	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success - Valid service order creation",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 1,
				VehicleID:  1,
			},
			setupMocks: func() {
				vehicleRepo.On("FindByID", uint(1)).Return(&dto.VehicleModel{ID: 1}, nil)
				customerRepo.On("GetByID", uint(1)).Return(&dto.CustomerModel{ID: 1}, nil)
				serviceOrderRepo.On("Create", mock.AnythingOfType("*entities.ServiceOrder")).Return(&entities.ServiceOrder{
					ID:         1,
					CustomerID: 1,
					VehicleID:  1,
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "Error - Vehicle not found",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 1,
				VehicleID:  1,
			},
			setupMocks: func() {
				vehicleRepo.On("FindByID", uint(1)).Return(nil, errors.New("vehicle not found"))
			},
			expectedError: errors.New("vehicle not found"),
		},
		{
			name: "Error - Customer not found",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 1,
				VehicleID:  1,
			},
			setupMocks: func() {
				vehicleRepo.On("FindByID", mock.Anything, uint(1)).Return(&dto.VehicleModel{ID: 1}, nil)
				customerRepo.On("GetByID", uint(1)).Return(nil, errors.New("customer not found"))
			},
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			r, err := useCase.CreateServiceOrder(context.Background(), tt.serviceOrder)
			if tt.expectedError != nil && err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NotNil(t, r)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateServiceOrder(t *testing.T) {
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceOrderRepo := new(MockServiceOrderRepository)
	serviceRepo := new(MockServiceRepository)
	partsSupplyRepo := new(MockPartsSupplyRepository)

	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	setupMocks := func() {
		serviceOrderRepo.On("GetByID", uint(1)).Return(&dto.ServiceOrderModel{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatus{
				ID:          1,
				Description: string(valueobject.StatusRecebida),
			},
		}, nil)
		serviceRepo.On("GetByID", mock.Anything, uint(1)).Return(entities.Service{ID: 1}, nil)
		partsSupplyRepo.On("GetByID", mock.Anything, uint(1)).Return(entities.PartsSupply{
			ID:              1,
			QuantityTotal:   10,
			QuantityReserve: 2,
		}, nil)
		partsSupplyRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.PartsSupply")).Return(nil)
		serviceOrderRepo.On("Update", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)
	}
	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		flow          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success - Update to EmDiagnostico",
			serviceOrder: entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEmDiagnostico,
				Services: []entities.Service{
					{ID: 1},
				},
				PartsSupplies: []entities.PartsSupply{
					{
						ID:              1,
						QuantityReserve: 2,
					},
				},
			},
			flow:          DIAGNOSIS,
			setupMocks:    setupMocks,
			expectedError: nil,
		},
		{
			name: "Error - Invalid Status Transition",
			serviceOrder: entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
			},
			flow: DIAGNOSIS,
			setupMocks: func() {
				serviceOrderRepo.On("GetByID", uint(1)).Return(&dto.ServiceOrderModel{
					ID: 1,
					ServiceOrderStatus: dto.ServiceOrderStatus{
						ID:          1,
						Description: string(valueobject.StatusRecebida),
					},
				}, nil)
			},
			expectedError: ErrInvalidTransitionStatusToDiagnosis,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			r, err := useCase.UpdateServiceOrder(context.Background(), tt.serviceOrder, tt.flow)
			if tt.expectedError != nil && err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NotNil(t, r)
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEstimate(t *testing.T) {
	tests := []struct {
		name            string
		request         *entities.ServiceOrder
		serviceOrderDTO *dto.ServiceOrderModel
		setupMocks      func(*MockPartsSupplyRepository, *MockServiceOrderRepository)
		expectedError   error
	}{
		{
			name: "Should approve estimate and release parts supply",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusAprovada,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: StatusAguardandoAprovacao,
				},
			},
			setupMocks: func(psRepo *MockPartsSupplyRepository, soRepo *MockServiceOrderRepository) {
				// Mock get parts supplies by service order ID
				psRepo.On("GetByServiceOrderID", context.Background(), uint(1)).Return([]entities.PartsSupply{
					{ID: 1, QuantityTotal: 10, QuantityReserve: 2},
				}, nil)

				// Mock get parts supply service order relation
				soRepo.On("GetPartsSupplyServiceOrder", uint(1), uint(1)).Return(&dto.PartsSupplyServiceOrder{
					PartsSupplyID:  1,
					ServiceOrderID: 1,
					Quantity:       2,
				}, nil)

				// Mock get parts supply by ID
				psRepo.On("GetByID", context.Background(), uint(1)).Return(entities.PartsSupply{
					ID:              1,
					QuantityTotal:   10,
					QuantityReserve: 2,
				}, nil)

				// Mock update parts supply
				psRepo.On("Update", context.Background(), mock.AnythingOfType("*entities.PartsSupply")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Should fail when getting parts supply relation fails",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusAprovada,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: StatusAguardandoAprovacao,
				},
			},
			setupMocks: func(psRepo *MockPartsSupplyRepository, soRepo *MockServiceOrderRepository) {
				psRepo.On("GetByServiceOrderID", context.Background(), uint(1)).Return([]entities.PartsSupply{
					{ID: 1},
				}, nil)

				soRepo.On("GetPartsSupplyServiceOrder", uint(1), uint(1)).Return(nil, errors.New("error getting relation"))
			},
			expectedError: errors.New("error getting relation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			partsSupplyRepo := new(MockPartsSupplyRepository)
			serviceOrderRepo := new(MockServiceOrderRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(partsSupplyRepo, serviceOrderRepo)
			}

			update := &entities.ServiceOrder{}
			current := toServiceOrderEntity(tt.serviceOrderDTO)
			result, err := ValidateEstimate(context.Background(), tt.request, current, update, partsSupplyRepo, serviceOrderRepo)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, valueobject.StatusAprovada, result.ServiceOrderStatus)
			}
		})
	}
}

func TestCalculateEstimate(t *testing.T) {
	tests := []struct {
		name           string
		services       []entities.Service
		partsSupplies  []entities.PartsSupply
		serviceRepo    *MockServiceRepository
		partsSuplyRepo *MockPartsSupplyRepository
		expected       float64
	}{
		{
			name: "Calculate with services and parts supplies",
			services: []entities.Service{
				{ID: 1},
				{ID: 2},
			},
			partsSupplies: []entities.PartsSupply{
				{ID: 1, QuantityReserve: 2},
				{ID: 2, QuantityReserve: 3},
			},
			serviceRepo:    &MockServiceRepository{},
			partsSuplyRepo: &MockPartsSupplyRepository{},
			expected:       350.0, // Updated to match actual calculation: (100 + 75) + (50*2 + 25*3) = 175 + 175 = 350
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks for services
			tt.serviceRepo.On("GetByID", mock.Anything, uint(1)).Return(entities.Service{
				ID:    1,
				Price: 100.0,
			}, nil)
			tt.serviceRepo.On("GetByID", mock.Anything, uint(2)).Return(entities.Service{
				ID:    2,
				Price: 75.0,
			}, nil)

			// Setup mocks for parts supplies
			tt.partsSuplyRepo.On("GetByID", mock.Anything, uint(1)).Return(entities.PartsSupply{
				ID:    1,
				Price: 50.0,
			}, nil)
			tt.partsSuplyRepo.On("GetByID", mock.Anything, uint(2)).Return(entities.PartsSupply{
				ID:    2,
				Price: 25.0,
			}, nil)

			result, err := CalculateEstimate(context.Background(), tt.services, tt.partsSupplies, tt.serviceRepo, tt.partsSuplyRepo)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateExecution(t *testing.T) {
	tests := []struct {
		name             string
		request          *entities.ServiceOrder
		serviceOrderDTO  *dto.ServiceOrderModel
		expectedError    error
		expectedDuration float64 // Espera-se que a duraÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â§ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â£o da execuÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â§ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â£o seja arredondada para 2 casas decimais
	}{
		{
			name: "Success - Start Execution",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEmExecucao,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: string(valueobject.StatusAprovada),
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - Finish Execution",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusFinalizada,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: string(valueobject.StatusEmExecucao),
				},
				StartedExecutionDate: func() *time.Time {
					// Cria uma data que ocorreu 1.25 horas no passado
					t := time.Now().Add(-75 * time.Minute)
					return &t
				}(),
			},
			expectedDuration: 1.25, // Exemplo de duraÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â§ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â£o esperada com 2 casas decimais
			expectedError:    nil,
		},
		{
			name: "Error - Invalid Status Transition",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusFinalizada,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: string(valueobject.StatusRecebida),
				},
			},
			expectedError: ErrInvalidTransitionStatusToExecution,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &entities.ServiceOrder{}
			current := toServiceOrderEntity(tt.serviceOrderDTO)
			result, err := ValidateExecution(context.Background(), tt.request, current, update)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.ServiceOrderStatus, result.ServiceOrderStatus)
				assert.Equal(t, tt.expectedDuration, result.ExecutionDurationInHours)
			}
		})
	}
}

func TestValidateDelivery(t *testing.T) {
	tests := []struct {
		name            string
		request         *entities.ServiceOrder
		serviceOrderDTO *dto.ServiceOrderModel
		expectedError   error
	}{
		{
			name: "Success - Complete Delivery",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: string(valueobject.StatusFinalizada),
				},
				Payment: &dto.PaymentModel{
					ID:           1,
					ServiceOrder: dto.ServiceOrderModel{ID: 1},
					PaymentDate:  time.Now(),
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Missing Payment Information",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: string(valueobject.StatusFinalizada),
				},
			},
			expectedError: errors.New("payment information is required for delivery"),
		},
		{
			name: "Error - Invalid Status Transition",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
				Payment: &entities.Payment{
					ID:           1,
					ServiceOrder: &entities.ServiceOrder{ID: 1},
					PaymentDate:  time.Now(),
				},
			},
			serviceOrderDTO: &dto.ServiceOrderModel{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatus{
					Description: string(valueobject.StatusEmDiagnostico),
				},
			},
			expectedError: ErrInvalidTransitionStatusToDelivery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &entities.ServiceOrder{}
			current := toServiceOrderEntity(tt.serviceOrderDTO)
			result, err := ValidateDelivery(context.Background(), tt.request, current, update)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, valueobject.StatusEntregue, result.ServiceOrderStatus)
			}
		})
	}
}

func TestInvalidServiceOrder(t *testing.T) {
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceOrderRepo := new(MockServiceOrderRepository)
	serviceRepo := new(MockServiceRepository)
	partsSupplyRepo := new(MockPartsSupplyRepository)

	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		flow          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Error - Service Order Not Found",
			serviceOrder: entities.ServiceOrder{
				ID:                 999,
				ServiceOrderStatus: valueobject.StatusEmDiagnostico,
			},
			flow: DIAGNOSIS,
			setupMocks: func() {
				serviceOrderRepo.On("GetByID", uint(999)).Return(nil, ErrServiceOrderNotFound)
			},
			expectedError: ErrServiceOrderNotFound,
		},
		{
			name: "Error - Invalid Flow Type",
			serviceOrder: entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEmDiagnostico,
			},
			flow: "invalid_flow",
			setupMocks: func() {
				serviceOrderRepo.On("GetByID", uint(1)).Return(&dto.ServiceOrderModel{
					ID: 1,
					ServiceOrderStatus: dto.ServiceOrderStatus{
						Description: string(valueobject.StatusRecebida),
					},
				}, nil)
			},
			expectedError: ErrInvalidFlow, // The update will return nil since no valid flow was matched
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			r, err := useCase.UpdateServiceOrder(context.Background(), tt.serviceOrder, tt.flow)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NotNil(t, r)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetServiceOrder(t *testing.T) {
	serviceOrderRepo := new(MockServiceOrderRepository)
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceRepo := new(MockServiceRepository)
	partsSupplyRepo := new(MockPartsSupplyRepository)
	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	ctx := context.Background()
	validID := uint(1)
	invalidID := uint(999)
	serviceOrderEntity := entities.ServiceOrder{ID: validID}
	serviceOrderDTO := &dto.ServiceOrderModel{ID: validID}

	t.Run("success", func(t *testing.T) {
		serviceOrderRepo.On("GetByID", validID).Return(serviceOrderDTO, nil)
		result, err := useCase.GetServiceOrder(ctx, serviceOrderEntity)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		serviceOrderRepo.AssertCalled(t, "GetByID", validID)
	})

	t.Run("not found", func(t *testing.T) {
		serviceOrderRepo.On("GetByID", invalidID).Return(nil, nil)
		result, err := useCase.GetServiceOrder(ctx, entities.ServiceOrder{ID: invalidID})
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrServiceOrderNotFound, err)
		serviceOrderRepo.AssertCalled(t, "GetByID", invalidID)
	})
}

func TestListServiceOrders(t *testing.T) {
	serviceOrderRepo := new(MockServiceOrderRepository)
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceRepo := new(MockServiceRepository)
	partsSupplyRepo := new(MockPartsSupplyRepository)
	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	ctx := context.Background()
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	timePtr := func(offsetDays int) *time.Time {
		ts := base.Add(time.Duration(offsetDays) * 24 * time.Hour)
		return &ts
	}

	serviceOrderDTOs := []dto.ServiceOrderModel{
		{ID: 1, CreatedAt: timePtr(0), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmExecucao)}},
		{ID: 5, CreatedAt: timePtr(2), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmExecucao)}},
		{ID: 2, CreatedAt: timePtr(1), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusAguardandoAprovacao)}},
		{ID: 3, CreatedAt: timePtr(3), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEmDiagnostico)}},
		{ID: 4, CreatedAt: timePtr(4), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusRecebida)}},
		{ID: 6, CreatedAt: timePtr(5), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusCancelada)}},
		{ID: 7, CreatedAt: timePtr(6), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusFinalizada)}},
		{ID: 8, CreatedAt: timePtr(7), ServiceOrderStatus: dto.ServiceOrderStatus{Description: string(valueobject.StatusEntregue)}},
	}

	t.Run("apply ordering rules", func(t *testing.T) {
		serviceOrderRepo.On("List").Return(serviceOrderDTOs, nil)
		result, err := useCase.ListServiceOrders(ctx)

		assert.NoError(t, err)
		// Finalizada and Entregue must be excluded from the listing.
		assert.Len(t, result, 6)

		expectedOrder := []uint{1, 5, 2, 3, 4, 6}
		for idx, expectedID := range expectedOrder {
			assert.Equal(t, expectedID, result[idx].ID)
		}

		// Ensure no filtered statuses leak through.
		for _, so := range result {
			assert.False(t, so.ServiceOrderStatus.IsFinalizada())
			assert.False(t, so.ServiceOrderStatus.IsEntregue())
		}

		serviceOrderRepo.AssertExpectations(t)
	})
}
