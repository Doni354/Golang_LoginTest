package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"Golang_LoginTest/database"
	"Golang_LoginTest/handlers"
)

func main() {
	// Hubungkan ke database
	database.Connect()

	// Inisialisasi Gin
	r := gin.Default()
	// Muat semua file template dari folder templates
	r.LoadHTMLGlob("templates/*")

	// Routing untuk login dan dashboard
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.ProcessLogin)

	// Kelompokkan route yang membutuhkan autentikasi
	auth := r.Group("/")
	auth.Use(handlers.AuthRequired)
	{
		auth.GET("/dashboard", handlers.Dashboard)
	}

	// Jika mengakses root, redirect ke login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// Jalankan server pada port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
