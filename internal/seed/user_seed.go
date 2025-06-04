package seed

import (
	"log"
	"time"

	"github.com/okamuuu/go-user-app/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	// DB接続
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	users := []repository.User{
		{
			Name:      "Test User 1",
			Email:     "test1@example.com",
			Password:  hash("password1"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Test User 2",
			Email:     "test2@example.com",
			Password:  hash("password2"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("⚠️ Failed to insert user: %v", err)
		}
	}
	log.Println("✅ Seeded users")
}

func hash(pw string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h)
}
