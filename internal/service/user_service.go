package service

import (
	"fmt"
	"log"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers(page, limit int) ([]*domain.User, error) {
	offset := (page - 1) * limit
	return s.repo.FindAll(offset, limit)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *domain.User) error {
	// 例: 重複チェック（emailユニークの場合）
	existing, _ := s.repo.FindByEmail(user.Email)
	if existing != nil {
		return fmt.Errorf("email already exists")
	}
	return s.repo.Create(user)
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

	existingUser, err := s.repo.FindByID(user.ID)
	log.Printf("[DEBUG] Updating user: %+v", existingUser)

	if err != nil {
		return err
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	}

	return s.repo.Update(user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
