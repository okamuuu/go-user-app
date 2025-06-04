package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/handler"
	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/service"
)

func main() {
	// .env ファイルから環境変数読み込み（任意）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	expireHoursStr := os.Getenv("JWT_EXPIRE_HOURS")

	expireHours, err := strconv.Atoi(expireHoursStr)

	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRE_HOURS: %v", err)
	}
	// DB接続（SQLiteの例）
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// マイグレーション（UserModel構造体をDBに反映）
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// リポジトリ・サービス・ハンドラーの初期化
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, []byte(jwtSecret), time.Duration(expireHours)*time.Hour)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// Ginルーター作成
	r := gin.Default()

	// APIルーティング設定
	api := r.Group("/api")

	// ユーザーCRUDルート
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	// 認証ルート
	api.POST("/signup", authHandler.Signup)
	api.POST("/login", authHandler.Login)

	// サーバー起動（ポートは環境変数PORTで指定可能。指定なければ8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
