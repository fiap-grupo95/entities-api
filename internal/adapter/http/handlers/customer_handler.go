package handlers

import (
	"errors"
	request "mecanica_xpto/internal/adapter/http/dto/request"
	response "mecanica_xpto/internal/adapter/http/dto/response"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	use_cases "mecanica_xpto/internal/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CustomerHandler handles HTTP requests for customer operations
// @title Customer API
// @version 1.0
// @description API for managing customers in the workshop management system
type CustomerHandler struct {
	ucCustomer use_cases.ICustomerUseCase
}

var (
	errInvalidCustomerID    = pkg.NewDomainErrorSimple("INVALID_CUSTOMER_ID", "Invalid customer ID", http.StatusBadRequest)
	errInvalidCustomerInput = pkg.NewDomainErrorSimple("INVALID_CUSTOMER_INPUT", "Invalid customer payload", http.StatusBadRequest)
)

// NewCustomerHandler creates a new customer http handler
func NewCustomerHandler(us use_cases.ICustomerUseCase) *CustomerHandler {
	return &CustomerHandler{ucCustomer: us}
}

func mapCustomerError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, use_cases.ErrCustomerNotFound):
		return pkg.NewDomainErrorSimple("CUSTOMER_NOT_FOUND", "Customer not found", http.StatusNotFound)
	case errors.Is(err, use_cases.ErrInvalidDocumentFormat):
		return pkg.NewDomainErrorSimple("INVALID_DOCUMENT", "Invalid document format", http.StatusBadRequest)
	case errors.Is(err, use_cases.ErrCustomerAlreadyExists):
		return pkg.NewDomainErrorSimple("CUSTOMER_EXISTS", "Customer already exists", http.StatusConflict)
	case errors.Is(err, use_cases.ErrInvalidCustomerID):
		return pkg.NewDomainErrorSimple("INVALID_CUSTOMER_ID", "Invalid customer ID", http.StatusBadRequest)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetCustomer godoc
// @Summary Get customer by ID
// @Description Retrieve a customer by their ID
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} response.CustomerResponse
// @Failure 404 {object} map[string]string "error":"customer not found"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{document} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	doc := c.Param("document")

	foundCustomer, err := h.ucCustomer.GetByDocument(doc)
	if err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toCustomerResponse(foundCustomer))
}

// GetCustomerByID godoc
// @Summary Get customer by ID
// @Description Retrieve a customer by numeric ID
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} response.CustomerResponse
// @Failure 400 {object} map[string]string "error":"invalid customer id"
// @Failure 404 {object} map[string]string "error":"customer not found"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	foundCustomer, err := h.ucCustomer.GetById(customerID)
	if err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toCustomerResponse(foundCustomer))
}

// GetFullCustomer godoc
// @Summary Get full customer by ID
// @Description Retrieve a full customer record by their numeric ID
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} response.CustomerResponse
// @Failure 400 {object} map[string]string "error":"invalid customer id"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/id/{id} [get]
func (h *CustomerHandler) GetFullCustomer(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	foundCustomer, err := h.ucCustomer.GetById(customerID)
	if err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toCustomerResponse(foundCustomer))
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Creates a new customer record
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param customer body request.CustomerCreateRequest true "Customer information"
// @Success 201 {object} response.CustomerResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "error message"
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var payload request.CustomerCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidCustomerInput.HTTPStatus, errInvalidCustomerInput.ToHTTPError())
		return
	}

	customer, err := toCustomerEntityFromCreate(payload)
	if err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	createdCustomer, err := h.ucCustomer.CreateCustomer(&customer)
	if err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, toCustomerResponse(createdCustomer))
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update an existing customer record by ID
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body request.CustomerUpdateRequest true "Customer information"
// @Success 200 {object} response.OperationMessageResponse
// @Failure 400 {object} map[string]string "error":"invalid customer id or input"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	var payload request.CustomerUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidCustomerInput.HTTPStatus, errInvalidCustomerInput.ToHTTPError())
		return
	}

	customer := toCustomerEntityFromUpdate(payload)
	if err := h.ucCustomer.UpdateCustomer(customerID, &customer); err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "Customer updated successfully"})
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer record by ID
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string "error":"invalid customer id"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	if err := h.ucCustomer.DeleteCustomer(customerID); err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListCustomer godoc
// @Summary List all customers
// @Description Retrieve a list of all customers
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} response.CustomerResponse
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers [get]
func (h *CustomerHandler) ListCustomer(c *gin.Context) {
	customers, err := h.ucCustomer.ListCustomer()
	if err != nil {
		appErr := mapCustomerError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toCustomerResponseList(customers))
}

func parseCustomerID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(errInvalidCustomerID.HTTPStatus, errInvalidCustomerID.ToHTTPError())
		return 0, false
	}
	return uint(id), true
}

func toCustomerEntityFromCreate(payload request.CustomerCreateRequest) (entities.Customer, error) {
	document, err := valueobject.NewCpfCnpj(payload.Document)
	if err != nil {
		return entities.Customer{}, use_cases.ErrInvalidDocumentFormat
	}

	return entities.Customer{
		FullName:    payload.FullName,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		CpfCnpj:     document,
	}, nil
}

func toCustomerEntityFromUpdate(payload request.CustomerUpdateRequest) entities.Customer {
	return entities.Customer{
		FullName:    payload.FullName,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}
}

func toCustomerResponse(customer *entities.Customer) response.CustomerResponse {
	if customer == nil {
		return response.CustomerResponse{}
	}

	response := response.CustomerResponse{
		ID:          customer.ID,
		UserID:      customer.UserID,
		FullName:    customer.FullName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Document:    customer.CpfCnpj.String(),
	}

	if len(customer.Vehicles) > 0 {
		response.Vehicles = mapCustomerVehicles(customer.Vehicles)
	}

	return response
}

func toCustomerResponseList(customers []entities.Customer) []response.CustomerResponse {
	result := make([]response.CustomerResponse, 0, len(customers))
	for _, customer := range customers {
		c := customer
		result = append(result, toCustomerResponse(&c))
	}
	return result
}

func mapCustomerVehicles(vehicles []entities.Vehicle) []response.CustomerVehicleResponse {
	result := make([]response.CustomerVehicleResponse, 0, len(vehicles))
	for _, vehicle := range vehicles {
		result = append(result, response.CustomerVehicleResponse{
			ID:    vehicle.ID,
			Brand: vehicle.Brand,
			Model: vehicle.Model,
			Year:  vehicle.Year,
			Plate: vehicle.Plate.String(),
		})
	}
	return result
}
