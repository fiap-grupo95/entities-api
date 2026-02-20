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
	errInvalidServiceID    = pkg.NewDomainErrorSimple("INVALID_SERVICE_ID", "Invalid service ID", http.StatusBadRequest)
	errInvalidServiceInput = pkg.NewDomainErrorSimple("INVALID_SERVICE_INPUT", "Invalid service payload", http.StatusBadRequest)
)

// ServiceHandler handles HTTP requests for service operations
// @title Service API
// @version 1.0
// @description API for managing services in the workshop management system
type ServiceHandler struct {
	usecase usecase.IServiceUseCase
}

func NewServiceHandler(usecase usecase.IServiceUseCase) *ServiceHandler {
	return &ServiceHandler{usecase: usecase}
}

func mapServiceError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrServiceNotFound):
		return pkg.NewDomainErrorSimple("SERVICE_NOT_FOUND", "Service not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid service ID", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrServiceAlreadyExists):
		return pkg.NewDomainErrorSimple("SERVICE_EXISTS", "Service already exists", http.StatusConflict)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetServiceByID godoc
// @Summary Get service by ID
// @Description Retrieve a service by its ID
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} response.ServiceResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services/{id} [get]
func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	serviceID, ok := parseServiceID(c)
	if !ok {
		return
	}

	service, err := h.usecase.GetServiceByID(c.Request.Context(), serviceID)
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, toServiceResponse(service))
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service record
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param service body request.ServiceCreateRequest true "Service Information"
// @Success 201 {object} response.ServiceResponse
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var payload request.ServiceCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidServiceInput.HTTPStatus, errInvalidServiceInput.ToHTTPError())
		return
	}
	entity := toServiceEntityFromCreate(payload)
	createdService, err := h.usecase.CreateService(c.Request.Context(), &entity)
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusCreated, toServiceResponse(createdService))
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service record
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param service body request.ServiceUpdateRequest true "Service Information"
// @Success 200 {object} response.OperationMessageResponse
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	serviceID, ok := parseServiceID(c)
	if !ok {
		return
	}
	var payload request.ServiceUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidServiceInput.HTTPStatus, errInvalidServiceInput.ToHTTPError())
		return
	}
	entity := toServiceEntityFromUpdate(serviceID, payload)
	if err := h.usecase.UpdateService(c.Request.Context(), &entity); err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Service updated successfully"})
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete an existing service record
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	serviceID, ok := parseServiceID(c)
	if !ok {
		return
	}
	if err := h.usecase.DeleteService(c.Request.Context(), serviceID); err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.Status(http.StatusNoContent)
}

// ListServices godoc
// @Summary List all services
// @Description Get a list of all services
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.ServiceResponse
// @Failure 500 {object} pkg.AppError
// @Router /services [get]
func (h *ServiceHandler) ListServices(c *gin.Context) {
	services, err := h.usecase.ListServices(c.Request.Context())
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, toServiceResponseList(services))
}

func parseServiceID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(errInvalidServiceID.HTTPStatus, errInvalidServiceID.ToHTTPError())
		return 0, false
	}
	return uint(id), true
}

func toServiceEntityFromCreate(payload request.ServiceCreateRequest) entities.Service {
	return entities.Service{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
	}
}

func toServiceEntityFromUpdate(id uint, payload request.ServiceUpdateRequest) entities.Service {
	return entities.Service{
		ID:          id,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
	}
}

func toServiceResponse(service entities.Service) response.ServiceResponse {
	return response.ServiceResponse{
		ID:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		Price:       service.Price,
	}
}

func toServiceResponseList(services []entities.Service) []response.ServiceResponse {
	result := make([]response.ServiceResponse, 0, len(services))
	for _, service := range services {
		result = append(result, toServiceResponse(service))
	}
	return result
}
