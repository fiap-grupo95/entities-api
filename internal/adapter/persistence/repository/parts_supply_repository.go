package repository

//go:generate mockgen -destination=../mocks/parts_supply_repository_mock.go -package=mocks mecanica_xpto/internal/domain PartsSupplyRepo

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/entities"
	dto "mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase/interfaces"

	"gorm.io/gorm"
)

type PartsSupplyRepository struct {
	db *gorm.DB
}

var _ interfaces.IPartsSupplyRepo = (*PartsSupplyRepository)(nil)

func NewPartsSupplyRepository(db *gorm.DB) *PartsSupplyRepository {
	return &PartsSupplyRepository{db: db}
}

func (s *PartsSupplyRepository) Create(ctx context.Context, ps *entities.PartsSupply) (entities.PartsSupply, error) {
	dto := dto.PartsSupplyModel{
		Name:            ps.Name,
		Description:     ps.Description,
		Price:           ps.Price,
		QuantityTotal:   ps.QuantityTotal,
		QuantityReserve: ps.QuantityReserve,
	}
	if err := s.db.WithContext(ctx).Create(&dto).Error; err != nil {
		return entities.PartsSupply{}, err
	}
	return dto.ToDomain(), nil
}

func (s *PartsSupplyRepository) GetByID(ctx context.Context, id uint) (entities.PartsSupply, error) {
	var dto dto.PartsSupplyModel
	if err := s.db.WithContext(ctx).First(&dto, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PartsSupply{}, nil
		}
		return entities.PartsSupply{}, err
	}
	return dto.ToDomain(), nil
}

func (s *PartsSupplyRepository) GetByName(ctx context.Context, name string) (entities.PartsSupply, error) {
	var dto dto.PartsSupplyModel
	if err := s.db.WithContext(ctx).Where("name = ?", name).Find(&dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PartsSupply{}, nil
		}
		return entities.PartsSupply{}, err
	}
	return dto.ToDomain(), nil
}

func (s *PartsSupplyRepository) Update(ctx context.Context, ps *entities.PartsSupply) error {
	var dtoDB dto.PartsSupplyModel
	if err := s.db.First(&dtoDB, ps.ID).Error; err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if ps.Name != dtoDB.Name {
		updates["name"] = ps.Name
	}
	if ps.Description != dtoDB.Description {
		updates["description"] = ps.Description
	}
	if ps.Price != dtoDB.Price {
		updates["price"] = ps.Price
	}
	if ps.QuantityTotal != dtoDB.QuantityTotal {
		updates["quantity_total"] = ps.QuantityTotal
	}
	if ps.QuantityReserve != dtoDB.QuantityReserve {
		updates["quantity_reserve"] = ps.QuantityReserve
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).
		Model(&dto.PartsSupplyModel{}).
		Where("id = ?", ps.ID).
		Updates(updates).Error
}

func (s *PartsSupplyRepository) Delete(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&dto.PartsSupplyModel{}, id).Error
}

func (s *PartsSupplyRepository) List(ctx context.Context) ([]entities.PartsSupply, error) {
	var dtos []dto.PartsSupplyModel
	if err := s.db.WithContext(ctx).Find(&dtos).Error; err != nil {
		return nil, err
	}
	result := make([]entities.PartsSupply, len(dtos))
	for i, dto := range dtos {
		result[i] = dto.ToDomain()
	}
	return result, nil
}

func (s *PartsSupplyRepository) GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error) {
	var dtos []dto.PartsSupplyModel
	if err := s.db.WithContext(ctx).
		Joins("JOIN parts_supply_service_orders ON parts_supply_service_orders.parts_supply_id = parts_supply.id").
		Where("parts_supply_service_orders.service_order_id = ?", serviceOrderID).
		Find(&dtos).Error; err != nil {
		return nil, err
	}

	var partsSupplies []entities.PartsSupply
	for _, dto := range dtos {
		partsSupplies = append(partsSupplies, dto.ToDomain())
	}
	return partsSupplies, nil
}
