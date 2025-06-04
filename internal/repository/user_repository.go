package repository

import (
	"time"

	"github.com/okamuuu/go-user-app/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save inserts a new user into the database
func (r *UserRepository) Save(user *domain.User) error {
	model := User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	return r.db.Create(&model).Error
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var model User
	result := r.db.Where("email = ?", email).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &domain.User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var model User
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &domain.User{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(user *domain.User) error {
	var model User
	if err := r.db.First(&model, "id = ?", user.ID).Error; err != nil {
		return err
	}

	model.Name = user.Name
	model.Email = user.Email
	model.Password = user.Password
	model.UpdatedAt = time.Now()

	return r.db.Save(&model).Error
}

// Delete removes a user from the database by ID
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
