package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	dto "mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase/interfaces"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ServiceOrderRepository implements IServiceOrderRepository interface
type ServiceOrderRepository struct {
	db *gorm.DB
}

var _ interfaces.IServiceOrderRepository = (*ServiceOrderRepository)(nil)

func NewServiceOrderRepository(db *gorm.DB) *ServiceOrderRepository {
	return &ServiceOrderRepository{db: db}
}

func (r *ServiceOrderRepository) Create(serviceOrder *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	if serviceOrder == nil {
		return nil, gorm.ErrInvalidData
	}

	dtoStatus, err := r.getStatus(serviceOrder.ServiceOrderStatus)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}

	if dtoStatus == nil {
		return nil, gorm.ErrInvalidData
	}

	// Begin transaction
	tx := r.db.Begin()

	serviceOrderDto := dto.ServiceOrderModel{
		ID:                   serviceOrder.ID,
		CustomerID:           serviceOrder.CustomerID,
		VehicleID:            serviceOrder.VehicleID,
		OSStatusID:           dtoStatus.ID,
		Estimate:             serviceOrder.Estimate,
		StartedExecutionDate: serviceOrder.StartedExecutionDate,
		FinalExecutionDate:   serviceOrder.FinalExecutionDate,
		CreatedAt:            serviceOrder.CreatedAt,
		UpdatedAt:            serviceOrder.UpdatedAt,
	}

	if err := tx.Create(&serviceOrderDto).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	serviceOrderDto.ServiceOrderStatus = *dtoStatus

	return serviceOrderDto.ToDomain(), nil
}

func (r *ServiceOrderRepository) GetByID(id uint) (*entities.ServiceOrder, error) {
	var serviceOrder dto.ServiceOrderModel
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessÃ¡rio
	err := r.db.Preload("Customer").
		Preload("Customer.User").
		Preload("Vehicle").
		Preload("ServiceOrderStatus").
		Preload("AdditionalRepairs").
		Preload("AdditionalRepairs.ARStatus").
		Preload("AdditionalRepairs.Services").
		Preload("AdditionalRepairs.PartsSupplies").
		Preload("Payment").
		Preload("PartsSupplies").
		Preload("Services").
		First(&serviceOrder, id).Error
	if err != nil {
		log.Error().Msgf("Error finding service order with id %d: %v", id, err)
		if strings.EqualFold(err.Error(), gorm.ErrRecordNotFound.Error()) {
			return nil, nil
		}
		return nil, err
	}
	return serviceOrder.ToDomain(), nil
}

func (r *ServiceOrderRepository) UpdateEstimate(id uint, estimate float64) error {
	var dtoDB dto.ServiceOrderModel
	if err := r.db.First(&dtoDB, id).Error; err != nil {
		return err
	}
	newEstimate := dtoDB.Estimate + estimate

	return r.db.Model(&dto.ServiceOrderModel{}).
		Where("id = ?", id).
		Update("estimate", newEstimate).Error
}

func (r *ServiceOrderRepository) Update(serviceOrder *entities.ServiceOrder) error {
	if serviceOrder == nil {
		return gorm.ErrInvalidData
	}

	dtoStatus, err := r.getStatus(serviceOrder.ServiceOrderStatus)
	if err != nil {
		return gorm.ErrInvalidData
	}

	tx := r.db.Begin()

	serviceOrderDto := dto.ServiceOrderModel{
		ID:                       serviceOrder.ID,
		CustomerID:               serviceOrder.CustomerID,
		VehicleID:                serviceOrder.VehicleID,
		OSStatusID:               dtoStatus.ID,
		Estimate:                 serviceOrder.Estimate,
		StartedExecutionDate:     serviceOrder.StartedExecutionDate,
		FinalExecutionDate:       serviceOrder.FinalExecutionDate,
		ExecutionDurationInHours: serviceOrder.ExecutionDurationInHours,
	}

	if err := tx.Model(&dto.ServiceOrderModel{}).Where("id = ?", serviceOrder.ID).Updates(&serviceOrderDto).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update PartsSupplies relationships
	if serviceOrder.PartsSupplies != nil {
		if err := tx.Where("service_order_id = ?", serviceOrder.ID).Delete(&dto.PartsSupplyServiceOrder{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, partsSupply := range serviceOrder.PartsSupplies {
			relation := dto.PartsSupplyServiceOrder{
				PartsSupplyID:  partsSupply.ID,
				ServiceOrderID: serviceOrder.ID,
				Quantity:       partsSupply.QuantityReserve,
			}
			if err := tx.Create(&relation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Update Services relationships
	if serviceOrder.Services != nil {
		if err := tx.Where("service_order_id = ?", serviceOrder.ID).Delete(&dto.ServiceServiceOrder{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, service := range serviceOrder.Services {
			relation := dto.ServiceServiceOrder{
				ServiceID:      service.ID,
				ServiceOrderID: serviceOrder.ID,
			}
			if err := tx.Create(&relation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *ServiceOrderRepository) List() ([]*entities.ServiceOrder, error) {
	var serviceOrders []dto.ServiceOrderModel
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessÃ¡rio
	err := r.db.
		Preload("Customer").
		Preload("Customer.User").
		Preload("Vehicle").
		Preload("ServiceOrderStatus").
		Preload("AdditionalRepairs").
		Preload("AdditionalRepairs.ARStatus").
		Preload("AdditionalRepairs.Services").
		Preload("AdditionalRepairs.PartsSupplies").
		Preload("Payment").
		Preload("PartsSupplies").
		Preload("Services").
		Find(&serviceOrders).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.ServiceOrder, 0, len(serviceOrders))
	for _, so := range serviceOrders {
		result = append(result, so.ToDomain())
	}
	return result, nil
}

func (r *ServiceOrderRepository) getStatus(status valueobject.ServiceOrderStatus) (*dto.ServiceOrderStatus, error) {
	var serviceOrderStatuses dto.ServiceOrderStatus
	err := r.db.Where("description = ?", status.String()).First(&serviceOrderStatuses).Error
	if err != nil {
		return nil, err
	}
	return &serviceOrderStatuses, nil
}

func (r *ServiceOrderRepository) GetPartsSupplyServiceOrder(partsSupplyID uint, serviceOrderID uint) (*entities.ServiceOrderPartsSupply, error) {
	var relation dto.PartsSupplyServiceOrder
	err := r.db.Where("parts_supply_id = ? AND service_order_id = ?", partsSupplyID, serviceOrderID).First(&relation).Error
	if err != nil {
		return nil, err
	}
	return &entities.ServiceOrderPartsSupply{
		PartsSupplyID:  relation.PartsSupplyID,
		ServiceOrderID: relation.ServiceOrderID,
		Quantity:       relation.Quantity,
	}, nil
}
