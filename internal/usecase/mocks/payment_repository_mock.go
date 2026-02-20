package mocks

import (
	"context"
	"mecanica_xpto/internal/domain/entities"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockIPaymentRepo provides a gomock mock for the payment repository.
type MockIPaymentRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIPaymentRepoMockRecorder
}

// MockIPaymentRepoMockRecorder records calls on MockIPaymentRepo.
type MockIPaymentRepoMockRecorder struct {
	mock *MockIPaymentRepo
}

// NewMockIPaymentRepo creates a new mock instance.
func NewMockIPaymentRepo(ctrl *gomock.Controller) *MockIPaymentRepo {
	mock := &MockIPaymentRepo{ctrl: ctrl}
	mock.recorder = &MockIPaymentRepoMockRecorder{mock}
	return mock
}

// EXPECT exposes recorder.
func (m *MockIPaymentRepo) EXPECT() *MockIPaymentRepoMockRecorder {
	return m.recorder
}

// Create mocks Create.
func (m *MockIPaymentRepo) Create(ctx context.Context, payment entities.Payment) (entities.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, payment)
	ret0, _ := ret[0].(entities.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create expectation.
func (mr *MockIPaymentRepoMockRecorder) Create(ctx, payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIPaymentRepo)(nil).Create), ctx, payment)
}

// GetByID mocks GetByID.
func (m *MockIPaymentRepo) GetByID(ctx context.Context, id uint) (entities.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(entities.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID expectation.
func (mr *MockIPaymentRepoMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockIPaymentRepo)(nil).GetByID), ctx, id)
}

// GetByServiceOrderID mocks GetByServiceOrderID.
func (m *MockIPaymentRepo) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) (entities.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByServiceOrderID", ctx, serviceOrderID)
	ret0, _ := ret[0].(entities.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByServiceOrderID expectation.
func (mr *MockIPaymentRepoMockRecorder) GetByServiceOrderID(ctx, serviceOrderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByServiceOrderID", reflect.TypeOf((*MockIPaymentRepo)(nil).GetByServiceOrderID), ctx, serviceOrderID)
}

// List mocks List.
func (m *MockIPaymentRepo) List(ctx context.Context) ([]entities.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]entities.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List expectation.
func (mr *MockIPaymentRepoMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIPaymentRepo)(nil).List), ctx)
}
