package repository

import "time"

type User struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
