//go:build legacy
// +build legacy

package handlers

import (
	"errors"
	"net/http"
	"strconv"

	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	diagnosisFlow = "diagnosis"
	estimateFlow  = "estimate"
	executionFlow = "execution"
	deliveryFlow  = "delivery"
)

type ServiceOrderHandler struct {
	serviceOrderUseCase usecase.IServiceOrderUseCase
}

func NewServiceOrderHandler(useCase usecase.IServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		serviceOrderUseCase: useCase,
	}
}

// CreateServiceOrder godoc
// @Summary Create a new service order
// @Description Create a new service order record
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param order body request.ServiceOrderCreateRequest true "Service Order Information"
// @Success 201 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders [post]
func (h *ServiceOrderHandler) CreateServiceOrder(c *gin.Context) {
	var req request.ServiceOrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.CreateServiceOrder(c.Request.Context(), req.ToEntity())
	if err != nil {
		writeServiceOrderError(c, err, "Failed to create service order")
		return
	}

	c.JSON(http.StatusCreated, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderDiagnosis godoc
// @Summary Update service order diagnosis
// @Description Update the diagnosis information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderDiagnosisUpdateRequest true "Service Order Diagnosis Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/diagnosis [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderDiagnosis(c *gin.Context) {
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		return
	}

	var req request.ServiceOrderDiagnosisUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(c.Request.Context(), req.ToEntity(id), diagnosisFlow)
	if err != nil {
		writeServiceOrderError(c, err, "Failed to update service order diagnosis")
		return
	}

	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderEstimate godoc
// @Summary Update service order estimate
// @Description Update the estimate information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderEstimateUpdateRequest true "Service Order Estimate Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/estimate [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderEstimate(c *gin.Context) {
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		return
	}

	var req request.ServiceOrderEstimateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(c.Request.Context(), req.ToEntity(id), estimateFlow)
	if err != nil {
		writeServiceOrderError(c, err, "Failed to update service order estimate")
		return
	}

	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderExecution godoc
// @Summary Update service order execution
// @Description Update the execution information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderExecutionUpdateRequest true "Service Order Execution Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/execution [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderExecution(c *gin.Context) {
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		return
	}

	var req request.ServiceOrderExecutionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(c.Request.Context(), req.ToEntity(id), executionFlow)
	if err != nil {
		writeServiceOrderError(c, err, "Failed to update service order execution")
		return
	}

	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// UpdateServiceOrderDelivery godoc
// @Summary Update service order delivery
// @Description Update the delivery information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body request.ServiceOrderDeliveryUpdateRequest true "Service Order Delivery Information"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/delivery [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderDelivery(c *gin.Context) {
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		return
	}

	var req request.ServiceOrderDeliveryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindingError(c, err)
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(c.Request.Context(), req.ToEntity(id), deliveryFlow)
	if err != nil {
		writeServiceOrderError(c, err, "Failed to update service order delivery")
		return
	}

	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// GetServiceOrder godoc
// @Summary Get service order by ID
// @Description Retrieve a service order by its ID
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Success 200 {object} response.ServiceOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id} [get]
func (h *ServiceOrderHandler) GetServiceOrder(c *gin.Context) {
	id, ok := parseServiceOrderIDParam(c)
	if !ok {
		return
	}

	serviceOrder := entities.ServiceOrder{ID: id}
	result, err := h.serviceOrderUseCase.GetServiceOrder(c.Request.Context(), serviceOrder)
	if err != nil {
		writeServiceOrderError(c, err, "Failed to retrieve service order")
		return
	}

	c.JSON(http.StatusOK, response.NewServiceOrderResponse(result))
}

// ListServiceOrders godoc
// @Summary List all service orders
// @Description Get a list of all service orders
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.ServiceOrderResponse
// @Failure 500 {object} map[string]string
// @Router /service-orders [get]
func (h *ServiceOrderHandler) ListServiceOrders(c *gin.Context) {
	serviceOrders, err := h.serviceOrderUseCase.ListServiceOrders(c.Request.Context())
	if err != nil {
		writeServiceOrderError(c, err, "Failed to retrieve service orders")
		return
	}

	c.JSON(http.StatusOK, response.NewServiceOrderListResponse(serviceOrders))
}

func parseServiceOrderIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service order ID"})
		return 0, false
	}
	return uint(id), true
}

func writeBindingError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Invalid input",
		"details": err.Error(),
	})
}

func writeServiceOrderError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, usecase.ErrServiceOrderNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Service order not found"})
	case errors.Is(err, usecase.ErrVehicleNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
	case errors.Is(err, usecase.ErrCustomerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
	case errors.Is(err, usecase.ErrServiceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
	case errors.Is(err, usecase.ErrPartsSupplyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "Parts supply not found"})
	case errors.Is(err, usecase.ErrInvalidCustomerID),
		errors.Is(err, usecase.ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case isBadRequestServiceOrderError(err):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Error().Err(err).Msg(fallbackMessage)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   fallbackMessage,
			"details": err.Error(),
		})
	}
}

func isBadRequestServiceOrderError(err error) bool {
	switch {
	case errors.Is(err, usecase.ErrInvalidTransitionStatusToDiagnosis),
		errors.Is(err, usecase.ErrInvalidTransitionStatusToEstimate),
		errors.Is(err, usecase.ErrInvalidTransitionStatusToExecution),
		errors.Is(err, usecase.ErrInvalidTransitionStatusToDelivery),
		errors.Is(err, usecase.ErrInvalidStatus),
		errors.Is(err, usecase.ErrInvalidFlow),
		errors.Is(err, usecase.ErrInsufficientPartsSupply):
		return true
	default:
		return false
	}
}
