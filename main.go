package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Golang_LoginTest/database"
	"Golang_LoginTest/handlers"
	"Golang_LoginTest/models"
)

// Seed data dummy: tambahkan user dummy jika tabel kosong
func seedUser() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Gunakan default profile picture dari URL
		defaultPP := "https://static.vecteezy.com/system/resources/previews/005/544/718/non_2x/profile-icon-design-free-vector.jpg"
		// User Admin
		adminPassword, err := bcrypt.GenerateFromPassword([]byte("wasd"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal hash password admin:", err)
		}
		admin := models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(adminPassword),
			Role:     "admin",
			PP:       defaultPP,
		}
		database.DB.Create(&admin)

		// User Member
		memberPassword, err := bcrypt.GenerateFromPassword([]byte("wasd"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal hash password member:", err)
		}
		member := models.User{
			Username: "member",
			Email:    "member@example.com",
			Password: string(memberPassword),
			Role:     "member",
			PP:       defaultPP,
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
	// Muat template HTML
	r.LoadHTMLGlob("templates/*")
	// Sajikan file statis untuk folder uploads
	r.Static("/uploads", "./uploads")

	// Route untuk login, registrasi, logout, dashboard, dan fitur
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.ProcessLogin)
	r.GET("/register", handlers.ShowRegistrationPage)
	r.POST("/register", handlers.ProcessRegistration)
	r.GET("/logout", handlers.Logout)

	// Route untuk user management (hanya untuk admin)
	admin := r.Group("/admin")
	admin.Use(handlers.AuthRequired)
	{
		// Halaman utama admin feature (dashboard admin)
		admin.GET("/dashboard", handlers.Dashboard)
		// Fitur user management
		admin.GET("/users", handlers.ManageUsers)
		admin.GET("/users/create", handlers.ShowCreateUserPage)
		admin.POST("/users/create", handlers.ProcessCreateUser)
		admin.GET("/users/edit/:id", handlers.ShowEditUserPage)
		admin.POST("/users/edit/:id", handlers.ProcessEditUser)
		admin.GET("/users/delete/:id", handlers.DeleteUser)
		

		
	}
	// Fitur-fitur lain (AdminFeature, dsb.) jika diperlukan
	
		// Group route yang memerlukan autentikasi
		auth := r.Group("/")
		auth.Use(handlers.AuthRequired)
		{
			auth.GET("/dashboard", handlers.Dashboard)
			auth.GET("/member/edit", handlers.ShowMemberEditPage)
			auth.POST("/member/edit", handlers.ProcessMemberEdit)
			auth.GET("/adminfeature", handlers.AdminFeature)
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