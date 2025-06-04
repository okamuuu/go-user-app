package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        *repository.UserRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewAuthService(repo *repository.UserRepository, jwtSecret []byte, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		repo:        repo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// SignUp hashes the password and saves the user
func (s *AuthService) SignUp(user *domain.User) error {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.repo.Save(user)
}

// Login verifies the user credentials and returns a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	// パスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// JWT発行
	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateJWT creates a signed JWT token for the user
func (s *AuthService) GenerateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT verifies the token string and returns the claims if valid
func (s *AuthService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名方法のチェック（HS256想定）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
