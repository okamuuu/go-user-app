package service

import (
	"errors"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(id, name, email, password string) (*domain.User, error) {
	// メールアドレスの重複チェック
	existing, _ := s.repo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	// ドメイン層の User を生成
	user, err := domain.NewUser(id, name, email, password)
	if err != nil {
		return nil, err
	}

	// 永続化
	err = s.repo.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
