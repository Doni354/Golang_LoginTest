package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Golang_LoginTest/database"
	"Golang_LoginTest/handlers"
	"Golang_LoginTest/models"
)

// Seed data: menambahkan user dummy jika tabel kosong
func seedUser() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Buat user admin
		adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal hash password admin:", err)
		}
		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(adminPassword),
			Role:     "admin",
		}
		database.DB.Create(&admin)

		// Buat user member
		memberPassword, err := bcrypt.GenerateFromPassword([]byte("member"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal hash password member:", err)
		}
		member := models.User{
			Username: "member",
			Email:    "member@example.com",
			Password: string(memberPassword),
			Role:     "member",
		}
		database.DB.Create(&member)

		log.Println("Seeded dummy users: admin & member")
	}
}
func main() {
	// Hubungkan ke database
	database.Connect()
	// Seed data dummy jika tabel kosong
	seedUser()

	// Inisialisasi Gin
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Route untuk login, logout, dan fitur umum
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.ProcessLogin)
	r.GET("/logout", handlers.Logout)

	// Group route yang memerlukan autentikasi
	auth := r.Group("/")
	auth.Use(handlers.AuthRequired)
	{
		auth.GET("/dashboard", handlers.Dashboard)
		auth.GET("/adminfeature", handlers.AdminFeature)
		auth.GET("/memberfeature", handlers.MemberFeature)
		auth.GET("/commonfeature", handlers.CommonFeature)
	}

	// Redirect root ke halaman login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})

	// Jalankan server pada port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
