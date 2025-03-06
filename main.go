package main

import (
	"log"
	"net/http"


	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Golang_LoginTest/database"
	"Golang_LoginTest/handlers"
	"Golang_LoginTest/models"
)

// seedUser: Menambahkan data user default jika tabel masih kosong
func seedUser() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		password := "john"
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		user := models.User{
			Username: "john",
			Email:    "john@example.com",
			Password: string(hash),
		}
		result := database.DB.Create(&user)
		if result.Error != nil {
			log.Fatal("Failed to seed user:", result.Error)
		}
		log.Println("Seeded default user: john / password123")
	}
}

func main() {
	// Hubungkan ke database
	database.Connect()
	seedUser()
	// Inisialisasi Gin
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Routing untuk login, logout, dan dashboard
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.ProcessLogin)
	r.GET("/logout", handlers.Logout)  // Route logout

	// Group route yang memerlukan autentikasi
	auth := r.Group("/")
	auth.Use(handlers.AuthRequired)
	{
		auth.GET("/dashboard", handlers.Dashboard)
	}

	// Redirect root ke halaman login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// Jalankan server pada port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
