package main

import (
	"fmt"

	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/service"
)

func main() {
	repo := repository.NewUserRepository()
	uc := service.NewUserService(repo)

	user, err := uc.RegisterUser("1", "Taro Yamada", "taro@example.com", "password123")
	if err != nil {
		panic(err)
	}

	fmt.Printf("User registered: %+v\n", user)
}
