package seed

import (
	"log"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/okamuuu/go-user-app/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, count int) {
	for i := 0; i < count; i++ {
		user := domain.User{
			Name:      faker.Name(),
			Email:     faker.Email(),
			Password:  faker.Password(), // 必要なら hash に変換
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user #%d: %v", i+1, err)
		}
	}
	log.Printf("Seeded %d users successfully", count)
}

func hash(pw string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h)
}
