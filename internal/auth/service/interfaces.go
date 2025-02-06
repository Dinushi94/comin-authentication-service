package service

import (
	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	Register(req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req *domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(refreshToken string) (*domain.AuthResponse, error)
	ValidateToken(token string) (*jwt.Claims, error)
}

type AuthRepositoryInterface interface {
	CreateUser(user *domain.User) error
	CreateOrganization(org *domain.Organization) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	UpdateUser(user *domain.User) error
}
