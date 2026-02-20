package usecase

import (
	"errors"
	"testing"

	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"mecanica_xpto/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVehicleService_GetAllVehicles(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	expected := []entities.Vehicle{
		{ID: 1, Plate: valueobject.ParsePlate("ABC1234")},
		{ID: 2, Plate: valueobject.ParsePlate("XYZ5678")},
	}

	mockRepo.On("FindAll").Return(expected, nil)

	result, err := service.GetAllVehicles()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetAllVehicles_Error(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	mockRepo.On("FindAll").Return(nil, errors.New("db error"))

	result, err := service.GetAllVehicles()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_GetVehicleByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := &entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByID", uint(1)).Return(vehicle, nil)

		result, err := service.GetVehicleByID(1)

		assert.NoError(t, err)
		assert.Equal(t, vehicle, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), nil)

		result, err := service.GetVehicleByID(1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrVehicleNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), errors.New("db error"))

		result, err := service.GetVehicleByID(1)

		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehicleByPlate(t *testing.T) {
	t.Run("invalid plate format", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		result, err := service.GetVehicleByPlate("INVALID")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrInvalidPlateFormat)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)
		vehicle := &entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByPlate", valueobject.ParsePlate("ABC1234")).Return(vehicle, nil)

		result, err := service.GetVehicleByPlate("ABC1234")

		assert.NoError(t, err)
		assert.Equal(t, vehicle, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByPlate", valueobject.ParsePlate("ABC1234")).Return((*entities.Vehicle)(nil), nil)

		result, err := service.GetVehicleByPlate("ABC1234")

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrVehicleNotFound)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehiclesByCustomerID(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	vehicles := []entities.Vehicle{{ID: 1}, {ID: 2}}
	mockRepo.On("FindByCustomerID", uint(10)).Return(vehicles, nil)

	result, err := service.GetVehiclesByCustomerID(10)

	assert.NoError(t, err)
	assert.Equal(t, vehicles, result)
	mockRepo.AssertExpectations(t)
}

func TestVehicleService_CreateVehicle(t *testing.T) {
	t.Run("invalid plate", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{Plate: valueobject.ParsePlate("INVALID")}

		result, err := service.CreateVehicle(vehicle)

		assert.ErrorIs(t, err, ErrInvalidPlateFormat)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("duplicate plate", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByPlate", vehicle.Plate).Return(&entities.Vehicle{ID: 1}, nil)

		result, err := service.CreateVehicle(vehicle)

		assert.ErrorIs(t, err, ErrVehicleAlreadyExists)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository find error", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByPlate", vehicle.Plate).Return((*entities.Vehicle)(nil), errors.New("db error"))

		result, err := service.CreateVehicle(vehicle)

		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("create success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			Plate:      valueobject.ParsePlate("ABC1234"),
			CustomerID: 1,
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
		}

		mockRepo.On("FindByPlate", vehicle.Plate).Return((*entities.Vehicle)(nil), nil)
		created := &entities.Vehicle{ID: 1, Plate: vehicle.Plate}
		mockRepo.On("Create", vehicle).Return(created, nil)

		result, err := service.CreateVehicle(vehicle)

		assert.NoError(t, err)
		assert.Equal(t, created, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_UpdateVehicle(t *testing.T) {
	t.Run("invalid plate", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{Plate: valueobject.ParsePlate("INVALID")}

		message, err := service.UpdateVehicle(vehicle)

		assert.ErrorIs(t, err, ErrInvalidPlateFormat)
		assert.Equal(t, MessageInvalidPlateFormat, message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), nil)

		message, err := service.UpdateVehicle(vehicle)

		assert.ErrorIs(t, err, ErrVehicleNotFound)
		assert.Equal(t, MessageVehicleNotFound, message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			ID:         1,
			Plate:      valueobject.ParsePlate("ABC1234"),
			Model:      "Civic",
			Brand:      "Honda",
			CustomerID: 1,
		}

		mockRepo.On("FindByID", uint(1)).Return(&vehicle, nil)
		mockRepo.On("Update", vehicle).Return(nil)

		message, err := service.UpdateVehicle(vehicle)

		assert.NoError(t, err)
		assert.Equal(t, MessageVehicleUpdatedSuccessfully, message)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_UpdateVehiclePartial(t *testing.T) {
	t.Run("vehicle not found", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), nil)

		message, err := service.UpdateVehiclePartial(1, map[string]interface{}{"model": "New"})

		assert.ErrorIs(t, err, ErrVehicleNotFound)
		assert.Equal(t, MessageVehicleNotFound, message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid plate", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		existing := &entities.Vehicle{ID: 1, Plate: valueobject.ParsePlate("ABC1234")}

		mockRepo.On("FindByID", uint(1)).Return(existing, nil)

		message, err := service.UpdateVehiclePartial(1, map[string]interface{}{"plate": "INVALID"})

		assert.ErrorIs(t, err, ErrInvalidPlateFormat)
		assert.Equal(t, MessageInvalidPlateFormat, message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		existing := &entities.Vehicle{
			ID:    1,
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Old",
		}

		mockRepo.On("FindByID", uint(1)).Return(existing, nil)
		mockRepo.On("Update", mock.MatchedBy(func(updated entities.Vehicle) bool {
			return updated.ID == 1 &&
				updated.Model == "New" &&
				updated.Plate.String() == "DEF5678"
		})).Return(nil)

		message, err := service.UpdateVehiclePartial(1, map[string]interface{}{
			"model": "New",
			"plate": "DEF5678",
		})

		assert.NoError(t, err)
		assert.Equal(t, MessageVehicleUpdatedSuccessfully, message)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_DeleteVehicle(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		err := service.DeleteVehicle(0)

		assert.ErrorIs(t, err, ErrInvalidID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return((*entities.Vehicle)(nil), nil)

		err := service.DeleteVehicle(1)

		assert.ErrorIs(t, err, ErrVehicleNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return(&entities.Vehicle{ID: 1}, nil)
		mockRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteVehicle(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
