// internal/auth/handler/organization_handler.go
package handler

import (
	"net/http"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	var req domain.OrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.orgService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

func (h *OrganizationHandler) GetProfile(c *gin.Context) {
	orgID, _ := c.Get("organizationID")
	id := orgID.(uuid.UUID)

	org, err := h.orgService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) UpdateProfile(c *gin.Context) {
	orgID, _ := c.Get("organizationID")
	id := orgID.(uuid.UUID)

	var req domain.OrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.orgService.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}
