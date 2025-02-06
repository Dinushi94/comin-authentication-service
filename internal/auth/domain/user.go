// internal/auth/domain/user.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid"`
	Email          string    `json:"email" gorm:"unique;not null"`
	Password       string    `json:"-" gorm:"not null"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Role           string    `json:"role" gorm:"not null"`
	Status         string    `json:"status" gorm:"default:'active'"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
