package dto

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"time"
)

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrder struct {
	PartsSupplyID  uint `gorm:"column:parts_supply_id;primaryKey"`
	ServiceOrderID uint `gorm:"column:service_order_id;primaryKey"`
	Quantity       int  `gorm:"column:quantity"`
}

// N:N relationship between Service and ServiceOrder
type ServiceServiceOrder struct {
	ServiceID      uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
}

type ServiceOrderStatus struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"size:50;not null"`
}

func (m *ServiceOrderStatus) TableName() string {
	return "service_order_status"
}

func (m *ServiceOrderStatus) ToDomain() valueobject.ServiceOrderStatus {
	return valueobject.ParseServiceOrderStatus(m.Description)
}

type ServiceOrderModel struct {
	ID                       uint               `gorm:"primaryKey"`
	CustomerID               uint               `gorm:"not null"`
	Customer                 CustomerModel      `gorm:"foreignKey:CustomerID"`
	VehicleID                uint               `gorm:"not null"`
	Vehicle                  VehicleModel       `gorm:"foreignKey:VehicleID"`
	OSStatusID               uint               `gorm:"not null"`
	ServiceOrderStatus       ServiceOrderStatus `gorm:"foreignKey:OSStatusID"`
	Estimate                 float64            `gorm:"type:decimal(10,2)"`
	StartedExecutionDate     *time.Time
	FinalExecutionDate       *time.Time
	ExecutionDurationInHours float64
	CreatedAt                *time.Time              `gorm:"autoCreateTime"`
	UpdatedAt                *time.Time              `gorm:"autoUpdateTime"`
	AdditionalRepairs        []AdditionalRepairModel `gorm:"foreignKey:ServiceOrderID"`
	Payment                  *PaymentModel           `gorm:"foreignKey:ServiceOrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PartsSupplies            []PartsSupplyModel      `gorm:"many2many:parts_supply_service_orders;foreignKey:ID;joinForeignKey:ServiceOrderID;references:ID;joinReferences:PartsSupplyID"`
	Services                 []ServiceModel          `gorm:"many2many:service_service_orders;foreignKey:ID;joinForeignKey:ServiceOrderID;references:ID;joinReferences:ServiceID"`
}

func (m *ServiceOrderModel) TableName() string {
	return "service_order"
}

func (m *ServiceOrderModel) ToDomain() *entities.ServiceOrder {
	var additionalRepairs []entities.AdditionalRepair
	var partsSupplies []entities.PartsSupply
	var services []entities.Service

	// Convert AdditionalRepairs
	for _, ar := range m.AdditionalRepairs {
		additionalRepairs = append(additionalRepairs, ar.ToDomain())
	}

	// Convert PartsSupplies
	for _, ps := range m.PartsSupplies {
		partsSupplies = append(partsSupplies, ps.ToDomain())
	}

	// Convert Services
	for _, s := range m.Services {
		services = append(services, s.ToDomain())
	}

	// Convert Payment if exists
	var payment *entities.Payment
	if m.Payment != nil {
		p := m.Payment.ToDomain()
		payment = p
	}

	// Convert Customer and Vehicle if they are loaded
	var customer entities.Customer
	var vehicle entities.Vehicle
	if m.Customer.ID != 0 {
		customer = *m.Customer.ToDomain()
	}
	if m.Vehicle.ID != 0 {
		vehicle = *m.Vehicle.ToDomain()
	}

	return &entities.ServiceOrder{
		ID:                       m.ID,
		CustomerID:               m.CustomerID,
		Customer:                 &customer,
		VehicleID:                m.VehicleID,
		Vehicle:                  &vehicle,
		ServiceOrderStatus:       m.ServiceOrderStatus.ToDomain(),
		Estimate:                 m.Estimate,
		StartedExecutionDate:     m.StartedExecutionDate,
		FinalExecutionDate:       m.FinalExecutionDate,
		ExecutionDurationInHours: m.ExecutionDurationInHours,
		CreatedAt:                m.CreatedAt,
		UpdatedAt:                m.UpdatedAt,
		AdditionalRepairs:        additionalRepairs,
		Payment:                  payment,
		PartsSupplies:            partsSupplies,
		Services:                 services,
	}
}
