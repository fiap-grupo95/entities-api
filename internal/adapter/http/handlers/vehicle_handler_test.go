package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	handlerpkg "mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/internal/adapter/http/handlers/mocks"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"mecanica_xpto/internal/usecase"
	"mecanica_xpto/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newVehicleRouter() (*mocks.MockVehicleService, *gin.Engine) {
	mockService := new(mocks.MockVehicleService)
	handler := handlerpkg.NewVehicleHandler(mockService)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/vehicles", handler.GetVehicles)
	r.GET("/vehicles/:id", handler.GetVehicleByID)
	r.GET("/vehicles/customer/:customerID", handler.GetVehicleByCustomerID)
	r.GET("/vehicles/plate/:plate", handler.GetVehicleByPlate)
	r.POST("/vehicles", handler.CreateVehicle)
	r.PATCH("/vehicles/:id", handler.UpdateVehicle)
	r.DELETE("/vehicles/:id", handler.DeleteVehicle)
	return mockService, r
}

func TestVehicleHandler_GetVehicles(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	vehicles := []entities.Vehicle{
		{ID: 1, CustomerID: 10, Brand: "Toyota", Model: "Corolla", Year: "2020", Plate: valueobject.ParsePlate("ABC1D23")},
		{ID: 2, CustomerID: 11, Brand: "Honda", Model: "Civic", Year: "2021", Plate: valueobject.ParsePlate("XYZ5678")},
	}
	mockService.On("GetAllVehicles").Return(vehicles, nil).Once()

	req, _ := http.NewRequest("GET", "/vehicles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []response.VehicleResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 2)

	mockService.On("GetAllVehicles").Return(nil, errors.New("fail")).Once()
	req, _ = http.NewRequest("GET", "/vehicles", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVehicleHandler_GetVehicleByID(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	expected := &entities.Vehicle{ID: 1, CustomerID: 10, Brand: "Toyota", Model: "Corolla", Year: "2020", Plate: valueobject.ParsePlate("ABC1D23")}
	mockService.On("GetVehicleByID", uint(1)).Return(expected, nil).Once()

	req, _ := http.NewRequest("GET", "/vehicles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockService.On("GetVehicleByID", uint(1)).Return(nil, usecase.ErrVehicleNotFound).Once()
	req, _ = http.NewRequest("GET", "/vehicles/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("GET", "/vehicles/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVehicleHandler_GetVehicleByCustomerID(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	doc, err := valueobject.NewCpfCnpj("52998224725")
	require.NoError(t, err)

	customer := &entities.Customer{
		ID:          10,
		FullName:    "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "123456789",
		CpfCnpj:     doc,
	}

	vehicles := []entities.Vehicle{
		{
			ID:         1,
			CustomerID: customer.ID,
			Brand:      "Toyota",
			Model:      "Corolla",
			Year:       "2020",
			Plate:      valueobject.ParsePlate("ABC1D23"),
			Customer:   customer,
		},
	}
	mockService.On("GetVehiclesByCustomerID", uint(10)).Return(vehicles, nil).Once()

	req, _ := http.NewRequest("GET", "/vehicles/customer/10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.VehiclesByCustomerResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.NotNil(t, resp.Customer)
	assert.Equal(t, customer.ID, resp.Customer.ID)
	assert.Equal(t, customer.FullName, resp.Customer.FullName)
	assert.Len(t, resp.Vehicles, 1)
	assert.Equal(t, vehicles[0].ID, resp.Vehicles[0].ID)

	req, _ = http.NewRequest("GET", "/vehicles/customer/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVehicleHandler_GetVehicleByPlate(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	expected := &entities.Vehicle{ID: 1, Brand: "Toyota", Model: "Corolla", Year: "2020", Plate: valueobject.ParsePlate("ABC1D23")}
	mockService.On("GetVehicleByPlate", "ABC1D23").Return(expected, nil).Once()

	req, _ := http.NewRequest("GET", "/vehicles/plate/ABC1D23", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Run("empty plate", func(t *testing.T) {
		h := handlerpkg.NewVehicleHandler(new(mocks.MockVehicleService))
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: ""}}

		h.GetVehicleByPlate(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestVehicleHandler_CreateVehicle(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	payload := request.VehicleCreateRequest{
		CustomerID: 10,
		Brand:      "Toyota",
		Model:      "Corolla",
		Year:       "2020",
		Plate:      "ABC1D23",
	}
	body, _ := json.Marshal(payload)

	createdVehicle := &entities.Vehicle{
		ID:         1,
		CustomerID: payload.CustomerID,
		Brand:      payload.Brand,
		Model:      payload.Model,
		Year:       payload.Year,
		Plate:      valueobject.ParsePlate(payload.Plate),
	}

	mockService.
		On("CreateVehicle", mock.MatchedBy(func(v entities.Vehicle) bool {
			return v.CustomerID == payload.CustomerID &&
				v.Brand == payload.Brand &&
				v.Model == payload.Model &&
				v.Year == payload.Year &&
				v.Plate.String() == payload.Plate
		})).
		Return(createdVehicle, nil).
		Once()

	req, _ := http.NewRequest("POST", "/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.VehicleResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, createdVehicle.ID, resp.ID)
	assert.Equal(t, createdVehicle.CustomerID, resp.CustomerID)
	assert.Equal(t, createdVehicle.Brand, resp.Brand)
	assert.Equal(t, createdVehicle.Model, resp.Model)
	assert.Equal(t, createdVehicle.Year, resp.Year)
	assert.Equal(t, createdVehicle.Plate.String(), resp.Plate)

	req, _ = http.NewRequest("POST", "/vehicles", bytes.NewBufferString(`{"customer_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.
		On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).
		Return((*entities.Vehicle)(nil), usecase.ErrVehicleAlreadyExists).
		Once()
	req, _ = http.NewRequest("POST", "/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestVehicleHandler_UpdateVehicle(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	payload := request.VehicleUpdateRequest{
		Brand: mockStringPtr("Toyota"),
		Plate: mockStringPtr("ABC1D23"),
	}
	body, _ := json.Marshal(payload)
	expectedMap := map[string]interface{}{
		"brand": "Toyota",
		"plate": "ABC1D23",
	}

	mockService.
		On("UpdateVehiclePartial", uint(1), expectedMap).
		Return(usecase.MessageVehicleUpdatedSuccessfully, nil).
		Once()

	req, _ := http.NewRequest("PATCH", "/vehicles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest("PATCH", "/vehicles/1", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("PATCH", "/vehicles/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.
		On("UpdateVehiclePartial", uint(1), expectedMap).
		Return("", usecase.ErrVehicleNotFound).
		Once()
	req, _ = http.NewRequest("PATCH", "/vehicles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVehicleHandler_DeleteVehicle(t *testing.T) {
	mockService, router := newVehicleRouter()
	defer mockService.AssertExpectations(t)

	mockService.On("DeleteVehicle", uint(1)).Return(nil).Once()
	req, _ := http.NewRequest("DELETE", "/vehicles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	req, _ = http.NewRequest("DELETE", "/vehicles/invalid", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockService.On("DeleteVehicle", uint(2)).Return(usecase.ErrVehicleNotFound).Once()
	req, _ = http.NewRequest("DELETE", "/vehicles/2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockService.On("DeleteVehicle", uint(3)).Return(pkg.NewDomainError("INTERNAL_ERROR", "fail", nil, http.StatusInternalServerError)).Once()
	req, _ = http.NewRequest("DELETE", "/vehicles/3", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func mockStringPtr(value string) *string {
	return &value
}
