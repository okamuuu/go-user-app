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

func (r *UserRepository) Save(user *domain.User) error {
	model := UserModel{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	// IDをuserに反映
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var model UserModel
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

func (r *UserRepository) Update(user *domain.User) error {
	model := UserModel{}
	if err := r.db.First(&model, "id = ?", user.ID).Error; err != nil {
		return err
	}

	model.Name = user.Name
	model.Email = user.Email
	model.Password = user.Password
	model.UpdatedAt = time.Now()

	return r.db.Save(&model).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&UserModel{}, "id = ?", id).Error
}
