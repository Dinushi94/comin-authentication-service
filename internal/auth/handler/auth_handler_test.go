// internal/auth/handler/auth_handler_test.go
package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (*domain.AuthResponse, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*jwt.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	reqBody := domain.RegisterRequest{
		Email:            "test@example.com",
		Password:         "password123",
		FirstName:        "Test",
		LastName:         "User",
		OrganizationName: "Test Org",
	}

	mockResponse := &domain.AuthResponse{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	mockService.On("Register", &reqBody).Return(mockResponse, nil)

	jsonData, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockResponse.AccessToken, response.AccessToken)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockResponse := &domain.AuthResponse{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	mockService.On("Login", &req).Return(mockResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBytes, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockResponse.AccessToken, response.AccessToken)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	req := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: "old-refresh-token",
	}

	mockResponse := &domain.AuthResponse{
		AccessToken: "new-test-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}

	mockService.On("RefreshToken", req.RefreshToken).Return(mockResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBytes, _ := json.Marshal(req)
	c.Request = httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockResponse.AccessToken, response.AccessToken)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_ValidateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/auth/validate", nil)
	c.Request.Header.Set("Authorization", "Bearer test-token")

	handler.ValidateToken(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
