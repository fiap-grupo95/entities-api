package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/usecase/interfaces"
)

type IPartsSupplyUseCase interface {
	GetPartsSupplyByID(ctx context.Context, id uint) (entities.PartsSupply, error)
	CreatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) (entities.PartsSupply, error)
	UpdatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) error
	DeletePartsSupply(ctx context.Context, id uint) error
	ListPartsSupplies(ctx context.Context) ([]entities.PartsSupply, error)
	GetPartsSupplyByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error)
	AuthorizeReserve(ctx context.Context, partsSupply []entities.PartsSupply) error
	Reserve(ctx context.Context, partsSupply []entities.PartsSupply) error
	Release(ctx context.Context, partsSupply []entities.PartsSupply) error
	WriteOff(ctx context.Context, partsSupply []entities.PartsSupply) error
}
type PartsSupplyUseCase struct {
	repo interfaces.IPartsSupplyRepo
}

var _ IPartsSupplyUseCase = (*PartsSupplyUseCase)(nil)

func NewPartsSupplyUseCase(repo interfaces.IPartsSupplyRepo) *PartsSupplyUseCase {
	return &PartsSupplyUseCase{repo: repo}
}

var (
	ErrPartsSupplyNotFound      = errors.New("parts supply not found")
	ErrPartsSupplyAlreadyExists = errors.New("parts supply already exists")
	ErrInsufficientStock        = errors.New("insufficient parts supply stock")
	ErrInsufficientReserved     = errors.New("insufficient reserved parts supply")
	ErrInvalidQuantity          = errors.New("invalid quantity")
)

func (h *PartsSupplyUseCase) AuthorizeReserve(ctx context.Context, partsSupply []entities.PartsSupply) error {
	for _, item := range partsSupply {
		current, err := h.repo.GetByID(ctx, item.ID)
		if err != nil {
			return err
		}

		if current.ID == 0 {
			return ErrPartsSupplyNotFound
		}

		if item.QuantityReserve <= 0 {
			return ErrInvalidQuantity
		}

		available := current.QuantityTotal - current.QuantityReserve
		if item.QuantityReserve > available {
			return ErrInsufficientStock
		}
	}

	return nil
}

func (h *PartsSupplyUseCase) Reserve(ctx context.Context, partsSupply []entities.PartsSupply) error {
	if err := h.AuthorizeReserve(ctx, partsSupply); err != nil {
		return err
	}

	for _, item := range partsSupply {
		current, err := h.repo.GetByID(ctx, item.ID)
		if err != nil {
			return err
		}

		current.QuantityReserve += item.QuantityReserve
		if err = h.repo.Update(ctx, &current); err != nil {
			return err
		}
	}

	return nil
}

func (h *PartsSupplyUseCase) Release(ctx context.Context, partsSupply []entities.PartsSupply) error {
	for _, item := range partsSupply {
		current, err := h.repo.GetByID(ctx, item.ID)
		if err != nil {
			return err
		}

		if current.ID == 0 {
			return ErrPartsSupplyNotFound
		}

		if item.QuantityReserve <= 0 {
			return ErrInvalidQuantity
		}

		if item.QuantityReserve > current.QuantityReserve {
			return ErrInsufficientReserved
		}

		current.QuantityReserve -= item.QuantityReserve
		if err = h.repo.Update(ctx, &current); err != nil {
			return err
		}
	}

	return nil
}

func (h *PartsSupplyUseCase) WriteOff(ctx context.Context, partsSupply []entities.PartsSupply) error {
	for _, item := range partsSupply {
		current, err := h.repo.GetByID(ctx, item.ID)
		if err != nil {
			return err
		}

		if current.ID == 0 {
			return ErrPartsSupplyNotFound
		}

		if item.QuantityReserve <= 0 {
			return ErrInvalidQuantity
		}

		if item.QuantityReserve > current.QuantityReserve {
			return ErrInsufficientReserved
		}

		if item.QuantityReserve > current.QuantityTotal {
			return ErrInsufficientStock
		}

		current.QuantityReserve -= item.QuantityReserve
		current.QuantityTotal -= item.QuantityReserve

		if err = h.repo.Update(ctx, &current); err != nil {
			return err
		}
	}

	return nil
}

func (h *PartsSupplyUseCase) GetPartsSupplyByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error) {
	partsSupplies, err := h.repo.GetByServiceOrderID(ctx, serviceOrderID)
	if err != nil {
		return nil, err
	}

	if len(partsSupplies) == 0 {
		return nil, ErrPartsSupplyNotFound
	}

	return partsSupplies, nil
}

func (h *PartsSupplyUseCase) GetPartsSupplyByID(ctx context.Context, id uint) (entities.PartsSupply, error) {
	foundPartsSupply, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return entities.PartsSupply{}, err
	}

	if foundPartsSupply.ID == 0 {
		return entities.PartsSupply{}, ErrPartsSupplyNotFound
	}

	return foundPartsSupply, nil
}

func (h *PartsSupplyUseCase) CreatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) (entities.PartsSupply, error) {
	//valide se jÃ¡ existe uma peÃ§a com o mesmo nome
	existingPartsSupply, err := h.repo.GetByName(ctx, partsSupply.Name)
	if err == nil && existingPartsSupply.ID != 0 {
		return entities.PartsSupply{}, ErrPartsSupplyAlreadyExists
	}

	return h.repo.Create(ctx, partsSupply)

}

func (h *PartsSupplyUseCase) UpdatePartsSupply(ctx context.Context, partsSupply *entities.PartsSupply) error {

	existingPartsSupply, err := h.repo.GetByID(ctx, partsSupply.ID)
	if err != nil {
		return errors.New("failed to retrieve parts supply")
	}

	if existingPartsSupply.ID == 0 {
		return ErrPartsSupplyNotFound
	}

	return h.repo.Update(ctx, partsSupply)

}

func (h *PartsSupplyUseCase) DeletePartsSupply(ctx context.Context, id uint) error {
	service, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if service.ID == 0 {
		return ErrPartsSupplyNotFound
	}
	return h.repo.Delete(ctx, id)
}

func (h *PartsSupplyUseCase) ListPartsSupplies(ctx context.Context) ([]entities.PartsSupply, error) {
	return h.repo.List(ctx)
}
