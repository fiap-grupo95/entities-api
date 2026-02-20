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
	usecase "mecanica_xpto/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupServiceRouter(handler *handlerpkg.ServiceHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/services/:id", handler.GetServiceByID)
	r.POST("/services", handler.CreateService)
	r.PUT("/services/:id", handler.UpdateService)
	r.DELETE("/services/:id", handler.DeleteService)
	r.GET("/services", handler.ListServices)
	return r
}

func TestGetServiceByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	expected := entities.Service{ID: 1, Name: "Troca de Óleo", Description: "Troca completa", Price: 100}
	mockUC.EXPECT().GetServiceByID(gomock.Any(), uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/services/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.ServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Name, resp.Name)
	assert.Equal(t, expected.Price, resp.Price)
}

func TestGetServiceByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	mockUC.EXPECT().GetServiceByID(gomock.Any(), uint(1)).Return(entities.Service{}, usecase.ErrServiceNotFound)

	req, _ := http.NewRequest("GET", "/services/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	req, _ = http.NewRequest("GET", "/services/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	payload := request.ServiceCreateRequest{
		Name:        "Troca de Óleo",
		Description: "Troca completa",
		Price:       100,
	}
	body, _ := json.Marshal(payload)
	expected := entities.Service{ID: 1, Name: payload.Name, Description: payload.Description, Price: payload.Price}

	mockUC.EXPECT().
		CreateService(gomock.Any(), gomock.AssignableToTypeOf(&entities.Service{})).
		DoAndReturn(func(_ interface{}, service *entities.Service) (entities.Service, error) {
			assert.Equal(t, payload.Name, service.Name)
			assert.Equal(t, payload.Description, service.Description)
			assert.Equal(t, payload.Price, service.Price)
			return expected, nil
		})

	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.ServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Name, resp.Name)
}

func TestCreateService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	payload := request.ServiceCreateRequest{
		Name:        "Troca de Óleo",
		Description: "Troca completa",
		Price:       100,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		CreateService(gomock.Any(), gomock.AssignableToTypeOf(&entities.Service{})).
		Return(entities.Service{}, usecase.ErrServiceAlreadyExists)

	req, _ := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	mockUC.EXPECT().
		CreateService(gomock.Any(), gomock.AssignableToTypeOf(&entities.Service{})).
		Return(entities.Service{}, errors.New("fail"))

	req, _ = http.NewRequest("POST", "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	req, _ = http.NewRequest("POST", "/services", bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	payload := request.ServiceUpdateRequest{
		Name:        "Troca de Óleo - Premium",
		Description: "Troca completa premium",
		Price:       150,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		UpdateService(gomock.Any(), gomock.AssignableToTypeOf(&entities.Service{})).
		DoAndReturn(func(_ interface{}, service *entities.Service) error {
			assert.Equal(t, uint(1), service.ID)
			assert.Equal(t, payload.Name, service.Name)
			assert.Equal(t, payload.Price, service.Price)
			return nil
		})

	req, _ := http.NewRequest("PUT", "/services/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.OperationMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Service updated successfully", resp.Message)
}

func TestUpdateService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	payload := request.ServiceUpdateRequest{
		Name:        "Troca de Óleo - Premium",
		Description: "Troca completa premium",
		Price:       150,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		UpdateService(gomock.Any(), gomock.AssignableToTypeOf(&entities.Service{})).
		Return(usecase.ErrServiceNotFound)

	req, _ := http.NewRequest("PUT", "/services/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUC.EXPECT().
		UpdateService(gomock.Any(), gomock.AssignableToTypeOf(&entities.Service{})).
		Return(errors.New("fail"))

	req, _ = http.NewRequest("PUT", "/services/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	req, _ = http.NewRequest("PUT", "/services/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req, _ = http.NewRequest("PUT", "/services/1", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	mockUC.EXPECT().DeleteService(gomock.Any(), uint(1)).Return(nil)
	req, _ := http.NewRequest("DELETE", "/services/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockUC.EXPECT().DeleteService(gomock.Any(), uint(1)).Return(usecase.ErrServiceNotFound)
	req, _ = http.NewRequest("DELETE", "/services/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUC.EXPECT().DeleteService(gomock.Any(), uint(1)).Return(errors.New("fail"))
	req, _ = http.NewRequest("DELETE", "/services/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	req, _ = http.NewRequest("DELETE", "/services/abc", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	handler := handlerpkg.NewServiceHandler(mockUC)
	router := setupServiceRouter(handler)

	expected := []entities.Service{
		{ID: 1, Name: "Troca de Óleo", Description: "Troca completa", Price: 100},
		{ID: 2, Name: "Alinhamento", Description: "Alinhamento completo", Price: 150},
	}
	mockUC.EXPECT().ListServices(gomock.Any()).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/services", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp []response.ServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)

	mockUC.EXPECT().ListServices(gomock.Any()).Return(nil, errors.New("fail"))
	req, _ = http.NewRequest("GET", "/services", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
