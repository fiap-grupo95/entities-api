//go:build legacy
// +build legacy

package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/usecase/interfaces"
	"time"

	"mecanica_xpto/internal/domain/entities"
)

var (
	ErrorPaymentNotFound         = errors.New("payment not found")
	ErrPaymentAlreadyExists      = errors.New("payment already exists")
	ErrPaymentAmountDoesNotMatch = errors.New("payment amount does not match service order estimate")
)

type IPaymentUseCase interface {
	CreatePayment(ctx context.Context, payment *entities.Payment) (*entities.Payment, error)
	GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error)
	ListPayments(ctx context.Context) ([]entities.Payment, error)
}

type PaymentUseCase struct {
	repo             interfaces.IPaymentRepo
	serviceOrderRepo interfaces.IServiceOrderRepository
}

var _ IPaymentUseCase = (*PaymentUseCase)(nil)

func NewPaymentUseCase(repo interfaces.IPaymentRepo, serviceOrderRepo interfaces.IServiceOrderRepository) *PaymentUseCase {
	return &PaymentUseCase{
		repo:             repo,
		serviceOrderRepo: serviceOrderRepo,
	}
}

func (p *PaymentUseCase) CreatePayment(ctx context.Context, payment *entities.Payment) (*entities.Payment, error) {
	serviceOrder, err := p.serviceOrderRepo.GetByID(payment.ServiceOrderID)
	if err != nil {
		return nil, err
	}
	if serviceOrder.Estimate != payment.Amount {
		return nil, ErrPaymentAmountDoesNotMatch
	}
	existingPayment, err := p.repo.GetByServiceOrderID(ctx, payment.ServiceOrderID)
	if err != nil {
		return nil, err
	}
	if existingPayment.ID != 0 {
		return nil, ErrPaymentAlreadyExists
	}

	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}

	created, err := p.repo.Create(ctx, *payment)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (p *PaymentUseCase) GetPaymentByID(ctx context.Context, id uint) (*entities.Payment, error) {
	payment, err := p.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment.ID == 0 {
		return &entities.Payment{}, ErrorPaymentNotFound
	}
	return &payment, nil
}

func (p *PaymentUseCase) ListPayments(ctx context.Context) ([]entities.Payment, error) {
	return p.repo.List(ctx)
}
