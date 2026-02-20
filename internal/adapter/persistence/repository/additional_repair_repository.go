package repository

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase/interfaces"

	"gorm.io/gorm"
)

// AdditionalRepairRepository implements IAdditionalRepairRepository interface
type AdditionalRepairRepository struct {
	db *gorm.DB
}

func NewAdditionalRepairRepository(db *gorm.DB) interfaces.IAdditionalRepairRepository {
	return &AdditionalRepairRepository{db: db}
}

func (r *AdditionalRepairRepository) Create(ctx context.Context, additionalRepair entities.AdditionalRepair) (entities.AdditionalRepair, error) {
	statusDTO, err := r.getStatus(ctx, additionalRepair.ARStatus.String())
	if err != nil {
		return entities.AdditionalRepair{}, err
	}

	model := dto.AdditionalRepairModel{
		ID:             additionalRepair.ID,
		Description:    additionalRepair.Description,
		ServiceOrderID: additionalRepair.ServiceOrderID,
		ARStatusID:     statusDTO.ID,
		Estimate:       additionalRepair.Estimate,
		Services:       mapServicesToModels(additionalRepair.Services),
		PartsSupplies:  mapPartsSuppliesToModels(additionalRepair.PartsSupplies),
	}
	model.ARStatus = *statusDTO

	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(&model).Error; err != nil {
		tx.Rollback()
		return entities.AdditionalRepair{}, err
	}

	if err := updatePartsSupplyQuantities(tx, model.ID, additionalRepair.PartsSupplies); err != nil {
		tx.Rollback()
		return entities.AdditionalRepair{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return entities.AdditionalRepair{}, err
	}

	created := model.ToDomain()
	return created, nil
}

func (r *AdditionalRepairRepository) GetByID(ctx context.Context, id uint) (entities.AdditionalRepair, error) {
	var additionalRepair dto.AdditionalRepairModel
	err := r.db.WithContext(ctx).
		Preload("ARStatus").
		Preload("PartsSupplies").
		Preload("Services").
		First(&additionalRepair, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.AdditionalRepair{}, nil
		}
		return entities.AdditionalRepair{}, err
	}
	return additionalRepair.ToDomain(), nil
}

func (r *AdditionalRepairRepository) AddPartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error {
	tx := r.db.WithContext(ctx).Begin()

	for _, svc := range services {
		relation := dto.ServiceAdditionalRepair{
			ServiceID:          svc.ID,
			AdditionalRepairID: additionalRepairID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, ps := range partsSupplies {
		relation := dto.PartsSupplyAdditionalRepair{
			PartsSupplyID:      ps.ID,
			AdditionalRepairID: additionalRepairID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&dto.PartsSupplyAdditionalRepair{}).
			Where("parts_supply_id = ? and additional_repair_id = ?", ps.ID, additionalRepairID).
			Update("quantity", ps.QuantityReserve).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&dto.AdditionalRepairModel{}).
		Where("id = ?", additionalRepairID).
		Update("estimate", newEstimate).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) ReplacePartSupplyAndService(ctx context.Context, additionalRepairID uint, services []entities.Service, partsSupplies []entities.PartsSupply, newEstimate float64) error {
	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Where("additional_repair_id = ?", additionalRepairID).
		Delete(&dto.PartsSupplyAdditionalRepair{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("additional_repair_id = ?", additionalRepairID).
		Delete(&dto.ServiceAdditionalRepair{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, svc := range services {
		relation := dto.ServiceAdditionalRepair{
			ServiceID:          svc.ID,
			AdditionalRepairID: additionalRepairID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, ps := range partsSupplies {
		relation := dto.PartsSupplyAdditionalRepair{
			PartsSupplyID:      ps.ID,
			AdditionalRepairID: additionalRepairID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&dto.PartsSupplyAdditionalRepair{}).
			Where("parts_supply_id = ? and additional_repair_id = ?", ps.ID, additionalRepairID).
			Update("quantity", ps.QuantityReserve).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&dto.AdditionalRepairModel{}).
		Where("id = ?", additionalRepairID).
		Update("estimate", newEstimate).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) UpdateStatus(ctx context.Context, id uint, status entities.AdditionalRepairStatusDTO) error {
	dtoStatus, err := r.getStatus(ctx, status.ApprovalStatus)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&dto.AdditionalRepairModel{}).
		Where("id = ?", id).
		Update("ar_status_id", dtoStatus.ID).Error
}

func (r *AdditionalRepairRepository) GetByServiceOrder(ctx context.Context, serviceOrderId uint) ([]entities.AdditionalRepair, error) {
	var additionalRepairs []dto.AdditionalRepairModel
	err := r.db.WithContext(ctx).
		Preload("ARStatus").
		Preload("PartsSupplies").
		Preload("Services").
		Where("service_order_id = ?", serviceOrderId).
		Find(&additionalRepairs).Error
	if err != nil {
		return nil, err
	}

	result := make([]entities.AdditionalRepair, 0, len(additionalRepairs))
	for _, ar := range additionalRepairs {
		result = append(result, ar.ToDomain())
	}
	return result, nil
}

func (r *AdditionalRepairRepository) GetPartsSupplyQuantity(ctx context.Context, partsSupplyID uint, additionalRepairID uint) (int, error) {
	var relation dto.PartsSupplyAdditionalRepair
	err := r.db.WithContext(ctx).
		Where("parts_supply_id = ? AND additional_repair_id = ?", partsSupplyID, additionalRepairID).
		First(&relation).Error
	if err != nil {
		return 0, err
	}
	return relation.Quantity, nil
}

func mapServicesToModels(services []entities.Service) []dto.ServiceModel {
	if len(services) == 0 {
		return nil
	}
	result := make([]dto.ServiceModel, 0, len(services))
	for _, svc := range services {
		result = append(result, dto.ServiceModel{
			ID:    svc.ID,
			Name:  svc.Name,
			Price: svc.Price,
		})
	}
	return result
}

func mapPartsSuppliesToModels(partsSupplies []entities.PartsSupply) []dto.PartsSupplyModel {
	if len(partsSupplies) == 0 {
		return nil
	}
	result := make([]dto.PartsSupplyModel, 0, len(partsSupplies))
	for _, ps := range partsSupplies {
		result = append(result, dto.PartsSupplyModel{
			ID:              ps.ID,
			Name:            ps.Name,
			Price:           ps.Price,
			QuantityTotal:   ps.QuantityTotal,
			QuantityReserve: ps.QuantityReserve,
		})
	}
	return result
}

func updatePartsSupplyQuantities(tx *gorm.DB, additionalRepairID uint, partsSupplies []entities.PartsSupply) error {
	for _, ps := range partsSupplies {
		if err := tx.Model(&dto.PartsSupplyAdditionalRepair{}).
			Where("parts_supply_id = ? and additional_repair_id = ?", ps.ID, additionalRepairID).
			Update("quantity", ps.QuantityReserve).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *AdditionalRepairRepository) getStatus(ctx context.Context, status string) (*dto.AdditionalRepairStatusModel, error) {
	var additionalRepairStatus dto.AdditionalRepairStatusModel
	err := r.db.WithContext(ctx).Where("description = ?", status).First(&additionalRepairStatus).Error
	if err != nil {
		return nil, err
	}
	return &additionalRepairStatus, nil
}
