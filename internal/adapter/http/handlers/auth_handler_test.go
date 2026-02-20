package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mecanica_xpto/internal/adapter/http/handlers"
	"mecanica_xpto/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock do usecase.AuthInterface
type mockAuthUsecase struct {
	loginFunc func(Email, Password string) (string, *pkg.AppError)
}

func (m *mockAuthUsecase) Login(Email, Password string) (string, *pkg.AppError) {
	return m.loginFunc(Email, Password)
}

func TestAuthHandlerLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           string
		mockLoginFunc  func(Email, Password string) (string, *pkg.AppError)
		wantStatusCode int
		wantBodySubstr string
	}{
		{
			name: "Login válido",
			body: `{"email":"user@example.com","password":"senha123"}`,
			mockLoginFunc: func(Email, Password string) (string, *pkg.AppError) {
				return "token_jwt_mock", nil
			},
			wantStatusCode: http.StatusOK,
			wantBodySubstr: `"token":"token_jwt_mock"`,
		},
		{
			name: "Corpo inválido JSON mal formado",
			body: `{"email":user@example.com,"password":"senha123"}`, // email sem aspas, inválido JSON
			mockLoginFunc: func(Email, Password string) (string, *pkg.AppError) {
				return "", nil // não será chamado
			},
			wantStatusCode: http.StatusBadRequest,
			wantBodySubstr: `"code":"INVALID_REQUEST"`,
		},
		{
			name: "Login falha - erro de autenticação",
			body: `{"email":"user@example.com","password":"senha123"}`,
			mockLoginFunc: func(Email, Password string) (string, *pkg.AppError) {
				return "", pkg.NewDomainErrorSimple("UNAUTHORIZED", "Credenciais inválidas", http.StatusUnauthorized)
			},
			wantStatusCode: http.StatusUnauthorized,
			wantBodySubstr: `"code":"UNAUTHORIZED"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup gin engine e request
			r := gin.New()
			mockUC := &mockAuthUsecase{loginFunc: tt.mockLoginFunc}
			handler := handlers.NewAuthHandler(mockUC)

			r.POST("/v1/login", handler.Login)

			req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBodySubstr)
		})
	}
}
