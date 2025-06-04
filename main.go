package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/okamuuu/go-user-app/internal/domain"
	"github.com/okamuuu/go-user-app/internal/repository"
)

func main() {
	// DB接続
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&repository.UserModel{})

	repo := repository.NewUserRepository(db)

	r := gin.Default()

	// ユーザー一覧取得
	r.GET("/users", func(c *gin.Context) {
		var users []domain.User

		var models []repository.UserModel
		if err := db.Find(&models).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, m := range models {
			users = append(users, domain.User{
				ID:        m.ID,
				Name:      m.Name,
				Email:     m.Email,
				Password:  m.Password,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			})
		}
		c.JSON(http.StatusOK, users)
	})

	// ユーザー登録
	r.POST("/users", func(c *gin.Context) {
		var user domain.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.Save(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	})

	// ユーザー更新
	r.PUT("/users/:id", func(c *gin.Context) {
		var user domain.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// URLパラメータのIDをuintに変換してuser.IDにセット
		// 省略する場合は検討が必要です

		if err := repo.Update(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// ユーザー削除
	r.DELETE("/users/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		idUint64, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}
		id := uint(idUint64)

		if err := repo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	r.Run(":8080")
}
