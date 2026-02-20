package handlers

import (
	"errors"
	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	"mecanica_xpto/internal/domain/entities"
	usecase "mecanica_xpto/internal/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// Domain errors
	errInvalidPartsSupplyID    = pkg.NewDomainErrorSimple("INVALID_PARTS_SUPPLY_ID", "Invalid parts supply ID", http.StatusBadRequest)
	errInvalidPartsSupplyInput = pkg.NewDomainErrorSimple("INVALID_PARTS_SUPPLY_INPUT", "Invalid parts supply payload", http.StatusBadRequest)
)

// PartsSupplyHandler handles HTTP requests for parts supply operations
// @title Parts Supply API
// @version 1.0
// @description API for managing parts supply in the workshop management system
type PartsSupplyHandler struct {
	usecase usecase.IPartsSupplyUseCase
}

func NewPartsSupplyHandler(usecase usecase.IPartsSupplyUseCase) *PartsSupplyHandler {
	return &PartsSupplyHandler{usecase: usecase}
}

func mapPartsSupplyError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrPartsSupplyNotFound):
		return pkg.NewDomainErrorSimple("PARTS_SUPPLY_NOT_FOUND", "Parts supply not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid parts supply ID", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrInvalidQuantity):
		return pkg.NewDomainErrorSimple("INVALID_QUANTITY", "Invalid quantity", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrInsufficientStock):
		return pkg.NewDomainErrorSimple("INSUFFICIENT_STOCK", "Insufficient stock", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrInsufficientReserved):
		return pkg.NewDomainErrorSimple("INSUFFICIENT_RESERVED", "Insufficient reserved quantity", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrPartsSupplyAlreadyExists):
		return pkg.NewDomainErrorSimple("PARTS_SUPPLY_EXISTS", "Parts supply already exists", http.StatusConflict)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

func (h *PartsSupplyHandler) AuthorizeReserve(c *gin.Context) {
	items, ok := parseStockOperationPayload(c)
	if !ok {
		return
	}

	if err := h.usecase.AuthorizeReserve(c.Request.Context(), items); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Parts supply reserve authorized successfully"})
}

func (h *PartsSupplyHandler) Reserve(c *gin.Context) {
	items, ok := parseStockOperationPayload(c)
	if !ok {
		return
	}

	if err := h.usecase.Reserve(c.Request.Context(), items); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Parts supply reserved successfully"})
}

func (h *PartsSupplyHandler) Release(c *gin.Context) {
	items, ok := parseStockOperationPayload(c)
	if !ok {
		return
	}

	if err := h.usecase.Release(c.Request.Context(), items); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Parts supply released successfully"})
}

func (h *PartsSupplyHandler) WriteOff(c *gin.Context) {
	items, ok := parseStockOperationPayload(c)
	if !ok {
		return
	}

	if err := h.usecase.WriteOff(c.Request.Context(), items); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Parts supply writeoff completed successfully"})
}

// GetPartsSupplyByID godoc
// @Summary Get parts supply by ID
// @Description Retrieve a parts supply by its ID
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Parts Supply ID"
// @Success 200 {object} response.PartsSupplyResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies/{id} [get]
func (h *PartsSupplyHandler) GetPartsSupplyByID(c *gin.Context) {
	partsSupplyID, ok := parsePartsSupplyID(c)
	if !ok {
		return
	}

	foundPartsSupply, err := h.usecase.GetPartsSupplyByID(c.Request.Context(), partsSupplyID)
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toPartsSupplyResponse(foundPartsSupply))
}

// CreatePartsSupply godoc
// @Summary Create a new parts supply
// @Description Create a new parts supply record
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param supply body request.PartsSupplyCreateRequest true "Parts Supply Information"
// @Success 201 {object} response.PartsSupplyResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies [post]
func (h *PartsSupplyHandler) CreatePartsSupply(c *gin.Context) {
	var payload request.PartsSupplyCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidPartsSupplyInput.HTTPStatus, errInvalidPartsSupplyInput.ToHTTPError())
		return
	}

	entity := toPartsSupplyEntityFromCreate(payload)
	createdPartsSupply, err := h.usecase.CreatePartsSupply(c.Request.Context(), &entity)
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, toPartsSupplyResponse(createdPartsSupply))
}

// UpdatePartsSupply godoc
// @Summary Update a parts supply
// @Description Update an existing parts supply record
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Parts Supply ID"
// @Param supply body request.PartsSupplyUpdateRequest true "Parts Supply Information"
// @Success 200 {object} response.OperationMessageResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies/{id} [put]
func (h *PartsSupplyHandler) UpdatePartsSupply(c *gin.Context) {
	partsSupplyID, ok := parsePartsSupplyID(c)
	if !ok {
		return
	}

	var payload request.PartsSupplyUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidPartsSupplyInput.HTTPStatus, errInvalidPartsSupplyInput.ToHTTPError())
		return
	}

	entity := toPartsSupplyEntityFromUpdate(partsSupplyID, payload)
	if err := h.usecase.UpdatePartsSupply(c.Request.Context(), &entity); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Parts supply updated successfully"})
}

// DeletePartsSupply godoc
// @Summary Delete a parts supply
// @Description Delete an existing parts supply record
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Parts Supply ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies/{id} [delete]
func (h *PartsSupplyHandler) DeletePartsSupply(c *gin.Context) {
	partsSupplyID, ok := parsePartsSupplyID(c)
	if !ok {
		return
	}

	if err := h.usecase.DeletePartsSupply(c.Request.Context(), partsSupplyID); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.Status(http.StatusNoContent)
}

// ListPartsSupplies godoc
// @Summary List all parts supplies
// @Description Get a list of all parts supplies
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.PartsSupplyResponse
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies [get]
func (h *PartsSupplyHandler) ListPartsSupplies(c *gin.Context) {
	partsSupplies, err := h.usecase.ListPartsSupplies(c.Request.Context())
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toPartsSupplyResponseList(partsSupplies))
}

func parsePartsSupplyID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(errInvalidPartsSupplyID.HTTPStatus, errInvalidPartsSupplyID.ToHTTPError())
		return 0, false
	}
	return uint(id), true
}

func toPartsSupplyEntityFromCreate(payload request.PartsSupplyCreateRequest) entities.PartsSupply {
	return entities.PartsSupply{
		Name:            payload.Name,
		Description:     payload.Description,
		Price:           payload.Price,
		QuantityTotal:   payload.QuantityTotal,
		QuantityReserve: payload.QuantityReserve,
	}
}

func toPartsSupplyEntityFromUpdate(id uint, payload request.PartsSupplyUpdateRequest) entities.PartsSupply {
	return entities.PartsSupply{
		ID:              id,
		Name:            payload.Name,
		Description:     payload.Description,
		Price:           payload.Price,
		QuantityTotal:   payload.QuantityTotal,
		QuantityReserve: payload.QuantityReserve,
	}
}

func toPartsSupplyResponse(partsSupply entities.PartsSupply) response.PartsSupplyResponse {
	return response.PartsSupplyResponse{
		ID:              partsSupply.ID,
		Name:            partsSupply.Name,
		Description:     partsSupply.Description,
		Price:           partsSupply.Price,
		QuantityTotal:   partsSupply.QuantityTotal,
		QuantityReserve: partsSupply.QuantityReserve,
	}
}

func toPartsSupplyResponseList(partsSupplies []entities.PartsSupply) []response.PartsSupplyResponse {
	result := make([]response.PartsSupplyResponse, 0, len(partsSupplies))
	for _, partsSupply := range partsSupplies {
		result = append(result, toPartsSupplyResponse(partsSupply))
	}
	return result
}

func parseStockOperationPayload(c *gin.Context) ([]entities.PartsSupply, bool) {
	var payload []request.PartsSupplyStockOperationRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidPartsSupplyInput.HTTPStatus, errInvalidPartsSupplyInput.ToHTTPError())
		return nil, false
	}

	result := make([]entities.PartsSupply, 0, len(payload))
	for _, item := range payload {
		id, err := strconv.ParseUint(item.ID, 10, 64)
		if err != nil {
			c.JSON(errInvalidPartsSupplyID.HTTPStatus, errInvalidPartsSupplyID.ToHTTPError())
			return nil, false
		}

		result = append(result, entities.PartsSupply{
			ID:              uint(id),
			QuantityReserve: item.Quantity,
		})
	}

	return result, true
}
