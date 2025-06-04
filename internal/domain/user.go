package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string // 本当はハッシュ化して扱う想定
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 新しいユーザーを作成するファクトリ関数
func NewUser(id, name, email, password string) (*User, error) {
	if id == "" || name == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}
	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
