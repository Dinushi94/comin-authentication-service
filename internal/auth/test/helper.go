// internal/auth/test/helper.go
package test

import (
	"time"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/google/uuid"
)

func MockOrganization() *domain.Organization {
	return &domain.Organization{
		Base: domain.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:   "Test Organization",
		Slug:   "test-organization",
		Domain: "test.com",
		Status: "active",
		Settings: domain.Settings{
			Theme:          "default",
			AllowedDomains: []string{"test.com"},
			MaxUsers:       10,
			Features:       []string{"basic"},
		},
	}
}

func MockUser() *domain.User {
	return &domain.User{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		OrganizationID: uuid.New(),
		Email:          "test@test.com",
		Password:       "hashedpassword",
		FirstName:      "Test",
		LastName:       "User",
		Role:           "admin",
		Status:         "active",
	}
}
