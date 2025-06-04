package main

import (
	"log"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// マイグレーション実行
	if err := db.AutoMigrate(&repository.UserModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	repo := repository.NewUserRepository(db)

	// アプリの処理（例：ユーザー作成）
	user := &domain.User{
		Name:  "Sample",
		Email: "sample@example.com",
	}
	if err := repo.Save(user); err != nil {
		log.Fatalf("failed to save user: %v", err)
	}
}
