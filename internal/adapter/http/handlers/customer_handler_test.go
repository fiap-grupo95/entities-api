package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	handlerpkg "mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/internal/adapter/http/handlers/mocks"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	usecase "mecanica_xpto/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *handlerpkg.CustomerHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/customers/:document", handler.GetCustomer)
	r.GET("/customers/id/:id", handler.GetFullCustomer)
	r.POST("/customers", handler.CreateCustomer)
	r.PUT("/customers/:id", handler.UpdateCustomer)
	r.DELETE("/customers/:id", handler.DeleteCustomer)
	r.GET("/customers", handler.ListCustomer)
	return r
}

func TestGetCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	doc, _ := valueobject.NewCpfCnpj("52998224725")
	customer := &entities.Customer{
		ID:          10,
		UserID:      5,
		FullName:    "Test",
		Email:       "test@example.com",
		PhoneNumber: "12345",
		CpfCnpj:     doc,
		Vehicles: []entities.Vehicle{
			{ID: 1, Brand: "Brand", Model: "Model", Plate: valueobject.ParsePlate("ABC1234"), Year: "2020"},
		},
	}
	mockUC.EXPECT().GetByDocument("123").Return(customer, nil)

	req, _ := http.NewRequest("GET", "/customers/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.CustomerResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, customer.ID, resp.ID)
	assert.Equal(t, customer.FullName, resp.FullName)
	assert.Equal(t, customer.Email, resp.Email)
	assert.Equal(t, customer.CpfCnpj.String(), resp.Document)
	assert.Len(t, resp.Vehicles, 1)
}

func TestGetCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().GetByDocument("123").Return(nil, usecase.ErrCustomerNotFound)

	req, _ := http.NewRequest("GET", "/customers/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	payload := request.CustomerCreateRequest{
		FullName:    "Test",
		Email:       "test@example.com",
		PhoneNumber: "12345",
		Document:    "52998224725",
	}
	body, _ := json.Marshal(payload)
	mockUC.EXPECT().
		CreateCustomer(gomock.AssignableToTypeOf(&entities.Customer{})).
		DoAndReturn(func(customer *entities.Customer) (*entities.Customer, error) {
			assert.Equal(t, payload.FullName, customer.FullName)
			assert.Equal(t, payload.Email, customer.Email)
			assert.Equal(t, payload.PhoneNumber, customer.PhoneNumber)
			assert.Equal(t, payload.Document, customer.CpfCnpj.String())
			customer.ID = 1
			return customer, nil
		})

	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.CustomerResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, payload.FullName, resp.FullName)
	assert.Equal(t, payload.Email, resp.Email)
	assert.Equal(t, payload.PhoneNumber, resp.PhoneNumber)
	assert.Equal(t, payload.Document, resp.Document)
}

func TestCreateCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	payload := request.CustomerCreateRequest{
		FullName:    "Test",
		Email:       "test@example.com",
		PhoneNumber: "12345",
		Document:    "52998224725",
	}
	body, _ := json.Marshal(payload)
	mockUC.EXPECT().
		CreateCustomer(gomock.AssignableToTypeOf(&entities.Customer{})).
		Return(nil, usecase.ErrCustomerAlreadyExists)

	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	payload := request.CustomerUpdateRequest{
		FullName:    "Test Updated",
		PhoneNumber: "9999",
	}
	body, _ := json.Marshal(payload)
	mockUC.EXPECT().
		UpdateCustomer(uint(1), gomock.AssignableToTypeOf(&entities.Customer{})).
		DoAndReturn(func(_ uint, customer *entities.Customer) error {
			assert.Equal(t, payload.FullName, customer.FullName)
			assert.Equal(t, payload.PhoneNumber, customer.PhoneNumber)
			return nil
		})

	req, _ := http.NewRequest("PUT", "/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.OperationMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Customer updated successfully", resp.Message)
}

func TestUpdateCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	payload := request.CustomerUpdateRequest{
		FullName: "Test Updated",
	}
	body, _ := json.Marshal(payload)
	mockUC.EXPECT().
		UpdateCustomer(uint(1), gomock.AssignableToTypeOf(&entities.Customer{})).
		Return(usecase.ErrCustomerNotFound)

	req, _ := http.NewRequest("PUT", "/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().DeleteCustomer(uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/customers/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().DeleteCustomer(uint(1)).Return(usecase.ErrCustomerNotFound)

	req, _ := http.NewRequest("DELETE", "/customers/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	doc, _ := valueobject.NewCpfCnpj("52998224725")
	customers := []entities.Customer{
		{
			ID:          1,
			FullName:    "John",
			Email:       "john@example.com",
			PhoneNumber: "1111",
			CpfCnpj:     doc,
		},
	}
	mockUC.EXPECT().ListCustomer().Return(customers, nil)

	req, _ := http.NewRequest("GET", "/customers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []response.CustomerResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, customers[0].FullName, resp[0].FullName)
}

func TestListCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := handlerpkg.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().ListCustomer().Return(nil, usecase.ErrGeneric)

	req, _ := http.NewRequest("GET", "/customers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
