package repository

import (
	"errors"
	"sync"

	"github.com/okamuuu/go-user-app/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepository() *UserRepository {

	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Save(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
