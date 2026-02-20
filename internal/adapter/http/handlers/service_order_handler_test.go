//go:build legacy
// +build legacy

package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mecanica_xpto/internal/adapter/http/handlers/mocks"
	"mecanica_xpto/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setupServiceOrderHandlerTest(t *testing.T) (*mocks.MockIServiceOrderUseCase, *ServiceOrderHandler, *gin.Engine) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockIServiceOrderUseCase(ctrl)
	h := &ServiceOrderHandler{serviceOrderUseCase: mockUC}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return mockUC, h, r
}

func TestCreateServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.POST("/os", h.CreateServiceOrder)

	// Success
	mockUC.EXPECT().CreateServiceOrder(gomock.Any(), gomock.Any()).Return(&entities.ServiceOrder{}, nil)
	jsonBody := `{"customer_id":1,"vehicle_id":1}`
	req, _ := http.NewRequest("POST", "/os", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	// Invalid input
	req, _ = http.NewRequest("POST", "/os", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	// Failure
	mockUC.EXPECT().CreateServiceOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))
	req, _ = http.NewRequest("POST", "/os", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetServiceOrder(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os/:id", h.GetServiceOrder)

	mockUC.EXPECT().GetServiceOrder(gomock.Any(), entities.ServiceOrder{ID: 1}).Return(&entities.ServiceOrder{ID: 1}, nil)
	req, _ := http.NewRequest("GET", "/os/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().GetServiceOrder(gomock.Any(), entities.ServiceOrder{ID: 2}).Return(nil, errors.New("fail"))
	req, _ = http.NewRequest("GET", "/os/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = http.NewRequest("GET", "/os/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListServiceOrders(t *testing.T) {
	mockUC, h, r := setupServiceOrderHandlerTest(t)
	r.GET("/os", h.ListServiceOrders)
	mockUC.EXPECT().ListServiceOrders(gomock.Any()).Return([]*entities.ServiceOrder{{ID: 1}}, nil)
	req, _ := http.NewRequest("GET", "/os", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().ListServiceOrders(gomock.Any()).Return(nil, errors.New("fail"))
	req, _ = http.NewRequest("GET", "/os", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
