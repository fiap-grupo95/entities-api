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
	errInvalidUserID    = pkg.NewDomainErrorSimple("INVALID_USER_ID", "Invalid user ID", http.StatusBadRequest)
	errInvalidUserInput = pkg.NewDomainErrorSimple("INVALID_USER_INPUT", "Invalid user payload", http.StatusBadRequest)
)

type UserHandler struct {
	usecase usecase.IUserUseCase
}

func NewUserHandler(usecase usecase.IUserUseCase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func mapUserError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return pkg.NewDomainErrorSimple("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrUserAlreadyExists):
		return pkg.NewDomainErrorSimple("USER_ALREADY_EXISTS", "User already exists", http.StatusConflict)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	user, err := h.usecase.GetByID(userID)
	if err != nil {
		appErr := mapUserError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var payload request.UserCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidUserInput.HTTPStatus, errInvalidUserInput.ToHTTPError())
		return
	}

	createdUser, err := h.usecase.Create(&entities.User{
		Email:    payload.Email,
		Password: payload.Password,
		UserType: valueobject.ParseUserType(payload.UserType),
	})
	if err != nil {
		appErr := mapUserError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(createdUser))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	var payload request.UserUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(errInvalidUserInput.HTTPStatus, errInvalidUserInput.ToHTTPError())
		return
	}

	err := h.usecase.Update(userID, &entities.User{
		Email:    payload.Email,
		Password: payload.Password,
		UserType: valueobject.ParseUserType(payload.UserType),
	})
	if err != nil {
		appErr := mapUserError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, response.OperationMessageResponse{Message: "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		return
	}

	if err := h.usecase.Delete(userID); err != nil {
		appErr := mapUserError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.usecase.List()
	if err != nil {
		appErr := mapUserError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	result := make([]response.UserResponse, 0, len(users))
	for _, user := range users {
		u := user
		result = append(result, toUserResponse(&u))
	}

	c.JSON(http.StatusOK, result)
}

func parseUserID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(errInvalidUserID.HTTPStatus, errInvalidUserID.ToHTTPError())
		return 0, false
	}

	return uint(id), true
}

func toUserResponse(user *entities.User) response.UserResponse {
	if user == nil {
		return response.UserResponse{}
	}

	return response.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		UserType: user.UserType.String(),
	}
}
