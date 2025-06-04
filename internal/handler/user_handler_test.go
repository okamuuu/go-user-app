package handler_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/handler"
	"github.com/okamuuu/go-user-app/internal/middleware"
	"github.com/okamuuu/go-user-app/internal/repository"
	"github.com/okamuuu/go-user-app/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() (*gin.Engine, *gorm.DB, *handler.UserHandler, *service.AuthService) {
	gin.SetMode(gin.TestMode)

	// テスト用DBをメモリ上に作成
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// マイグレーション
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		panic(err)
	}

	// リポジトリ、サービス、ハンドラー作成
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	jwtSecret := []byte("test-secret")
	expireHours := 1000

	authService := service.NewAuthService(userRepo, []byte(jwtSecret), time.Duration(expireHours)*time.Hour)
	userHandler := handler.NewUserHandler(userService)

	// ルーター作成
	r := gin.Default()

	authMiddleware := middleware.AuthMiddleware([]byte(jwtSecret))

	// 認証ミドルウェアを適用したルートグループ
	authorized := r.Group("/")
	authorized.Use(authMiddleware)
	authorized.PUT("/api/users/:id", userHandler.UpdateUser)

	return r, db, userHandler, authService
}

func TestUpdateUser(t *testing.T) {
	r, db, _, authService := setupRouter()

	// 事前にユーザー作成
	user := &domain.User{
		Name:     "Old Name",
		Email:    "old@example.com",
		Password: "oldpassword",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	token, err := authService.GenerateJWT(user)
	if err != nil {
		log.Fatalf("failed to generate token: %v", err)
	}

	// 更新リクエストボディ
	updateData := map[string]string{
		"name":     "New Name",
		"email":    "new@example.com",
		"password": "newpassword",
	}
	jsonBody, _ := json.Marshal(updateData)

	req, _ := http.NewRequest(http.MethodPut, "/api/users/"+strconv.Itoa(int(user.ID)), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// レスポンス取得用Recorder
	w := httptest.NewRecorder()

	// リクエスト処理
	r.ServeHTTP(w, req)

	// ステータスコードチェック
	assert.Equal(t, http.StatusOK, w.Code)

	// DBの値も確認
	var updatedUser domain.User
	if err := db.First(&updatedUser, user.ID).Error; err != nil {
		t.Fatalf("failed to find updated user: %v", err)
	}

	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "new@example.com", updatedUser.Email)
	// パスワードはハッシュ化されているはずなので値は異なる
	assert.NotEqual(t, "newpassword", updatedUser.Password)
}
