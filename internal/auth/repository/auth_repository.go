// internal/auth/repository/auth_repository.go
package repository

import (
	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *domain.User) error
	CreateOrganization(org *domain.Organization) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	UpdateUser(user *domain.User) error
}

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (r *authRepositoryImpl) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *authRepositoryImpl) CreateOrganization(org *domain.Organization) error {
	return r.db.Create(org).Error
}

func (r *authRepositoryImpl) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}
