// internal/auth/service/organization_service.go
package service

import (
	"errors"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/repository"
	commonerror "github.com/Axontik/comin-authentication-service/internal/common/errors"
	"github.com/google/uuid"
)

type OrganizationService struct {
	repo repository.OrganizationRepository
}

func NewOrganizationService(repo repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) Create(req *domain.OrganizationRequest) (*domain.OrganizationResponse, error) {
	org := &domain.Organization{
		Name:   req.Name,
		Slug:   generateSlug(req.Name),
		Domain: req.Domain,
		Status: "active",
		Settings: domain.Settings{
			Theme:          "default",
			AllowedDomains: []string{},
			MaxUsers:       10, // Default limit
			Features:       []string{"basic"},
		},
	}

	if err := s.repo.Create(org); err != nil {
		return nil, err
	}

	return toOrganizationResponse(org), nil
}

func (s *OrganizationService) GetByID(id uuid.UUID) (*domain.OrganizationResponse, error) {
	org, err := s.repo.GetByID(id)
	if err != nil {
		return nil, commonerror.NewNotFoundError("Organization not found")
	}
	return toOrganizationResponse(org), nil
}

func (s *OrganizationService) Update(id uuid.UUID, req *domain.OrganizationRequest) (*domain.OrganizationResponse, error) {
	org, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	org.Name = req.Name
	org.Domain = req.Domain

	if err := s.repo.Update(org); err != nil {
		return nil, err
	}

	return toOrganizationResponse(org), nil
}

func (s *OrganizationService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func toOrganizationResponse(org *domain.Organization) *domain.OrganizationResponse {
	return &domain.OrganizationResponse{
		ID:     org.ID,
		Name:   org.Name,
		Slug:   org.Slug,
		Domain: org.Domain,
		Status: org.Status,
	}
}
