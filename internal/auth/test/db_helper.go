// internal/auth/test/db_helper.go
package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Axontik/comin-authentication-service/internal/auth/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	TestDBHost     = "localhost"
	TestDBUser     = "comin_owner"
	TestDBPassword = "password"
	TestDBName     = "comin_test"
	TestDBPort     = "5432"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		TestDBHost, TestDBUser, TestDBPassword, TestDBName, TestDBPort,
	)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate test tables
	err = db.AutoMigrate(&domain.Organization{}, &domain.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func CleanupTestDB(db *gorm.DB) {
	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	// Clean up tables
	db.Exec("TRUNCATE TABLE users CASCADE")
	db.Exec("TRUNCATE TABLE organizations CASCADE")

	// Close connection
	sqlDB.Close()
}

func WaitForDB(t *testing.T, maxRetries int) {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		TestDBHost, TestDBUser, TestDBPassword, TestDBName, TestDBPort,
	)

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					return
				}
			}
		}
		t.Logf("Waiting for database... attempt %d/%d", i+1, maxRetries)
		time.Sleep(time.Second)
	}
	t.Fatal("Database not available")
}
