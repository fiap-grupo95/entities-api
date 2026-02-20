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

func setupPartsSupplyRouter(handler *handlerpkg.PartsSupplyHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/parts-supplies/:id", handler.GetPartsSupplyByID)
	r.POST("/parts-supplies", handler.CreatePartsSupply)
	r.PUT("/parts-supplies/:id", handler.UpdatePartsSupply)
	r.DELETE("/parts-supplies/:id", handler.DeletePartsSupply)
	r.GET("/parts-supplies", handler.ListPartsSupplies)
	return r
}

func TestGetPartsSupplyByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	expected := entities.PartsSupply{
		ID:              1,
		Name:            "Filtro",
		Description:     "Filtro de óleo",
		Price:           12.5,
		QuantityTotal:   10,
		QuantityReserve: 2,
	}
	mockUC.EXPECT().GetPartsSupplyByID(gomock.Any(), uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/parts-supplies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.PartsSupplyResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Name, resp.Name)
	assert.Equal(t, expected.QuantityReserve, resp.QuantityReserve)
}

func TestGetPartsSupplyByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	mockUC.EXPECT().GetPartsSupplyByID(gomock.Any(), uint(1)).Return(entities.PartsSupply{}, usecase.ErrPartsSupplyNotFound)

	req, _ := http.NewRequest("GET", "/parts-supplies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPartsSupplyByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := handlerpkg.NewPartsSupplyHandler(mocks.NewMockIPartsSupplyUseCase(ctrl))
	router := setupPartsSupplyRouter(handler)

	req, _ := http.NewRequest("GET", "/parts-supplies/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePartsSupply_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	payload := request.PartsSupplyCreateRequest{
		Name:            "Filtro",
		Description:     "Filtro de óleo",
		Price:           12.5,
		QuantityTotal:   10,
		QuantityReserve: 2,
	}
	body, _ := json.Marshal(payload)
	expected := entities.PartsSupply{
		ID:              1,
		Name:            payload.Name,
		Description:     payload.Description,
		Price:           payload.Price,
		QuantityTotal:   payload.QuantityTotal,
		QuantityReserve: payload.QuantityReserve,
	}

	mockUC.EXPECT().
		CreatePartsSupply(gomock.Any(), gomock.AssignableToTypeOf(&entities.PartsSupply{})).
		DoAndReturn(func(_ interface{}, ps *entities.PartsSupply) (entities.PartsSupply, error) {
			assert.Equal(t, payload.Name, ps.Name)
			assert.Equal(t, payload.Description, ps.Description)
			assert.Equal(t, payload.Price, ps.Price)
			assert.Equal(t, payload.QuantityTotal, ps.QuantityTotal)
			assert.Equal(t, payload.QuantityReserve, ps.QuantityReserve)
			return expected, nil
		})

	req, _ := http.NewRequest("POST", "/parts-supplies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.PartsSupplyResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Name, resp.Name)
}

func TestCreatePartsSupply_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	payload := request.PartsSupplyCreateRequest{
		Name:            "Filtro",
		Description:     "Filtro de óleo",
		Price:           12.5,
		QuantityTotal:   10,
		QuantityReserve: 2,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		CreatePartsSupply(gomock.Any(), gomock.AssignableToTypeOf(&entities.PartsSupply{})).
		Return(entities.PartsSupply{}, usecase.ErrPartsSupplyAlreadyExists)

	req, _ := http.NewRequest("POST", "/parts-supplies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreatePartsSupply_InvalidPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := handlerpkg.NewPartsSupplyHandler(mocks.NewMockIPartsSupplyUseCase(ctrl))
	router := setupPartsSupplyRouter(handler)

	req, _ := http.NewRequest("POST", "/parts-supplies", bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePartsSupply_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	payload := request.PartsSupplyUpdateRequest{
		Name:            "Filtro Atualizado",
		Description:     "Filtro atualizado",
		Price:           15.0,
		QuantityTotal:   20,
		QuantityReserve: 5,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		UpdatePartsSupply(gomock.Any(), gomock.AssignableToTypeOf(&entities.PartsSupply{})).
		DoAndReturn(func(_ interface{}, ps *entities.PartsSupply) error {
			assert.Equal(t, uint(1), ps.ID)
			assert.Equal(t, payload.Name, ps.Name)
			assert.Equal(t, payload.Price, ps.Price)
			return nil
		})

	req, _ := http.NewRequest("PUT", "/parts-supplies/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.OperationMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Parts supply updated successfully", resp.Message)
}

func TestUpdatePartsSupply_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	payload := request.PartsSupplyUpdateRequest{
		Name:            "Filtro Atualizado",
		Description:     "Filtro atualizado",
		Price:           15.0,
		QuantityTotal:   20,
		QuantityReserve: 5,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		UpdatePartsSupply(gomock.Any(), gomock.AssignableToTypeOf(&entities.PartsSupply{})).
		Return(usecase.ErrPartsSupplyNotFound)

	req, _ := http.NewRequest("PUT", "/parts-supplies/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdatePartsSupply_InvalidPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := handlerpkg.NewPartsSupplyHandler(mocks.NewMockIPartsSupplyUseCase(ctrl))
	router := setupPartsSupplyRouter(handler)

	req, _ := http.NewRequest("PUT", "/parts-supplies/1", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeletePartsSupply_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	mockUC.EXPECT().DeletePartsSupply(gomock.Any(), uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/parts-supplies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeletePartsSupply_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	mockUC.EXPECT().DeletePartsSupply(gomock.Any(), uint(1)).Return(usecase.ErrPartsSupplyNotFound)

	req, _ := http.NewRequest("DELETE", "/parts-supplies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListPartsSupplies_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	expected := []entities.PartsSupply{
		{ID: 1, Name: "Filtro", Description: "Filtro de óleo", Price: 10, QuantityTotal: 5},
		{ID: 2, Name: "Pastilha", Description: "Pastilha de freio", Price: 30, QuantityTotal: 8},
	}
	mockUC.EXPECT().ListPartsSupplies(gomock.Any()).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/parts-supplies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []response.PartsSupplyResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, expected[0].Name, resp[0].Name)
}

func TestListPartsSupplies_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	handler := handlerpkg.NewPartsSupplyHandler(mockUC)
	router := setupPartsSupplyRouter(handler)

	mockUC.EXPECT().ListPartsSupplies(gomock.Any()).Return(nil, errors.New("fail"))

	req, _ := http.NewRequest("GET", "/parts-supplies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
