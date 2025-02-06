// internal/auth/service/organization_service_test.go
package service

import (
	"testing"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/repository"
	"github.com/Axontik/comin-authentication-service/internal/auth/test"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ repository.OrganizationRepository = (*MockOrganizationRepo)(nil)

type MockOrganizationRepo struct {
	mock.Mock
}

func (m *MockOrganizationRepo) Create(org *domain.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepo) GetByID(id uuid.UUID) (*domain.Organization, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrganizationRepo) Update(org *domain.Organization) error {
	args := m.Called(org)
	return args.Error(0)
}

func (m *MockOrganizationRepo) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestOrganizationService_Create(t *testing.T) {
	mockRepo := new(MockOrganizationRepo)
	service := NewOrganizationService(mockRepo)

	req := &domain.OrganizationRequest{
		Name:   "Test Org",
		Domain: "test.com",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Organization")).Return(nil)

	response, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationService_GetByID(t *testing.T) {
	mockRepo := new(MockOrganizationRepo)
	service := NewOrganizationService(mockRepo)

	mockOrg := test.MockOrganization()
	mockRepo.On("GetByID", mockOrg.ID).Return(mockOrg, nil)

	response, err := service.GetByID(mockOrg.ID)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, mockOrg.ID, response.ID)
	assert.Equal(t, mockOrg.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestOrganizationService_Update(t *testing.T) {
	mockRepo := new(MockOrganizationRepo)
	service := NewOrganizationService(mockRepo)

	mockOrg := test.MockOrganization()
	req := &domain.OrganizationRequest{
		Name:   "Updated Org",
		Domain: "updated.com",
	}

	mockRepo.On("GetByID", mockOrg.ID).Return(mockOrg, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Organization")).Return(nil)

	response, err := service.Update(mockOrg.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.Domain, response.Domain)
	mockRepo.AssertExpectations(t)
}
