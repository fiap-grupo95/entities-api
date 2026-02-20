//go:build legacy
// +build legacy

package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"mecanica_xpto/internal/domain/entities"
	domainmocks "mecanica_xpto/internal/domain/mocks"
	dto "mecanica_xpto/internal/infrastructure/database/model"
	usecasemocks "mecanica_xpto/internal/usecase/mocks"
)

func TestPaymentUseCase_CreatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockPaymentRepo := usecasemocks.NewMockIPaymentRepo(ctrl)
	mockServiceOrderRepo := &domainmocks.MockServiceOrderRepository{}
	u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

	mockServiceOrderDTO := &dto.ServiceOrderModel{ID: 1, Estimate: 100.0}
	payment := &entities.Payment{ID: 1, ServiceOrderID: 1, Amount: 100.0}

	t.Run("success", func(t *testing.T) {
		mockServiceOrderRepo.On("GetByID", uint(1)).Return(mockServiceOrderDTO, nil)
		mockPaymentRepo.EXPECT().GetByServiceOrderID(ctx, uint(1)).Return(entities.Payment{}, nil)
		mockPaymentRepo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(entities.Payment{})).Return(entities.Payment{
			ID:             1,
			ServiceOrderID: 1,
			Amount:         100,
			PaymentDate:    time.Now(),
		}, nil)
		result, err := u.CreatePayment(ctx, payment)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID != 1 {
			t.Fatalf("unexpected result: %+v", result)
		}
		mockServiceOrderRepo.AssertExpectations(t)
	})

	t.Run("service order not found", func(t *testing.T) {
		badPayment := &entities.Payment{ID: 2, ServiceOrderID: 2, Amount: 100.0}
		mockServiceOrderRepo.On("GetByID", uint(2)).Return(nil, errors.New("not found"))
		_, err := u.CreatePayment(ctx, badPayment)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		mockServiceOrderRepo.AssertExpectations(t)
	})

	t.Run("amount does not match", func(t *testing.T) {
		badPayment := &entities.Payment{ID: 3, ServiceOrderID: 1, Amount: 200.0}
		mockServiceOrderRepo.On("GetByID", uint(1)).Return(mockServiceOrderDTO, nil)
		_, err := u.CreatePayment(ctx, badPayment)
		if !errors.Is(err, ErrPaymentAmountDoesNotMatch) {
			t.Fatalf("expected ErrPaymentAmountDoesNotMatch, got %v", err)
		}
		mockServiceOrderRepo.AssertExpectations(t)
	})

	t.Run("payment already exists", func(t *testing.T) {
		mockServiceOrderRepo.On("GetByID", uint(1)).Return(mockServiceOrderDTO, nil)
		mockPaymentRepo.EXPECT().GetByServiceOrderID(ctx, uint(1)).Return(entities.Payment{ID: 99}, nil)
		badPayment := &entities.Payment{ID: 4, ServiceOrderID: 1, Amount: 100.0}
		_, err := u.CreatePayment(ctx, badPayment)
		if !errors.Is(err, ErrPaymentAlreadyExists) {
			t.Fatalf("expected ErrPaymentAlreadyExists, got %v", err)
		}
		mockServiceOrderRepo.AssertExpectations(t)
	})

	t.Run("repo create error", func(t *testing.T) {
		mockServiceOrderRepo.On("GetByID", uint(1)).Return(mockServiceOrderDTO, nil)
		mockPaymentRepo.EXPECT().GetByServiceOrderID(ctx, uint(1)).Return(entities.Payment{}, nil)
		mockPaymentRepo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(entities.Payment{})).Return(entities.Payment{}, errors.New("db error"))
		_, err := u.CreatePayment(ctx, payment)
		if err == nil || err.Error() != "db error" {
			t.Fatalf("expected db error, got %v", err)
		}
		mockServiceOrderRepo.AssertExpectations(t)
	})
}

func TestPaymentUseCase_GetPaymentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockPaymentRepo := usecasemocks.NewMockIPaymentRepo(ctrl)
	mockServiceOrderRepo := &domainmocks.MockServiceOrderRepository{}
	u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo.EXPECT().GetByID(ctx, uint(1)).Return(entities.Payment{ID: 1}, nil)
		result, err := u.GetPaymentByID(ctx, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID != 1 {
			t.Fatalf("unexpected result: %+v", result)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockPaymentRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.Payment{}, nil)
		result, err := u.GetPaymentByID(ctx, 2)
		if !errors.Is(err, ErrorPaymentNotFound) {
			t.Fatalf("expected ErrorPaymentNotFound, got %v", err)
		}
		if result == nil {
			t.Fatalf("expected empty payment, got nil")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockPaymentRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.Payment{}, errors.New("db error"))
		_, err := u.GetPaymentByID(ctx, 3)
		if err == nil || err.Error() != "db error" {
			t.Fatalf("expected db error, got %v", err)
		}
	})
}

func TestPaymentUseCase_ListPayments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	mockPaymentRepo := usecasemocks.NewMockIPaymentRepo(ctrl)
	mockServiceOrderRepo := &domainmocks.MockServiceOrderRepository{}
	u := NewPaymentUseCase(mockPaymentRepo, mockServiceOrderRepo)

	t.Run("success", func(t *testing.T) {
		mockPaymentRepo.EXPECT().List(ctx).Return([]entities.Payment{{ID: 1}, {ID: 2}}, nil)
		result, err := u.ListPayments(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 payments, got %d", len(result))
		}
	})

	t.Run("repo error", func(t *testing.T) {
		mockPaymentRepo.EXPECT().List(ctx).Return(nil, errors.New("db error"))
		_, err := u.ListPayments(ctx)
		if err == nil || err.Error() != "db error" {
			t.Fatalf("expected db error, got %v", err)
		}
	})
}
