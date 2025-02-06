// internal/auth/domain/organization.go
package domain

import "github.com/google/uuid"

type Organization struct {
    Base
    Name     string   `json:"name" gorm:"not null"`
    Slug     string   `json:"slug" gorm:"unique;not null"`
    Domain   string   `json:"domain"`
    Settings Settings `json:"settings" gorm:"type:jsonb;default:'{}'"`
    Status   string   `json:"status" gorm:"default:'active'"`
}

type Settings struct {
    Theme          string   `json:"theme"`
    AllowedDomains []string `json:"allowed_domains"`
    MaxUsers       int      `json:"max_users"`
    Features       []string `json:"features"`
}

type OrganizationRequest struct {
    Name   string `json:"name" binding:"required"`
    Domain string `json:"domain"`
}

type OrganizationResponse struct {
    ID     uuid.UUID `json:"id"`
    Name   string    `json:"name"`
    Slug   string    `json:"slug"`
    Domain string    `json:"domain"`
    Status string    `json:"status"`
}