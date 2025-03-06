package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"not null;unique"`
	Email     string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null"` // disimpan dalam bentuk hash
	Role      string    `gorm:"not null"` // "admin" atau "member"
	PP        string    `gorm:"not null"` // URL foto profile
	CreatedAt time.Time
}
