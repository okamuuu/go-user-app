package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	docs "github.com/okamuuu/go-user-app/cmd/docs"
	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/handler"
	"github.com/okamuuu/go-user-app/internal/middleware"
	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/seed"
	"github.com/okamuuu/go-user-app/internal/service"
)

// @title           Go User App API
// @version         1.0
// @description     シンプルなユーザー管理API（Gin + GORM + JWT）
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT形式: Bearer <token>
func main() {
	// .env ファイル読み込み（あれば）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	expireHoursStr := os.Getenv("JWT_EXPIRE_HOURS")
	expireHours, err := strconv.Atoi(expireHoursStr)
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRE_HOURS: %v", err)
	}

	// DB接続
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 起動時に全テーブルをドロップしてから再作成
	// 学習用なので都度DBをリセットしている
	db.Migrator().DropTable(&domain.User{})
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	seed.SeedUsers(db, 100)

	// リポジトリ、サービス、ハンドラー初期化
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, []byte(jwtSecret), time.Duration(expireHours)*time.Hour)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// Ginルーター作成
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	api := r.Group("/api")

	// 認証不要ルート（サインアップ・ログイン）
	api.POST("/signup", authHandler.Signup)
	api.POST("/login", authHandler.Login)

	// 認証必要ルート
	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware([]byte(jwtSecret)))
	authorized.GET("/me", userHandler.Me)

	// ユーザーCRUDルート
	userRoutes := authorized.Group("/users")
	{
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
		userRoutes.GET("", userHandler.GetUsers)
		userRoutes.POST("", userHandler.CreateUser)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to run server: %v", err)
	}
}
