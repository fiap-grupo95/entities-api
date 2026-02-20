//go:build legacy
// +build legacy

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"mecanica_xpto/internal/domain/entities"
	dto "mecanica_xpto/internal/infrastructure/database/model"
	"mecanica_xpto/internal/usecase"
	"mecanica_xpto/internal/usecase/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdditionalRepair_ServiceOrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdditionalRepairRepository(ctrl)
	mockRepoOS := mocks.NewMockIServiceOrderRepository(ctrl)
	mockServiceRepo := mocks.NewMockIServiceRepo(ctrl)
	mockPartsSupplyRepo := mocks.NewMockIPartsSupplyRepo(ctrl)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: 2}

	mockRepoOS.EXPECT().GetByID(uint(2)).Return(nil, errors.New("not found"))

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, entities.AdditionalRepair{}, result)
}

func TestCreateAdditionalRepair_ServiceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdditionalRepairRepository(ctrl)
	mockRepoOS := mocks.NewMockIServiceOrderRepository(ctrl)
	mockServiceRepo := mocks.NewMockIServiceRepo(ctrl)
	mockPartsSupplyRepo := mocks.NewMockIPartsSupplyRepo(ctrl)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: 1,
		Services:       []entities.Service{{ID: 99}},
	}

	mockRepoOS.EXPECT().GetByID(uint(1)).Return((&dto.ServiceOrderModel{ID: 1}).ToDomain(), nil)
	mockServiceRepo.EXPECT().GetByID(gomock.Any(), uint(99)).Return(entities.Service{}, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, entities.AdditionalRepair{}, result)
}

func TestCreateAdditionalRepair_PartsSupplyNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdditionalRepairRepository(ctrl)
	mockRepoOS := mocks.NewMockIServiceOrderRepository(ctrl)
	mockServiceRepo := mocks.NewMockIServiceRepo(ctrl)
	mockPartsSupplyRepo := mocks.NewMockIPartsSupplyRepo(ctrl)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: 1,
		PartsSupplies:  []entities.PartsSupply{{ID: 88}},
	}

	mockRepoOS.EXPECT().GetByID(uint(1)).Return((&dto.ServiceOrderModel{ID: 1}).ToDomain(), nil)
	mockPartsSupplyRepo.EXPECT().GetByID(gomock.Any(), uint(88)).Return(entities.PartsSupply{}, nil)

	result, err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
	assert.Equal(t, entities.AdditionalRepair{}, result)
}
