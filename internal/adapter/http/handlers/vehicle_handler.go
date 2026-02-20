package handlers

import (
	"errors"
	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"mecanica_xpto/internal/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// Domain errors
	errInvalidVehicleID = pkg.NewDomainErrorSimple("INVALID_VEHICLE_ID", "Invalid vehicle ID", http.StatusBadRequest)
	errInvalidInput     = pkg.NewDomainErrorSimple("INVALID_INPUT", "Invalid input data", http.StatusBadRequest)
	errEmptyPlate       = pkg.NewDomainErrorSimple("EMPTY_PLATE", "Plate cannot be empty", http.StatusBadRequest)
)

// VehicleHandler handles HTTP requests for vehicle operations
// @title Vehicle API
// @version 1.0
// @description API for managing vehicles in the workshop management system
type VehicleHandler struct {
	service usecase.VehicleServiceInterface
}

func NewVehicleHandler(service usecase.VehicleServiceInterface) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

func mapVehicleError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrVehicleNotFound):
		return pkg.NewDomainErrorSimple("VEHICLE_NOT_FOUND", "Vehicle not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidPlateFormat):
		return pkg.NewDomainErrorSimple("INVALID_PLATE_FORMAT", "Invalid plate format", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrVehicleAlreadyExists):
		return pkg.NewDomainErrorSimple("VEHICLE_EXISTS", "Vehicle already exists", http.StatusConflict)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid vehicle ID", http.StatusBadRequest)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetVehicles godoc
// @Summary Get all vehicles
// @Description Retrieves a list of all vehicles
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.VehicleResponse
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles [get]
func (v VehicleHandler) GetVehicles(c *gin.Context) {
	vehicles, err := v.service.GetAllVehicles()
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, toVehicleResponseList(vehicles))
}

// GetVehicleByCustomerID godoc
// @Summary Get vehicles by customer ID
// @Description Retrieves all vehicles belonging to a specific customer
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param customerID path int true "Customer ID"
// @Success 200 {array} response.VehicleResponse
// @Failure 400 {object} pkg.ErrorResponse "Invalid customer ID format"
// @Failure 404 {object} pkg.ErrorResponse "Customer not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/customer/{customerID} [get]
func (v VehicleHandler) GetVehicleByCustomerID(c *gin.Context) {
	customerID, ok := parseVehicleIDParam(c, "customerID")
	if !ok {
		return
	}

	vehicles, err := v.service.GetVehiclesByCustomerID(customerID)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, toVehiclesByCustomerResponse(vehicles))
}

// GetVehicleByID godoc
// @Summary Get vehicle by ID
// @Description Retrieves a vehicle by its ID
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} response.VehicleResponse
// @Failure 400 {object} pkg.ErrorResponse "Invalid vehicle ID format"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/{id} [get]
func (v VehicleHandler) GetVehicleByID(c *gin.Context) {
	id, ok := parseVehicleIDParam(c, "id")
	if !ok {
		return
	}

	vehicle, err := v.service.GetVehicleByID(id)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toVehicleResponseFromPtr(vehicle))
}

// GetVehicleByPlate godoc
// @Summary Get vehicle by plate
// @Description Retrieves a vehicle by its license plate
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param plate path string true "Vehicle license plate"
// @Success 200 {object} response.VehicleResponse
// @Failure 400 {object} pkg.ErrorResponse "Invalid plate format or empty plate"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/plate/{plate} [get]
func (v VehicleHandler) GetVehicleByPlate(c *gin.Context) {
	plate := c.Param("plate")
	if plate == "" {
		c.JSON(errEmptyPlate.HTTPStatus, errEmptyPlate.ToHTTPError())
		return
	}

	vehicle, err := v.service.GetVehicleByPlate(plate)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, toVehicleResponseFromPtr(vehicle))
}

// CreateVehicle godoc
// @Summary Create a new vehicle
// @Description Creates a new vehicle record
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param vehicle body request.VehicleCreateRequest true "Vehicle information"
// @Success 201 {object} response.VehicleResponse "Vehicle created successfully"
// @Failure 400 {object} pkg.ErrorResponse "Invalid input data or plate format"
// @Failure 409 {object} pkg.ErrorResponse "Vehicle already exists"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles [post]
func (v VehicleHandler) CreateVehicle(c *gin.Context) {
	var payload request.VehicleCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
		return
	}

	vehicle := toVehicleEntityFromCreate(payload)
	createdVehicle, err := v.service.CreateVehicle(vehicle)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusCreated, toVehicleResponseFromPtr(createdVehicle))
}

// UpdateVehicle godoc
// @Summary Update a vehicle partially
// @Description Updates specific fields of an existing vehicle
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Param vehicle body request.VehicleUpdateRequest true "Vehicle information to update"
// @Success 200 {object} response.OperationMessageResponse "Vehicle updated successfully"
// @Failure 400 {object} pkg.ErrorResponse "Invalid input data, ID format or plate format"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/{id} [patch]
func (v VehicleHandler) UpdateVehicle(c *gin.Context) {
	id, ok := parseVehicleIDParam(c, "id")
	if !ok {
		return
	}

	var payload request.VehicleUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
		return
	}

	updates := toVehicleUpdateMap(payload)
	if len(updates) == 0 {
		appErr := pkg.NewDomainErrorSimple("EMPTY_UPDATE_PAYLOAD", "At least one field must be provided", http.StatusBadRequest)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	result, err := v.service.UpdateVehiclePartial(id, updates)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: result})
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Deletes a vehicle by its ID
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.ErrorResponse "Invalid vehicle ID format"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/{id} [delete]
func (v VehicleHandler) DeleteVehicle(c *gin.Context) {
	id, ok := parseVehicleIDParam(c, "id")
	if !ok {
		return
	}

	if err := v.service.DeleteVehicle(id); err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.Status(http.StatusNoContent)
}

func parseVehicleIDParam(c *gin.Context, param string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		c.JSON(errInvalidVehicleID.HTTPStatus, errInvalidVehicleID.ToHTTPError())
		return 0, false
	}
	return uint(id), true
}

func toVehicleEntityFromCreate(payload request.VehicleCreateRequest) entities.Vehicle {
	return entities.Vehicle{
		CustomerID: payload.CustomerID,
		Brand:      payload.Brand,
		Model:      payload.Model,
		Year:       payload.Year,
		Plate:      valueobject.ParsePlate(payload.Plate),
	}
}

func toVehicleUpdateMap(payload request.VehicleUpdateRequest) map[string]interface{} {
	updates := make(map[string]interface{})
	if payload.CustomerID != nil {
		updates["customer_id"] = float64(*payload.CustomerID)
	}
	if payload.Brand != nil {
		updates["brand"] = *payload.Brand
	}
	if payload.Model != nil {
		updates["model"] = *payload.Model
	}
	if payload.Year != nil {
		updates["year"] = *payload.Year
	}
	if payload.Plate != nil {
		updates["plate"] = *payload.Plate
	}
	return updates
}

func toVehicleResponse(vehicle entities.Vehicle) response.VehicleResponse {
	return response.VehicleResponse{
		ID:         vehicle.ID,
		CustomerID: vehicle.CustomerID,
		Brand:      vehicle.Brand,
		Model:      vehicle.Model,
		Year:       vehicle.Year,
		Plate:      vehicle.Plate.String(),
	}
}

func toVehicleResponseFromPtr(vehicle *entities.Vehicle) response.VehicleResponse {
	if vehicle == nil {
		return response.VehicleResponse{}
	}
	return toVehicleResponse(*vehicle)
}

func toVehicleResponseList(vehicles []entities.Vehicle) []response.VehicleResponse {
	result := make([]response.VehicleResponse, 0, len(vehicles))
	for _, vehicle := range vehicles {
		result = append(result, toVehicleResponse(vehicle))
	}
	return result
}

func toVehiclesByCustomerResponse(vehicles []entities.Vehicle) response.VehiclesByCustomerResponse {
	response := response.VehiclesByCustomerResponse{
		Vehicles: toVehicleResponseList(vehicles),
	}

	for _, vehicle := range vehicles {
		if vehicle.Customer != nil {
			response.Customer = toVehicleCustomerSummary(vehicle.Customer)
			break
		}
	}

	return response
}

func toVehicleCustomerSummary(customer *entities.Customer) *response.VehicleCustomerSummary {
	if customer == nil {
		return nil
	}

	return &response.VehicleCustomerSummary{
		ID:          customer.ID,
		FullName:    customer.FullName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Document:    customer.CpfCnpj.String(),
	}
}
