package repository_test

import (
	"testing"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&repository.UserModel{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
	return db
}

func TestUserRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	// --- Create ---
	user := &domain.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secure123",
	}

	err := repo.Save(user)
	assert.NoError(t, err, "save user")

	// --- Read (FindByEmail) ---
	found, err := repo.FindByEmail("john@example.com")
	assert.NoError(t, err, "find user by email")
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Name, found.Name)

	// --- Read (FindByID) ---
	byID, err := repo.FindByID(found.ID)
	assert.NoError(t, err, "find user by ID")
	assert.Equal(t, found.ID, byID.ID)

	// --- Update ---
	byID.Name = "Updated Name"
	byID.Password = "newpass"
	err = repo.Update(byID)
	assert.NoError(t, err, "update user")

	updated, err := repo.FindByID(byID.ID)
	assert.NoError(t, err, "read updated user")
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "newpass", updated.Password)
	assert.True(t, updated.UpdatedAt.After(updated.CreatedAt), "UpdatedAt should be after CreatedAt")

	// --- Delete ---
	err = repo.Delete(byID.ID)
	assert.NoError(t, err, "delete user")

	deleted, err := repo.FindByID(byID.ID)
	assert.Error(t, err, "should not find deleted user")
	assert.Nil(t, deleted)
}
