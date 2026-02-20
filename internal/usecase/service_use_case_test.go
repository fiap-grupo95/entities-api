package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/usecase/mocks"

	"github.com/golang/mock/gomock"
)

func TestGetServiceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()

	// Sucesso
	service := entities.Service{ID: 1, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(service, nil)
	result, err := uc.GetServiceByID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, service) {
		t.Errorf("expected %v, got %v", service, result)
	}

	// NÃƒÆ’Ã‚Â£o encontrado
	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.Service{}, nil)
	_, err = uc.GetServiceByID(ctx, 2)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}

	// Erro do repo
	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.Service{}, errors.New("fail"))
	_, err = uc.GetServiceByID(ctx, 3)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()

	// JÃƒÆ’Ã‚Â¡ existe
	service := &entities.Service{ID: 1, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}
	mockRepo.EXPECT().GetByName(ctx, service.Name).Return(entities.Service{ID: 2, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}, nil)
	_, err := uc.CreateService(ctx, service)
	if !errors.Is(err, ErrServiceAlreadyExists) {
		t.Errorf("expected ErrServiceAlreadyExists, got %v", err)
	}

	// Erro ao buscar por nome
	mockRepo.EXPECT().GetByName(ctx, service.Name).Return(entities.Service{}, errors.New("fail"))
	_, err = uc.CreateService(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Sucesso
	mockRepo.EXPECT().GetByName(ctx, service.Name).Return(entities.Service{}, nil)
	mockRepo.EXPECT().Create(ctx, service).Return(*service, nil)
	result, err := uc.CreateService(ctx, service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, *service) {
		t.Errorf("expected %v, got %v", *service, result)
	}

	// Erro ao criar
	mockRepo.EXPECT().GetByName(ctx, service.Name).Return(entities.Service{}, nil)
	mockRepo.EXPECT().Create(ctx, service).Return(entities.Service{}, errors.New("fail"))
	_, err = uc.CreateService(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()
	// Erro ao buscar por ID
	service := &entities.Service{ID: 1, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}
	mockRepo.EXPECT().GetByID(ctx, service.ID).Return(entities.Service{}, errors.New("fail"))
	err := uc.UpdateService(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// NÃƒÆ’Ã‚Â£o encontrado
	mockRepo.EXPECT().GetByID(ctx, service.ID).Return(entities.Service{}, nil)
	err = uc.UpdateService(ctx, service)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}

	// Sucesso
	mockRepo.EXPECT().GetByID(ctx, service.ID).Return(*service, nil)
	mockRepo.EXPECT().Update(ctx, service).Return(nil)
	err = uc.UpdateService(ctx, service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Erro ao atualizar
	mockRepo.EXPECT().GetByID(ctx, service.ID).Return(*service, nil)
	mockRepo.EXPECT().Update(ctx, service).Return(errors.New("fail"))
	err = uc.UpdateService(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()

	// Erro ao buscar por ID
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(entities.Service{}, errors.New("fail"))
	err := uc.DeleteService(ctx, 1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// NÃƒÆ’Ã‚Â£o encontrado
	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.Service{}, nil)
	err = uc.DeleteService(ctx, 2)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}

	// Sucesso
	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.Service{ID: 3, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}, nil)
	mockRepo.EXPECT().Delete(ctx, uint(3)).Return(nil)
	err = uc.DeleteService(ctx, 3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Erro ao deletar
	mockRepo.EXPECT().GetByID(ctx, uint(4)).Return(entities.Service{ID: 4, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}, nil)
	mockRepo.EXPECT().Delete(ctx, uint(4)).Return(errors.New("fail"))
	err = uc.DeleteService(ctx, 4)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()
	services := []entities.Service{{ID: 1, Name: "Troca de ÃƒÆ’Ã‚Â³leo"}, {ID: 2, Name: "Alinhamento"}}

	mockRepo.EXPECT().List(ctx).Return(services, nil)
	result, err := uc.ListServices(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, services) {
		t.Errorf("expected %v, got %v", services, result)
	}

	mockRepo.EXPECT().List(ctx).Return(nil, errors.New("fail"))
	_, err = uc.ListServices(ctx)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
