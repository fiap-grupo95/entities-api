//go:build legacy
// +build legacy

package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	handler "mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/internal/adapter/http/handlers/mocks"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupADRRouter(h *handler.AdditionalRepairHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/additional-repair", h.CreateAdditionalRepair)
	r.GET("/additional-repair/:id", h.GetAdditionalRepair)
	r.POST("/additional-repair/:id/part", h.AddPartSupplyAndService)
	r.DELETE("/additional-repair/:id/part", h.RemovePartSupplyAndService)
	r.POST("/additional-repair/:id/approval", h.CustomerApproval)
	return r
}

func TestCreateSOAdditionalRepair_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairCreateRequest{ServiceOrderID: 1, Description: "desc"}
	expected := entities.AdditionalRepair{ServiceOrderID: 1, Description: "desc"}
	created := entities.AdditionalRepair{
		ID:             10,
		ServiceOrderID: expected.ServiceOrderID,
		Description:    expected.Description,
		ARStatus:       valueobject.AdditionalRepairStatus("IN_ANALYSIS"),
	}
	mockUC.EXPECT().CreateAdditionalRepair(gomock.Any(), expected).Return(created, nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.AdditionalRepairResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, resp.ID)
	assert.Equal(t, created.ServiceOrderID, resp.ServiceOrderID)
	assert.Equal(t, created.Description, resp.Description)
	assert.Equal(t, created.ARStatus.String(), resp.Status)
}

func TestCreateSOAdditionalRepair_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairCreateRequest{ServiceOrderID: 1, Description: "desc"}
	expected := entities.AdditionalRepair{ServiceOrderID: 1, Description: "desc"}
	mockUC.EXPECT().CreateAdditionalRepair(gomock.Any(), expected).Return(entities.AdditionalRepair{}, errors.New("fail"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAdditionalRepair_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	now := time.Now()
	expected := entities.AdditionalRepair{
		ID:             1,
		ServiceOrderID: 10,
		Description:    "desc",
		ARStatus:       valueobject.StatusAAprovada,
		Estimate:       123.45,
		CreatedAt:      now,
		UpdatedAt:      now,
		Services: []entities.Service{
			{ID: 2, Name: "service", Price: 80},
		},
		PartsSupplies: []entities.PartsSupply{
			{ID: 3, Name: "part", Price: 43.45, QuantityReserve: 2, QuantityTotal: 10},
		},
	}
	mockUC.EXPECT().GetAdditionalRepair(gomock.Any(), uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/additional-repair/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.AdditionalRepairResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.ServiceOrderID, resp.ServiceOrderID)
	assert.Equal(t, expected.Description, resp.Description)
	assert.Equal(t, expected.ARStatus.String(), resp.Status)
	assert.Equal(t, expected.Estimate, resp.Estimate)
	assert.Len(t, resp.Services, 1)
	assert.Equal(t, expected.Services[0].ID, resp.Services[0].ID)
	assert.Len(t, resp.PartsSupplies, 1)
	assert.Equal(t, expected.PartsSupplies[0].ID, resp.PartsSupplies[0].ID)
}

func TestGetAdditionalRepair_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	// Return empty struct and error
	mockUC.EXPECT().GetAdditionalRepair(gomock.Any(), uint(999)).Return(entities.AdditionalRepair{}, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/additional-repair/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddPartSupplyAndService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairItemsRequest{
		Description: "update",
		Services:    []request.AdditionalRepairServiceItem{{ID: 4}},
		PartsSupplies: []request.AdditionalRepairPartsSupplyItem{
			{ID: 5, QuantityReserve: 3},
		},
	}
	expected := entities.AdditionalRepair{
		Description: "update",
		Services:    []entities.Service{{ID: 4}},
		PartsSupplies: []entities.PartsSupply{
			{ID: 5, QuantityReserve: 3},
		},
	}
	mockUC.EXPECT().AddPartSupplyAndService(gomock.Any(), uint(1), expected).Return(nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.OperationMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Additional repair updated successfully", resp.Message)
}

func TestAddPartSupplyAndService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairItemsRequest{
		Description: "update",
		Services:    []request.AdditionalRepairServiceItem{{ID: 4}},
		PartsSupplies: []request.AdditionalRepairPartsSupplyItem{
			{ID: 5, QuantityReserve: 3},
		},
	}
	expected := entities.AdditionalRepair{
		Description: "update",
		Services:    []entities.Service{{ID: 4}},
		PartsSupplies: []entities.PartsSupply{
			{ID: 5, QuantityReserve: 3},
		},
	}
	mockUC.EXPECT().AddPartSupplyAndService(gomock.Any(), uint(1), expected).Return(errors.New("fail"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRemovePartSupplyAndService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairItemsRequest{
		Description: "remove",
		Services:    []request.AdditionalRepairServiceItem{{ID: 4}},
		PartsSupplies: []request.AdditionalRepairPartsSupplyItem{
			{ID: 5, QuantityReserve: 1},
		},
	}
	expected := entities.AdditionalRepair{
		Description: "remove",
		Services:    []entities.Service{{ID: 4}},
		PartsSupplies: []entities.PartsSupply{
			{ID: 5, QuantityReserve: 1},
		},
	}
	mockUC.EXPECT().RemovePartSupplyAndService(gomock.Any(), uint(1), expected).Return(nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("DELETE", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.OperationMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Additional repair updated successfully", resp.Message)
}

func TestRemovePartSupplyAndService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairItemsRequest{
		Description: "remove",
		Services:    []request.AdditionalRepairServiceItem{{ID: 4}},
		PartsSupplies: []request.AdditionalRepairPartsSupplyItem{
			{ID: 5, QuantityReserve: 1},
		},
	}
	expected := entities.AdditionalRepair{
		Description: "remove",
		Services:    []entities.Service{{ID: 4}},
		PartsSupplies: []entities.PartsSupply{
			{ID: 5, QuantityReserve: 1},
		},
	}
	mockUC.EXPECT().RemovePartSupplyAndService(gomock.Any(), uint(1), expected).Return(errors.New("fail"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("DELETE", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerApproval_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairApprovalRequest{ApprovalStatus: "APPROVED"}
	expected := entities.AdditionalRepairStatusDTO{ApprovalStatus: "APPROVED"}
	mockUC.EXPECT().CustomerApprovalStatus(gomock.Any(), uint(1), expected).Return(nil)

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair/1/approval", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.OperationMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Additional repair updated successfully", resp.Message)
}

func TestCustomerApproval_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	payload := request.AdditionalRepairApprovalRequest{ApprovalStatus: "DENIED"}
	expected := entities.AdditionalRepairStatusDTO{ApprovalStatus: "DENIED"}
	mockUC.EXPECT().CustomerApprovalStatus(gomock.Any(), uint(1), expected).Return(errors.New("fail"))

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/additional-repair/1/approval", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
