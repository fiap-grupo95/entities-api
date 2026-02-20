package repository

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase/interfaces"
	"time"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

var _ interfaces.IPaymentRepo = (*PaymentRepository)(nil)

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (p *PaymentRepository) Create(ctx context.Context, payment entities.Payment) (entities.Payment, error) {
	model := dto.PaymentModel{
		ServiceOrderID: payment.ServiceOrderID,
		PaymentDate:    payment.PaymentDate,
		Amount:         payment.Amount,
	}
	if model.PaymentDate.IsZero() {
		model.PaymentDate = time.Now()
	}

	if err := p.db.WithContext(ctx).Create(&model).Error; err != nil {
		return entities.Payment{}, err
	}

	return modelToPayment(model), nil
}

func (p *PaymentRepository) GetByID(ctx context.Context, id uint) (entities.Payment, error) {
	var model dto.PaymentModel
	if err := p.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Payment{}, nil
		}
		return entities.Payment{}, err
	}
	return modelToPayment(model), nil
}

func (p *PaymentRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) (entities.Payment, error) {
	var model dto.PaymentModel
	if err := p.db.WithContext(ctx).Where("service_order_id = ?", serviceOrderID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Payment{}, nil
		}
		return entities.Payment{}, err
	}
	return modelToPayment(model), nil
}

func (p *PaymentRepository) List(ctx context.Context) ([]entities.Payment, error) {
	var models []dto.PaymentModel
	if err := p.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]entities.Payment, 0, len(models))
	for _, model := range models {
		result = append(result, modelToPayment(model))
	}
	return result, nil
}

func modelToPayment(model dto.PaymentModel) entities.Payment {
	return entities.Payment{
		ID:             model.ID,
		ServiceOrderID: model.ServiceOrderID,
		PaymentDate:    model.PaymentDate,
		Amount:         model.Amount,
	}
}
