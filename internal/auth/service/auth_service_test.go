// internal/auth/service/auth_service_test.go
package service

import (
	"testing"
	"time"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/repository"
	"github.com/Axontik/comin-authentication-service/internal/auth/test"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var _ repository.AuthRepository = (*MockAuthRepo)(nil)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepo) CreateOrganization(organization *domain.Organization) error {
	args := m.Called(organization)
	return args.Error(0)
}

func (m *MockAuthRepo) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthRepo) GetUserByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthRepo) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	mockJWT := jwt.NewJWTService("test-secret", "test-issuer")
	service := NewAuthService(mockRepo, mockJWT)

	req := &domain.RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	mockRepo.On("GetUserByEmail", req.Email).Return(nil, nil)
	mockRepo.On("CreateOrganization", mock.AnythingOfType("*domain.Organization")).Return(nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil)

	response, err := service.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	mockJWT := jwt.NewJWTService("test-secret", "test-issuer")
	service := NewAuthService(mockRepo, mockJWT)

	// Create mock user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := test.MockUser()
	mockUser.Password = string(hashedPassword)

	req := &domain.LoginRequest{
		Email:    mockUser.Email,
		Password: "password123",
	}

	// Setup mock expectations
	mockRepo.On("GetUserByEmail", req.Email).Return(mockUser, nil)

	// Perform login
	response, err := service.Login(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	mockJWT := jwt.NewJWTService("test-secret", "test-issuer")
	service := NewAuthService(mockRepo, mockJWT)

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Setup mock expectations for invalid credentials
	mockRepo.On("GetUserByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound)

	// Perform login
	response, err := service.Login(req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	mockJWT := jwt.NewJWTService("test-secret", "test-issuer")
	service := NewAuthService(mockRepo, mockJWT)

	mockUser := test.MockUser()
	token, _ := mockJWT.GenerateToken(mockUser.ID, mockUser.OrganizationID, mockUser.Email, mockUser.Role, time.Hour)

	mockRepo.On("GetUserByID", mockUser.ID).Return(mockUser, nil)

	response, err := service.RefreshToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := new(MockAuthRepo)
	mockJWT := jwt.NewJWTService("test-secret", "test-issuer")
	service := NewAuthService(mockRepo, mockJWT)

	mockUser := test.MockUser()
	token, _ := mockJWT.GenerateToken(mockUser.ID, mockUser.OrganizationID, mockUser.Email, mockUser.Role, time.Hour)

	mockRepo.On("GetUserByID", mockUser.ID).Return(mockUser, nil)

	claims, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	mockRepo.AssertExpectations(t)
}
