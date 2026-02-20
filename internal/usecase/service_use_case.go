package usecase

import (
	"context"
	"errors"
	"fmt"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/usecase/interfaces"
)

type IServiceUseCase interface {
	GetServiceByID(ctx context.Context, id uint) (entities.Service, error)
	CreateService(ctx context.Context, service *entities.Service) (entities.Service, error)
	UpdateService(ctx context.Context, service *entities.Service) error
	DeleteService(ctx context.Context, id uint) error
	ListServices(ctx context.Context) ([]entities.Service, error)
}
type ServiceUseCase struct {
	repo interfaces.IServiceRepo
}

var _ IServiceUseCase = (*ServiceUseCase)(nil)

func NewServiceUseCase(repo interfaces.IServiceRepo) *ServiceUseCase {
	return &ServiceUseCase{repo: repo}
}

var (
	ErrServiceNotFound      = errors.New("service not found")
	ErrServiceAlreadyExists = errors.New("service already exists")
)

func (h *ServiceUseCase) GetServiceByID(ctx context.Context, id uint) (entities.Service, error) {
	service, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return entities.Service{}, err
	}
	if service.ID == 0 {
		return entities.Service{}, ErrServiceNotFound
	}
	return service, nil
}

func (h *ServiceUseCase) CreateService(ctx context.Context, service *entities.Service) (entities.Service, error) {
	serviceFound, err := h.repo.GetByName(ctx, service.Name)
	if err != nil {
		return entities.Service{}, err
	}
	if serviceFound.ID != 0 {
		return entities.Service{}, ErrServiceAlreadyExists
	}

	return h.repo.Create(ctx, service)
}

func (h *ServiceUseCase) UpdateService(ctx context.Context, service *entities.Service) error {
	existingService, err := h.repo.GetByID(ctx, service.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve service: %w", err)
	}
	if existingService.ID == 0 {
		return ErrServiceNotFound
	}

	return h.repo.Update(ctx, service)
}

func (h *ServiceUseCase) DeleteService(ctx context.Context, id uint) error {
	service, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if service.ID == 0 {
		return ErrServiceNotFound
	}
	return h.repo.Delete(ctx, id)
}

func (h *ServiceUseCase) ListServices(ctx context.Context) ([]entities.Service, error) {
	return h.repo.List(ctx)
}
