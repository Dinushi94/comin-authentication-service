// test/integration/auth_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/handler"
	"github.com/Axontik/comin-authentication-service/internal/auth/repository"
	"github.com/Axontik/comin-authentication-service/internal/auth/service"
	"github.com/Axontik/comin-authentication-service/internal/auth/test"
	"github.com/Axontik/comin-authentication-service/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	db     *gorm.DB
	router *gin.Engine
}

func setupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	db := test.SetupTestDB(t)

	jwtService := jwt.NewJWTService("test-secret", "test-issuer")
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
	}

	return &IntegrationTestSuite{
		db:     db,
		router: router,
	}
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupTestDB(db)

	db = db.Debug()

	gin.SetMode(gin.TestMode)

	jwtService := jwt.NewJWTService("test-secret", "test-issuer")
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.New()
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
	}

	org := &domain.Organization{
		Name:   "Test Direct Org",
		Slug:   "test-direct-org",
		Domain: "test.com",
		Status: "active",
		Settings: domain.Settings{
			Theme:          "default",
			AllowedDomains: []string{"test.com"},
			MaxUsers:       10,
			Features:       []string{"basic"},
		},
	}

	err := db.Create(org).Error
	assert.NoError(t, err, "Direct organization creation should work")

	// Test Register
	registerReq := domain.RegisterRequest{
		OrganizationName: "Test Org",
		Email:            "test@example.com",
		Password:         "password123",
		FirstName:        "Test",
		LastName:         "User",
	}

	registerBody, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	t.Logf("Register Response Code: %d", w.Code)
	t.Logf("Register Response Body: %s", w.Body.String())
	// Log the error response if status is not 201
	if w.Code != http.StatusCreated {
		// Try to parse error response
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &errResp); err == nil {
			t.Logf("Error message: %s", errResp.Error)
		}
	}

	assert.Equal(t, http.StatusCreated, w.Code)

	var registerResp domain.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &registerResp)
	if err != nil {
		t.Logf("Register response unmarshal error: %v", err)
	}
	assert.NoError(t, err)
	assert.NotEmpty(t, registerResp.AccessToken)

	// Test Login with proper error logging
	loginReq := domain.LoginRequest{
		Email:    registerReq.Email,
		Password: registerReq.Password,
	}

	loginBody, _ := json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Login Response Body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp domain.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp.AccessToken)
}

func TestIntegration_RefreshToken(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupTestDB(db)

	gin.SetMode(gin.TestMode)

	jwtService := jwt.NewJWTService("test-secret", "test-issuer")
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.New()
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}
	}

	// First register and login to get valid tokens
	registerReq := domain.RegisterRequest{
		OrganizationName: "Test Org",
		Email:            "test@example.com",
		Password:         "password123",
		FirstName:        "Test",
		LastName:         "User",
	}

	registerBody, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Logf("Register Response Body: %s", w.Body.String())
		t.FailNow()
	}

	var registerResp domain.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NoError(t, err)

	// Test refresh token
	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: registerResp.RefreshToken,
	}

	refreshBody, _ := json.Marshal(refreshReq)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Refresh Response Body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var refreshResp domain.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &refreshResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshResp.AccessToken)
}
