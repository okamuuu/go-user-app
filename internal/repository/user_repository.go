package repository

import (
	"errors"
	"sync"

	"github.com/okamuuu/go-user-app/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*UserModel
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*UserModel),
	}
}

func (r *UserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = ToUserModel(user)
	return nil
}
