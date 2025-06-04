package repository

import "github.com/okamuuu/go-user-app/internal/domain"

// ドメインモデル → DBモデル
func ToUserModel(u *domain.User) *UserModel {
	return &UserModel{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// DBモデル → ドメインモデル
func ToDomainUser(um *UserModel) *domain.User {
	return &domain.User{
		ID:        um.ID,
		Name:      um.Name,
		Email:     um.Email,
		Password:  um.Password,
		CreatedAt: um.CreatedAt,
		UpdatedAt: um.UpdatedAt,
	}
}
