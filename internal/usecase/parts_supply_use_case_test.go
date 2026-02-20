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

func TestGetPartsSupplyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()

	ps := entities.PartsSupply{ID: 1, Name: "Filtro"}
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(ps, nil)
	result, err := uc.GetPartsSupplyByID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, ps) {
		t.Errorf("expected %v, got %v", ps, result)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.PartsSupply{}, nil)
	_, err = uc.GetPartsSupplyByID(ctx, 2)
	if !errors.Is(err, ErrPartsSupplyNotFound) {
		t.Errorf("expected ErrPartsSupplyNotFound, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.PartsSupply{}, errors.New("fail"))
	_, err = uc.GetPartsSupplyByID(ctx, 3)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCreatePartsSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()

	// Test: jÃƒÆ’Ã‚Â¡ existe
	ps := &entities.PartsSupply{ID: 1, Name: "Filtro"}
	mockRepo.EXPECT().GetByName(ctx, "Filtro").Return(entities.PartsSupply{ID: 2, Name: "Filtro"}, nil)
	_, err := uc.CreatePartsSupply(ctx, ps)
	if !errors.Is(err, ErrPartsSupplyAlreadyExists) {
		t.Errorf("expected ErrPartsSupplyAlreadyExists, got %v", err)
	}

	// Test: erro ao buscar por nome (deve criar normalmente)
	ps2 := &entities.PartsSupply{ID: 3, Name: "Novo"}
	mockRepo.EXPECT().GetByName(ctx, "Novo").Return(entities.PartsSupply{}, errors.New("not found"))
	mockRepo.EXPECT().Create(ctx, ps2).Return(*ps2, nil)
	_, err = uc.CreatePartsSupply(ctx, ps2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test: erro ao criar
	ps3 := &entities.PartsSupply{ID: 4, Name: "Falha"}
	mockRepo.EXPECT().GetByName(ctx, "Falha").Return(entities.PartsSupply{}, errors.New("not found"))
	mockRepo.EXPECT().Create(ctx, ps3).Return(entities.PartsSupply{}, errors.New("fail"))
	_, err = uc.CreatePartsSupply(ctx, ps3)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdatePartsSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()

	// Test: erro ao buscar por ID
	ps := &entities.PartsSupply{ID: 1, Name: "Filtro"}
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(entities.PartsSupply{}, errors.New("fail"))
	err := uc.UpdatePartsSupply(ctx, ps)
	if err == nil || err.Error() != "failed to retrieve parts supply" {
		t.Errorf("expected failed to retrieve parts supply, got %v", err)
	}

	// Test: nÃƒÆ’Ã‚Â£o encontrado
	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.PartsSupply{}, nil)
	ps.ID = 2
	err = uc.UpdatePartsSupply(ctx, ps)
	if !errors.Is(err, ErrPartsSupplyNotFound) {
		t.Errorf("expected ErrPartsSupplyNotFound, got %v", err)
	}

	// Test: sucesso
	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.PartsSupply{ID: 3, Name: "Filtro"}, nil)
	mockRepo.EXPECT().Update(ctx, ps).Return(nil)
	ps.ID = 3
	err = uc.UpdatePartsSupply(ctx, ps)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test: erro ao atualizar
	mockRepo.EXPECT().GetByID(ctx, uint(4)).Return(entities.PartsSupply{ID: 4, Name: "Falha"}, nil)
	mockRepo.EXPECT().Update(ctx, ps).Return(errors.New("fail"))
	ps.ID = 4
	err = uc.UpdatePartsSupply(ctx, ps)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeletePartsSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()

	// Test: erro ao buscar por ID
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(entities.PartsSupply{}, errors.New("fail"))
	err := uc.DeletePartsSupply(ctx, 1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Test: nÃƒÆ’Ã‚Â£o encontrado
	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.PartsSupply{}, nil)
	err = uc.DeletePartsSupply(ctx, 2)
	if !errors.Is(err, ErrPartsSupplyNotFound) {
		t.Errorf("expected ErrPartsSupplyNotFound, got %v", err)
	}

	// Test: sucesso
	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.PartsSupply{ID: 3, Name: "Filtro"}, nil)
	mockRepo.EXPECT().Delete(ctx, uint(3)).Return(nil)
	err = uc.DeletePartsSupply(ctx, 3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test: erro ao deletar
	mockRepo.EXPECT().GetByID(ctx, uint(4)).Return(entities.PartsSupply{ID: 4, Name: "Falha"}, nil)
	mockRepo.EXPECT().Delete(ctx, uint(4)).Return(errors.New("fail"))
	err = uc.DeletePartsSupply(ctx, 4)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListPartsSupplies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()
	parts := []entities.PartsSupply{{ID: 1, Name: "Filtro"}, {ID: 2, Name: "Pastilha"}}

	mockRepo.EXPECT().List(ctx).Return(parts, nil)
	result, err := uc.ListPartsSupplies(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, parts) {
		t.Errorf("expected %v, got %v", parts, result)
	}

	mockRepo.EXPECT().List(ctx).Return(nil, errors.New("fail"))
	_, err = uc.ListPartsSupplies(ctx)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
