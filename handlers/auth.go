package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Golang_LoginTest/database"
	"Golang_LoginTest/models"
)

// Tampilkan halaman login
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// Proses form login
func ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Cari user berdasarkan username
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Username tidak ditemukan"})
		return
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Password salah"})
		return
	}

	// Simpan user_id di cookie sebagai tanda sesi login
	c.SetCookie("user_id", strconv.Itoa(int(user.ID)), 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/dashboard")
}

// Middleware untuk mengecek autentikasi
func AuthRequired(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil || userID == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

// Halaman dashboard (setelah login)
// Mengirim data role sehingga template bisa menampilkan fitur sesuai role
func Dashboard(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	var user models.User
	id, _ := strconv.Atoi(userID)
	if err := database.DB.First(&user, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}

// Logout: hapus cookie sesi dan redirect ke halaman login
func Logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
