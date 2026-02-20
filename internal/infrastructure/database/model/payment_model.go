package dto

import (
	"mecanica_xpto/internal/domain/entities"
	"time"
)

type PaymentModel struct {
	ID             uint              `gorm:"primaryKey"`
	ServiceOrderID uint              `gorm:"unique;not null"`
	ServiceOrder   ServiceOrderModel `gorm:"foreignKey:ServiceOrderID;references:ID"`
	PaymentDate    time.Time         `gorm:"not null"`
	Amount         float64           `gorm:"not null"`
}

func (pm *PaymentModel) TableName() string {
	return "payment"
}

func (pm *PaymentModel) ToDomain() *entities.Payment {
	return &entities.Payment{
		ID:             pm.ID,
		ServiceOrderID: pm.ServiceOrder.ToDomain().ID,
		PaymentDate:    pm.PaymentDate,
		Amount:         pm.Amount,
	}
}
