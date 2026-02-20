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
	handlerpkg "mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/internal/adapter/http/handlers/mocks"
	"mecanica_xpto/internal/domain/entities"
	usecase "mecanica_xpto/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupPaymentRouter(handler *handlerpkg.PaymentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/payments/:id", handler.GetPaymentByID)
	r.GET("/payments", handler.ListPayments)
	r.POST("/payments", handler.CreatePayment)
	return r
}

func TestGetPaymentByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	handler := handlerpkg.NewPaymentHandler(mockUC)
	router := setupPaymentRouter(handler)

	expected := &entities.Payment{
		ID:             1,
		ServiceOrderID: 10,
		PaymentDate:    time.Now(),
		Amount:         500,
	}
	mockUC.EXPECT().GetPaymentByID(gomock.Any(), uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/payments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp response.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.ServiceOrderID, resp.ServiceOrderID)
	assert.Equal(t, expected.Amount, resp.Amount)
}

func TestGetPaymentByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	handler := handlerpkg.NewPaymentHandler(mockUC)
	router := setupPaymentRouter(handler)

	mockUC.EXPECT().GetPaymentByID(gomock.Any(), uint(1)).Return(&entities.Payment{}, usecase.ErrorPaymentNotFound)

	req, _ := http.NewRequest("GET", "/payments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPaymentByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := handlerpkg.NewPaymentHandler(mocks.NewMockIPaymentUseCase(ctrl))
	router := setupPaymentRouter(handler)

	req, _ := http.NewRequest("GET", "/payments/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListPayments_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	handler := handlerpkg.NewPaymentHandler(mockUC)
	router := setupPaymentRouter(handler)

	now := time.Now()
	expected := []entities.Payment{
		{ID: 1, ServiceOrderID: 10, PaymentDate: now, Amount: 100},
		{ID: 2, ServiceOrderID: 11, PaymentDate: now, Amount: 200},
	}
	mockUC.EXPECT().ListPayments(gomock.Any()).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []response.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, expected[0].Amount, resp[0].Amount)
}

func TestListPayments_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	handler := handlerpkg.NewPaymentHandler(mockUC)
	router := setupPaymentRouter(handler)

	mockUC.EXPECT().ListPayments(gomock.Any()).Return(nil, errors.New("fail"))

	req, _ := http.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreatePayment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	handler := handlerpkg.NewPaymentHandler(mockUC)
	router := setupPaymentRouter(handler)

	now := time.Now().UTC().Truncate(time.Second)
	payload := request.PaymentCreateRequest{
		ServiceOrderID: 10,
		PaymentDate:    now.Format(time.RFC3339),
		Amount:         500.50,
	}
	body, _ := json.Marshal(payload)

	expected := &entities.Payment{
		ID:             1,
		ServiceOrderID: payload.ServiceOrderID,
		PaymentDate:    now,
		Amount:         payload.Amount,
	}

	mockUC.EXPECT().
		CreatePayment(gomock.Any(), gomock.AssignableToTypeOf(&entities.Payment{})).
		DoAndReturn(func(_ interface{}, payment *entities.Payment) (*entities.Payment, error) {
			assert.Equal(t, payload.ServiceOrderID, payment.ServiceOrderID)
			assert.Equal(t, payload.Amount, payment.Amount)
			assert.True(t, payment.PaymentDate.Equal(now))
			return expected, nil
		})

	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp response.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Amount, resp.Amount)
}

func TestCreatePayment_InvalidPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler := handlerpkg.NewPaymentHandler(mocks.NewMockIPaymentUseCase(ctrl))
	router := setupPaymentRouter(handler)

	req, _ := http.NewRequest("POST", "/payments", bytes.NewBufferString(`{"service_order_id":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePayment_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	handler := handlerpkg.NewPaymentHandler(mockUC)
	router := setupPaymentRouter(handler)

	now := time.Now().UTC().Truncate(time.Second)
	payload := request.PaymentCreateRequest{
		ServiceOrderID: 10,
		PaymentDate:    now.Format(time.RFC3339),
		Amount:         500.50,
	}
	body, _ := json.Marshal(payload)

	mockUC.EXPECT().
		CreatePayment(gomock.Any(), gomock.AssignableToTypeOf(&entities.Payment{})).
		Return(nil, usecase.ErrPaymentAmountDoesNotMatch)

	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockUC.EXPECT().
		CreatePayment(gomock.Any(), gomock.AssignableToTypeOf(&entities.Payment{})).
		Return(nil, usecase.ErrPaymentAlreadyExists)

	req, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	mockUC.EXPECT().
		CreatePayment(gomock.Any(), gomock.AssignableToTypeOf(&entities.Payment{})).
		Return(nil, errors.New("fail"))

	req, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
