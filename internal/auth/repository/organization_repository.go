// internal/auth/repository/organization_repository.go
package repository

import (
	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(org *domain.Organization) error
	GetByID(id uuid.UUID) (*domain.Organization, error)
	Update(org *domain.Organization) error
	Delete(id uuid.UUID) error
}

type organizationRepositoryImpl struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepositoryImpl{db: db}
}

func (r *organizationRepositoryImpl) Create(org *domain.Organization) error {
	return r.db.Create(org).Error
}

func (r *organizationRepositoryImpl) GetByID(id uuid.UUID) (*domain.Organization, error) {
	var org domain.Organization
	err := r.db.First(&org, "id = ?", id).Error
	return &org, err
}

func (r *organizationRepositoryImpl) Update(org *domain.Organization) error {
	return r.db.Save(org).Error
}

func (r *organizationRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Organization{}, "id = ?", id).Error
}

func (r *organizationRepositoryImpl) GetBySlug(slug string) (*domain.Organization, error) {
	var org domain.Organization
	err := r.db.First(&org, "slug = ?", slug).Error
	return &org, err
}
