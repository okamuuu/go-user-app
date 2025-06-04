package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/okamuuu/go-user-app/internal/handler"
	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&repository.UserModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	r := gin.Default()
	r.SetTrustedProxies(nil)

	h.RegisterRoutes(r)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
