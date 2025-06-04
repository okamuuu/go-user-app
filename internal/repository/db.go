package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// User テーブルがまだなければ自動で作成する
	err = db.AutoMigrate(&UserModel{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
