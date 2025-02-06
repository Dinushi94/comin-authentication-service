// internal/auth/repository/auth_repository_test.go
package repository

import (
	"testing"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"github.com/Axontik/comin-authentication-service/internal/auth/test"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database connection
	dsn := "host=localhost user=comin_owner password=password dbname=comin_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate tables
	err = db.AutoMigrate(&domain.Organization{}, &domain.User{})
	if err != nil {
		t.Fatalf("Failed to migrate tables: %v", err)
	}

	return db
}

func cleanupTestDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE users CASCADE")
	db.Exec("TRUNCATE TABLE organizations CASCADE")
}

func TestAuthRepository_CreateUser(t *testing.T) {
	db := test.SetupTestDB(t)
	defer test.CleanupTestDB(db)

	repo := NewAuthRepository(db)
	mockUser := test.MockUser()

	err := repo.CreateUser(mockUser)
	assert.NoError(t, err)

	var savedUser domain.User
	err = db.First(&savedUser, "email = ?", mockUser.Email).Error
	assert.NoError(t, err)
	assert.Equal(t, mockUser.Email, savedUser.Email)
}

func TestAuthRepository_GetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	repo := NewAuthRepository(db)
	mockUser := test.MockUser()

	err := repo.CreateUser(mockUser)
	assert.NoError(t, err)

	// Test getting user
	foundUser, err := repo.GetUserByEmail(mockUser.Email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, mockUser.Email, foundUser.Email)
}

func TestAuthRepository_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	repo := NewAuthRepository(db)
	mockUser := test.MockUser()

	// Create user first
	err := repo.CreateUser(mockUser)
	assert.NoError(t, err)

	// Test getting user by ID
	foundUser, err := repo.GetUserByID(mockUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, mockUser.ID, foundUser.ID)
}

func TestAuthRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(db)

	repo := NewAuthRepository(db)
	mockUser := test.MockUser()

	// Create user first
	err := repo.CreateUser(mockUser)
	assert.NoError(t, err)

	// Update user
	mockUser.FirstName = "Updated"
	err = repo.UpdateUser(mockUser)
	assert.NoError(t, err)

	// Verify update
	foundUser, err := repo.GetUserByID(mockUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", foundUser.FirstName)
}
