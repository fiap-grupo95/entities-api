//go:build legacy
// +build legacy

package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/usecase/interfaces"

	"github.com/rs/zerolog/log"

	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
)

var (
	ErrAdditionalRepairNotFound = errors.New("additional repair not found")
	ErrStatusNotPermitted       = errors.New("additional repair status not permitted")
)

type IAdditionalRepairUseCase interface {
	CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) (entities.AdditionalRepair, error)
	AddPartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	RemovePartSupplyAndService(ctx context.Context, adrId uint, adr entities.AdditionalRepair) error
	GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error)
	CustomerApprovalStatus(ctx context.Context, additionalRepairId uint, status entities.AdditionalRepairStatusDTO) error
}

type AdditionalRepairUseCase struct {
	repo            interfaces.IAdditionalRepairRepository
	repoOS          interfaces.IServiceOrderRepository
	serviceRepo     interfaces.IServiceRepo
	partsSupplyRepo interfaces.IPartsSupplyRepo
}

var _ IAdditionalRepairUseCase = (*AdditionalRepairUseCase)(nil)

func NewSOAdditionalRepairUseCase(repo interfaces.IAdditionalRepairRepository, repoOS interfaces.IServiceOrderRepository, serviceRepo interfaces.IServiceRepo, partsSupplyRepo interfaces.IPartsSupplyRepo) *AdditionalRepairUseCase {
	return &AdditionalRepairUseCase{
		repo:            repo,
		repoOS:          repoOS,
		serviceRepo:     serviceRepo,
		partsSupplyRepo: partsSupplyRepo,
	}
}

func (u *AdditionalRepairUseCase) CreateAdditionalRepair(ctx context.Context, adr entities.AdditionalRepair) (entities.AdditionalRepair, error) {
	_, err := u.repoOS.GetByID(adr.ServiceOrderID)
	if err != nil {
		log.Error().Msgf("error finding service order with id %d: %v", adr.ServiceOrderID, err)
		return entities.AdditionalRepair{}, err
	}

	services, estimatedService, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return entities.AdditionalRepair{}, err
	}

	partsSupplies, estimatedPartsSupply, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return entities.AdditionalRepair{}, err
	}

	newAdditionalRepair := entities.AdditionalRepair{
		ServiceOrderID: adr.ServiceOrderID,
		Description:    adr.Description,
		ARStatus:       valueobject.AdditionalRepairStatus("IN_ANALYSIS"),
		Estimate:       estimatedService + estimatedPartsSupply,
		Services:       services,
		PartsSupplies:  partsSupplies,
	}

	created, err := u.repo.Create(ctx, newAdditionalRepair)
	if err != nil {
		log.Error().Msgf("Error creating additional repair: %v", err)
		return entities.AdditionalRepair{}, err
	}
	return created, nil
}

func (u *AdditionalRepairUseCase) AddPartSupplyAndService(ctx context.Context, additionalRepairId uint, adr entities.AdditionalRepair) error {
	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return err
	}
	if additionalRepair.ID == 0 {
		return ErrAdditionalRepairNotFound
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepair.ARStatus.String()); err != nil {
		return err
	}

	services, estimatedService, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}

	partsSupplies, estimatedPartsSupply, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	newEstimate := additionalRepair.Estimate + estimatedService + estimatedPartsSupply
	if err := u.repo.AddPartSupplyAndService(ctx, additionalRepairId, services, partsSupplies, newEstimate); err != nil {
		log.Error().Msgf("Error adding part supply and services for additional repair: %v", err)
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) RemovePartSupplyAndService(ctx context.Context, additionalRepairId uint, adr entities.AdditionalRepair) error {
	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return err
	}
	if additionalRepair.ID == 0 {
		return ErrAdditionalRepairNotFound
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepair.ARStatus.String()); err != nil {
		return err
	}

	services, estimatedService, err := u.addServiceToAdditionalRepair(ctx, adr.Services)
	if err != nil {
		return err
	}
	partsSupplies, estimatedPartsSupply, err := u.addPartsSupplyToAdditionalRepair(ctx, adr.PartsSupplies)
	if err != nil {
		return err
	}

	newEstimate := estimatedService + estimatedPartsSupply
	if err := u.repo.ReplacePartSupplyAndService(ctx, additionalRepairId, services, partsSupplies, newEstimate); err != nil {
		log.Error().Msgf("Error updating parts supply and services for additional repair: %v", err)
		return err
	}
	return nil
}

func (u *AdditionalRepairUseCase) GetAdditionalRepair(ctx context.Context, additionalRepairId uint) (entities.AdditionalRepair, error) {
	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return entities.AdditionalRepair{}, err
	}
	if additionalRepair.ID == 0 {
		log.Error().Msgf("additional repair with id %d not found", additionalRepairId)
		return entities.AdditionalRepair{}, ErrAdditionalRepairNotFound
	}
	return additionalRepair, nil
}

func (u *AdditionalRepairUseCase) CustomerApprovalStatus(ctx context.Context, additionalRepairId uint, status entities.AdditionalRepairStatusDTO) error {
	additionalRepair, err := u.repo.GetByID(ctx, additionalRepairId)
	if err != nil {
		log.Error().Msgf("error finding additional repair with id %d: %v", additionalRepairId, err)
		return err
	}
	if additionalRepair.ID == 0 {
		return ErrAdditionalRepairNotFound
	}
	if err := u.ValidateAdditionalRepairStatus(additionalRepair.ARStatus.String()); err != nil {
		return err
	}

	if status.ApprovalStatus == "DENIED" {
		for _, ps := range additionalRepair.PartsSupplies {
			if err := unreservePartsSupply(ctx, ps, u.partsSupplyRepo); err != nil {
				log.Error().Msgf("Error unreserving parts supply: %v", err)
			}
		}
		log.Info().Msgf("Customer rejected additional repair with id %d", additionalRepairId)
	} else {
		for _, ps := range additionalRepair.PartsSupplies {
			quantity, err := u.repo.GetPartsSupplyQuantity(ctx, ps.ID, additionalRepair.ID)
			if err != nil {
				log.Error().Msgf("Error getting parts supply additional repair relation: %v", err)
				return err
			}

			entity := entities.PartsSupply{
				ID:              ps.ID,
				QuantityReserve: quantity,
				QuantityTotal:   quantity,
			}
			if err := releaseReservedPartsSupply(ctx, entity, u.partsSupplyRepo); err != nil {
				log.Error().Msgf("Error releasing reserved parts supply: %v", err)
				return err
			}
		}
	}

	if err := u.repo.UpdateStatus(ctx, additionalRepairId, status); err != nil {
		log.Error().Msgf("error updating customer approval with id %d: %v", additionalRepairId, err)
		return err
	}

	if err := u.repoOS.UpdateEstimate(additionalRepair.ServiceOrderID, additionalRepair.Estimate); err != nil {
		log.Error().Msgf("error updating service order estimate with id %d: %v", additionalRepair.ServiceOrderID, err)
		return err
	}

	return nil
}

func (u *AdditionalRepairUseCase) addPartsSupplyToAdditionalRepair(ctx context.Context, partsSupplyId []entities.PartsSupply) ([]entities.PartsSupply, float64, error) {
	var listPartsSupply []entities.PartsSupply
	var estimatedPrice float64

	for _, ps := range partsSupplyId {
		err := reservePartsSupply(ctx, ps, u.partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error reserving parts supply: %v", err)
			return nil, 0, err
		}
		psDto, err := u.partsSupplyRepo.GetByID(ctx, ps.ID)
		if err != nil {
			log.Error().Msgf("error finding parts supply with id %d: %v", ps.ID, err)
		}
		if psDto.ID == 0 {
			log.Error().Msgf("parts supply with id %d not found", ps.ID)
			return listPartsSupply, estimatedPrice, ErrServiceNotFound
		}
		psDto.QuantityReserve = ps.QuantityReserve
		if ps.QuantityTotal > 0 {
			psDto.QuantityTotal = ps.QuantityTotal
		}
		estimatedPrice += psDto.Price * float64(psDto.QuantityReserve)
		listPartsSupply = append(listPartsSupply, psDto)
	}
	return listPartsSupply, estimatedPrice, nil
}

func (u *AdditionalRepairUseCase) addServiceToAdditionalRepair(ctx context.Context, services []entities.Service) ([]entities.Service, float64, error) {
	var listServices []entities.Service
	var estimatedPrice float64

	for _, s := range services {
		serviceDto, err := u.serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			log.Error().Msgf("error finding service with id %d: %v", s.ID, err)
		}
		if serviceDto.ID == 0 {
			log.Error().Msgf("service with id %d not found", s.ID)
			return listServices, estimatedPrice, ErrServiceNotFound
		}
		estimatedPrice += serviceDto.Price
		listServices = append(listServices, serviceDto)
	}
	return listServices, estimatedPrice, nil
}

func (u *AdditionalRepairUseCase) ValidateAdditionalRepairStatus(status string) error {
	if status != "IN_ANALYSIS" {
		log.Error().Msgf("invalid additional repair status: %s", status)
		return ErrStatusNotPermitted
	}
	return nil
}
