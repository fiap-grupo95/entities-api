package interfaces

import (
	"context"
	"mecanica_xpto/internal/domain/entities"
)

type IServiceRepo interface {
	Create(ctx context.Context, so *entities.Service) (entities.Service, error)
	GetByID(ctx context.Context, id uint) (entities.Service, error)
	GetByName(ctx context.Context, name string) (entities.Service, error)
	Update(ctx context.Context, so *entities.Service) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]entities.Service, error)
}
