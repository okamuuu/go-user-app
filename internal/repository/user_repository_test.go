package repository_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// マイグレーション
	err = testDB.AutoMigrate(&repository.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	os.Exit(m.Run())
}

func createTestUser() *domain.User {
	return &domain.User{
		Name:      "Test User",
		Email:     uuid.New().String() + "@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestUpdateUser(t *testing.T) {
	repo := repository.NewUserRepository(testDB)
	user := createTestUser()

	// Save
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	// Update
	user.Name = "Updated Name"
	user.Email = "updated@example.com"
	err = repo.Update(user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	// Verify update
	updated, err := repo.FindByEmail("updated@example.com")
	if err != nil {
		t.Fatalf("failed to find updated user: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected name to be 'Updated Name', got '%s'", updated.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	repo := repository.NewUserRepository(testDB)
	user := createTestUser()

	// Save
	err := repo.Save(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	// Delete
	err = repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify delete
	_, err = repo.FindByEmail(user.Email)
	if err == nil {
		t.Fatal("expected error after deleting user, got none")
	}
}
