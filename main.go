package main

import (
	"fmt"
	"log"

	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/service"
)

func main() {
	repository.InitDB()

	repo := repository.NewUserRepository()
	userService := service.NewUserService(repo)

	// ユーザー登録の例
	user, err := userService.RegisterUser("Taro Yamada", "taro@example.com", "password123")
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}
	fmt.Printf("User registered: %+v\n", user)

	// ユーザー検索の例
	foundUser, err := userService.FindUserByEmail("taro@example.com")
	if err != nil {
		log.Fatalf("User not found: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)
}
