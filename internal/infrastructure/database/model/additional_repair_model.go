package dto

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"time"
)

// N:N relationship between PartsSupply and AdditionalRepair
type PartsSupplyAdditionalRepair struct {
	PartsSupplyID      uint `gorm:"column:parts_supply_id;primaryKey"`
	AdditionalRepairID uint `gorm:"column:additional_repair_id;primaryKey"`
	Quantity           int  `gorm:"column:quantity"`
}

func (m PartsSupplyAdditionalRepair) TableName() string {
	return "parts_supply_additional_repairs"
}

// N:N relationship between Service and AdditionalRepair
type ServiceAdditionalRepair struct {
	ServiceID          uint `gorm:"primaryKey"`
	AdditionalRepairID uint `gorm:"primaryKey"`
}

func (m ServiceAdditionalRepair) TableName() string {
	return "service_additional_repairs"
}

// 1:N relationship between AdditionalRepair and AdditionalRepairStatus
type AdditionalRepairStatusModel struct {
	ID                uint                    `gorm:"primaryKey"`
	Description       string                  `gorm:"size:50;not null"`
	AdditionalRepairs []AdditionalRepairModel `gorm:"foreignKey:ARStatusID"`
}

func (arsm *AdditionalRepairStatusModel) TableName() string {
	return "additional_repair_status"
}

func (arsm *AdditionalRepairStatusModel) ToDomain() valueobject.AdditionalRepairStatus {
	return valueobject.ParseAdditionalRepairStatus(arsm.Description)
}

type AdditionalRepairModel struct {
	ID             uint                        `gorm:"primaryKey"`
	Description    string                      `gorm:"column:description;not null"`
	ServiceOrderID uint                        `gorm:"column:service_order_id;not null"`
	ServiceOrder   ServiceOrderModel           `gorm:"foreignKey:ServiceOrderID"`
	ARStatusID     uint                        `gorm:"not null"`
	ARStatus       AdditionalRepairStatusModel `gorm:"foreignKey:ARStatusID"`
	Estimate       float64                     `gorm:"type:decimal(10,2)"`
	CreatedAt      time.Time                   `gorm:"autoCreateTime"`
	UpdatedAt      time.Time                   `gorm:"autoUpdateTime"`
	PartsSupplies  []PartsSupplyModel          `gorm:"many2many:parts_supply_additional_repairs;foreignKey:ID;joinForeignKey:AdditionalRepairID;references:ID;joinReferences:PartsSupplyID"`
	Services       []ServiceModel              `gorm:"many2many:service_additional_repairs;foreignKey:ID;joinForeignKey:AdditionalRepairID;references:ID;joinReferences:ServiceID"`
}

func (arm *AdditionalRepairModel) TableName() string {
	return "additional_repair"
}

func (arm *AdditionalRepairModel) ToDomain() entities.AdditionalRepair {
	var partsSupplies []entities.PartsSupply
	var services []entities.Service

	for _, ps := range arm.PartsSupplies {
		partsSupplies = append(partsSupplies, ps.ToDomain())
	}

	for _, s := range arm.Services {
		services = append(services, s.ToDomain())
	}

	return entities.AdditionalRepair{
		ID:             arm.ID,
		Description:    arm.Description,
		ARStatus:       arm.ARStatus.ToDomain(),
		Estimate:       arm.Estimate,
		CreatedAt:      arm.CreatedAt,
		UpdatedAt:      arm.UpdatedAt,
		ServiceOrderID: arm.ServiceOrderID,
		ServiceOrder:   arm.ServiceOrder.ToDomain(),
		PartsSupplies:  partsSupplies,
		Services:       services,
	}
}
