package dto

import (
	"mecanica_xpto/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type ServiceModel struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"size:100;not null"`
	Description string         `gorm:"type:text"`
	Price       float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (m *ServiceModel) TableName() string {
	return "service"
}

func (m *ServiceModel) ToDomain() entities.Service {
	return entities.Service{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt: func() *time.Time {
			if m.DeletedAt.Valid {
				return &m.DeletedAt.Time
			}
			return nil
		}(),
	}
}
