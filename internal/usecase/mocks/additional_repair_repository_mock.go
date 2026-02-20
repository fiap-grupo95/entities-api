// Code generated manually to match updated repository interface.
// Package mocks provides gomock implementations for use case dependencies.
package mocks

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"mecanica_xpto/internal/domain/entities"
)

// MockIAdditionalRepairRepository is a mock of IAdditionalRepairRepository interface.
type MockIAdditionalRepairRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAdditionalRepairRepositoryMockRecorder
}

// MockIAdditionalRepairRepositoryMockRecorder records invocations on MockIAdditionalRepairRepository.
type MockIAdditionalRepairRepositoryMockRecorder struct {
	mock *MockIAdditionalRepairRepository
}

// NewMockIAdditionalRepairRepository creates a new mock instance.
func NewMockIAdditionalRepairRepository(ctrl *gomock.Controller) *MockIAdditionalRepairRepository {
	mock := &MockIAdditionalRepairRepository{ctrl: ctrl}
	mock.recorder = &MockIAdditionalRepairRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns recorder for expectations.
func (m *MockIAdditionalRepairRepository) EXPECT() *MockIAdditionalRepairRepositoryMockRecorder {
	return m.recorder
}

// Create mocks the Create method.
func (m *MockIAdditionalRepairRepository) Create(ctx context.Context, additionalRepair entities.AdditionalRepair) (entities.AdditionalRepair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, additionalRepair)
	ret0, _ := ret[0].(entities.AdditionalRepair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) Create(ctx, additionalRepair interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).Create), ctx, additionalRepair)
}

// GetByID mocks the GetByID method.
func (m *MockIAdditionalRepairRepository) GetByID(ctx context.Context, id uint) (entities.AdditionalRepair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(entities.AdditionalRepair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).GetByID), ctx, id)
}

// AddPartSupplyAndService mocks the AddPartSupplyAndService method.
func (m *MockIAdditionalRepairRepository) AddPartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPartSupplyAndService", ctx, additionalRepairID, services, partsSupplies, newEstimate)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPartSupplyAndService indicates an expected call of AddPartSupplyAndService.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) AddPartSupplyAndService(ctx, additionalRepairID, services, partsSupplies, newEstimate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPartSupplyAndService", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).AddPartSupplyAndService), ctx, additionalRepairID, services, partsSupplies, newEstimate)
}

// ReplacePartSupplyAndService mocks the ReplacePartSupplyAndService method.
func (m *MockIAdditionalRepairRepository) ReplacePartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplacePartSupplyAndService", ctx, additionalRepairID, services, partsSupplies, newEstimate)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplacePartSupplyAndService indicates an expected call of ReplacePartSupplyAndService.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) ReplacePartSupplyAndService(ctx, additionalRepairID, services, partsSupplies, newEstimate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplacePartSupplyAndService", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).ReplacePartSupplyAndService), ctx, additionalRepairID, services, partsSupplies, newEstimate)
}

// GetByServiceOrder mocks the GetByServiceOrder method.
func (m *MockIAdditionalRepairRepository) GetByServiceOrder(ctx context.Context, serviceOrderId uint) ([]entities.AdditionalRepair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByServiceOrder", ctx, serviceOrderId)
	ret0, _ := ret[0].([]entities.AdditionalRepair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByServiceOrder indicates an expected call of GetByServiceOrder.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) GetByServiceOrder(ctx, serviceOrderId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByServiceOrder", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).GetByServiceOrder), ctx, serviceOrderId)
}

// UpdateStatus mocks the UpdateStatus method.
func (m *MockIAdditionalRepairRepository) UpdateStatus(ctx context.Context, id uint, status entities.AdditionalRepairStatusDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) UpdateStatus(ctx, id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).UpdateStatus), ctx, id, status)
}

// GetPartsSupplyQuantity mocks the GetPartsSupplyQuantity method.
func (m *MockIAdditionalRepairRepository) GetPartsSupplyQuantity(ctx context.Context, partsSupplyID uint, additionalRepairID uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPartsSupplyQuantity", ctx, partsSupplyID, additionalRepairID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPartsSupplyQuantity indicates an expected call of GetPartsSupplyQuantity.
func (mr *MockIAdditionalRepairRepositoryMockRecorder) GetPartsSupplyQuantity(ctx, partsSupplyID, additionalRepairID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPartsSupplyQuantity", reflect.TypeOf((*MockIAdditionalRepairRepository)(nil).GetPartsSupplyQuantity), ctx, partsSupplyID, additionalRepairID)
}
