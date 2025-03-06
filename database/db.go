package database

import (
	"log"

	"Golang_LoginTest/models" // Import model untuk AutoMigrate
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect membuka koneksi ke MySQL dan menjalankan AutoMigrate untuk model User
func Connect() {
	// DSN: gunakan user root tanpa password, host 127.0.0.1, dan database golang_test
	dsn := "root@tcp(127.0.0.1:3306)/golang_test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// AutoMigrate model User
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("AutoMigrate gagal:", err)
	}
}
