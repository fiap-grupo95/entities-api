package usecase_test

import (
	"errors"
	use_cases "mecanica_xpto/internal/usecase"
	"mecanica_xpto/internal/usecase/interfaces/mocks"
	"testing"

	"mecanica_xpto/internal/domain/entities"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	customer := &entities.Customer{FullName: "Updated", CpfCnpj: "123", PhoneNumber: "999"}

	mockRepo.EXPECT().GetByID(uint(1)).Return(customer, nil)
	mockRepo.EXPECT().Update(customer).Return(nil)

	err := uc.UpdateCustomer(1, customer)
	assert.NoError(t, err)
}

func TestUpdateCustomer_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().GetByID(uint(1)).Return(nil, errors.New("not found"))

	err := uc.UpdateCustomer(1, &entities.Customer{})
	assert.Error(t, err)
}

func TestDeleteCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().Delete(uint(1)).Return(nil)

	err := uc.DeleteCustomer(1)
	assert.NoError(t, err)
}

func TestDeleteCustomer_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().Delete(uint(1)).Return(errors.New("fail"))

	err := uc.DeleteCustomer(1)
	assert.Error(t, err)
}
