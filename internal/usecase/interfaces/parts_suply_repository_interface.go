package interfaces

import (
	"context"
	"mecanica_xpto/internal/domain/entities"
)

type IPartsSupplyRepo interface {
	Create(ctx context.Context, ps *entities.PartsSupply) (entities.PartsSupply, error)
	GetByID(ctx context.Context, id uint) (entities.PartsSupply, error)
	GetByName(ctx context.Context, name string) (entities.PartsSupply, error)
	Update(ctx context.Context, ps *entities.PartsSupply) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]entities.PartsSupply, error)
	GetByServiceOrderID(ctx context.Context, serviceOrderID uint) ([]entities.PartsSupply, error)
}
