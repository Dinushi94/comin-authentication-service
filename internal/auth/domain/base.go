// internal/auth/domain/base.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;ON_UPDATE:CURRENT_TIMESTAMP"`
}
