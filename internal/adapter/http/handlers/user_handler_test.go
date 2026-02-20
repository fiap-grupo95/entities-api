package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/valueobject"
	"mecanica_xpto/internal/usecase"
	"net/http"
	"net/http/httptest"
	"testing"

	request "mecanica_xpto/internal/adapter/http/dto/request"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeUserUseCase struct {
	getFn    func(id uint) (*entities.User, error)
	createFn func(user *entities.User) (*entities.User, error)
	updateFn func(id uint, user *entities.User) error
	deleteFn func(id uint) error
	listFn   func() ([]entities.User, error)
}

func (f fakeUserUseCase) GetByID(id uint) (*entities.User, error)            { return f.getFn(id) }
func (f fakeUserUseCase) Create(user *entities.User) (*entities.User, error) { return f.createFn(user) }
func (f fakeUserUseCase) Update(id uint, user *entities.User) error          { return f.updateFn(id, user) }
func (f fakeUserUseCase) Delete(id uint) error                               { return f.deleteFn(id) }
func (f fakeUserUseCase) List() ([]entities.User, error)                     { return f.listFn() }

func TestUserHandler_CRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewUserHandler(fakeUserUseCase{
		getFn: func(id uint) (*entities.User, error) {
			if id == 1 {
				return &entities.User{ID: 1, Email: "user@xpto.com", UserType: valueobject.ParseUserType("admin")}, nil
			}
			return nil, usecase.ErrUserNotFound
		},
		createFn: func(user *entities.User) (*entities.User, error) {
			if user.Email == "exists@xpto.com" {
				return nil, usecase.ErrUserAlreadyExists
			}
			user.ID = 2
			return user, nil
		},
		updateFn: func(id uint, user *entities.User) error {
			if id == 99 {
				return usecase.ErrUserNotFound
			}
			return nil
		},
		deleteFn: func(id uint) error {
			if id == 99 {
				return usecase.ErrUserNotFound
			}
			return nil
		},
		listFn: func() ([]entities.User, error) {
			return []entities.User{{ID: 1, Email: "user@xpto.com", UserType: valueobject.ParseUserType("admin")}}, nil
		},
	})

	r := gin.New()
	r.GET("/users/:id", handler.GetUser)
	r.POST("/users", handler.CreateUser)
	r.PATCH("/users/:id", handler.UpdateUser)
	r.DELETE("/users/:id", handler.DeleteUser)
	r.GET("/users", handler.ListUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users/abc", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	body, _ := json.Marshal(request.UserCreateRequest{Email: "new@xpto.com", Password: "123456", UserType: "admin"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	body, _ = json.Marshal(request.UserCreateRequest{Email: "exists@xpto.com", Password: "123456", UserType: "admin"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	body, _ = json.Marshal(request.UserUpdateRequest{Email: "update@xpto.com"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPatch, "/users/99", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/users/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/users/99", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/users", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMapUserError_Default(t *testing.T) {
	err := mapUserError(errors.New("boom"))
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
}
