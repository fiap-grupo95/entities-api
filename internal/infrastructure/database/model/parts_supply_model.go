package dto

import (
	"mecanica_xpto/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

// N:N relationship between PartsSupply and ServiceOrder
// 1:N relationship between PartsSupply and AdditionalRepair
type PartsSupplyModel struct {
	ID              uint           `gorm:"primaryKey"`
	Name            string         `gorm:"size:100;not null"`
	Description     string         `gorm:"type:text"`
	Price           float64        `gorm:"type:decimal(10,2);not null"`
	QuantityTotal   int            `gorm:"not null;default:0"`
	QuantityReserve int            `gorm:"not null;default:0"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (m *PartsSupplyModel) TableName() string {
	return "parts_supply"
}

func (m *PartsSupplyModel) ToDomain() entities.PartsSupply {
	return entities.PartsSupply{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		Price:           m.Price,
		QuantityTotal:   m.QuantityTotal,
		QuantityReserve: m.QuantityReserve,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt: func() *time.Time {
			if m.DeletedAt.Valid {
				return &m.DeletedAt.Time
			}
			return nil
		}(),
	}
}
