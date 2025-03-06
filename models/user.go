package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"not null;unique"`
	Email     string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null"` // Password disimpan dalam bentuk hash
	CreatedAt time.Time
}
