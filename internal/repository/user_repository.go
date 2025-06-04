package repository

import (
	"sync"

	"github.com/okamuuu/go-user-app/internal/domain"
)

type UserRepository struct {
	mu sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	model := UserModel{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	return DB.Create(&model).Error
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var model UserModel
	result := DB.Where("email = ?", email).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return &domain.User{
		ID:       model.ID,
		Name:     model.Name,
		Email:    model.Email,
		Password: model.Password,
	}, nil
}
