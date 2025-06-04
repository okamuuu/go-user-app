package service_test

import (
	"testing"
	"time"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&repository.User{})
	return db
}

func TestAuthService_Login_Success(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	// 事前にユーザー作成しておく
	userRepo.Create(&domain.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: service.HashPassword("secret123"),
	})

	authService := service.NewAuthService(userRepo, []byte("testsecret"), time.Hour)

	// 実行
	token, err := authService.Login("test@example.com", "secret123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	db := setupTestDB()
	userRepo := repository.NewUserRepository(db)

	userRepo.Create(&domain.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: service.HashPassword("secret123"),
	})

	authService := service.NewAuthService(userRepo, []byte("testsecret"), time.Hour)

	token, err := authService.Login("test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
}
