package service

import (
	"fmt"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *domain.User) error {
	// 例: 重複チェック（emailユニークの場合）
	existing, _ := s.repo.FindByEmail(user.Email)
	if existing != nil {
		return fmt.Errorf("email already exists")
	}
	return s.repo.Save(user)
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

// GetUserByEmail fetches user by email
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.repo.FindByEmail(email)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *domain.User) error {
	// 存在確認など必要に応じて
	return s.repo.Update(user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
