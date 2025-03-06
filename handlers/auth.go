package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"Golang_LoginTest/database"
	"Golang_LoginTest/models"
)

// ShowLoginPage menampilkan halaman login
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// ShowRegistrationPage menampilkan halaman registrasi
func ShowRegistrationPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

// ProcessRegistration memproses registrasi user baru (default role: member)
// Jika file foto tidak diupload, gunakan default image.
func ProcessRegistration(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Default foto profile
	defaultPP := "https://static.vecteezy.com/system/resources/previews/005/544/718/non_2x/profile-icon-design-free-vector.jpg"
	var fileURL string

	// Coba ambil file yang diupload (field "pp")
	file, err := c.FormFile("pp")
	if err != nil {
		// Jika tidak ada file, gunakan default
		fileURL = defaultPP
	} else {
		uploadPath := "./uploads"
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s_%d%s", username, time.Now().Unix(), ext)
		savePath := filepath.Join(uploadPath, newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Gagal menyimpan foto profile"})
			return
		}
		fileURL = "/uploads/" + newFileName
	}

	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Semua field harus diisi"})
		return
	}

	// Cek apakah username atau email sudah terdaftar
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Username atau email sudah terdaftar"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Gagal membuat password"})
		return
	}

	newUser := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "member", // Registrasi default: member
		PP:       fileURL,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Gagal mendaftar"})
		return
	}

	c.HTML(http.StatusOK, "register.html", gin.H{"success": "Registrasi berhasil, silakan login."})
}

// ProcessLogin memproses login user
func ProcessLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Username tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Password salah"})
		return
	}

	c.SetCookie("user_id", strconv.Itoa(int(user.ID)), 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, "/dashboard")
}

// AuthRequired middleware untuk mengecek sesi login
func AuthRequired(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil || userID == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	c.Next()
}

// Dashboard menampilkan halaman dashboard dengan data user
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
		"pp":       user.PP,
	})
}

// AdminFeature hanya dapat diakses oleh admin
func AdminFeature(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	var user models.User
	id, _ := strconv.Atoi(userID)
	if err := database.DB.First(&user, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if user.Role != "admin" {
		c.String(http.StatusForbidden, "Access denied: hanya admin yang dapat mengakses halaman ini")
		return
	}
	c.HTML(http.StatusOK, "admin_feature.html", gin.H{
		"username": user.Username,
		"role":     user.Role,
	})
}

// MemberFeature hanya dapat diakses oleh member

// ShowMemberEditPage menampilkan halaman edit untuk member
func ShowMemberEditPage(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	var user models.User
	id, _ := strconv.Atoi(userID)
	if err := database.DB.First(&user, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	// Hanya member yang boleh mengedit akun sendiri
	if user.Role != "member" {
		c.String(http.StatusForbidden, "Hanya member yang dapat mengedit akun mereka sendiri")
		return
	}
	c.HTML(http.StatusOK, "member_edit.html", gin.H{
		"user": user,
	})
}

// ProcessMemberEdit memproses update data akun member (tidak mengubah role)
func ProcessMemberEdit(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	var user models.User
	id, _ := strconv.Atoi(userID)
	if err := database.DB.First(&user, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if user.Role != "member" {
		c.String(http.StatusForbidden, "Hanya member yang dapat mengedit akun mereka sendiri")
		return
	}

	// Ambil data dari form
	newUsername := c.PostForm("username")
	newEmail := c.PostForm("email")
	newPassword := c.PostForm("password")

	// Proses file foto profile jika diupload
	// Jika tidak, biarkan tetap sama
	fileURL := user.PP
	file, err := c.FormFile("pp")
	if err == nil {
		uploadPath := "./uploads"
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s_%d%s", newUsername, time.Now().Unix(), ext)
		savePath := filepath.Join(uploadPath, newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.HTML(http.StatusInternalServerError, "member_edit.html", gin.H{"error": "Gagal menyimpan foto profile", "user": user})
			return
		}
		fileURL = "/uploads/" + newFileName
	}

	// Update field (role tidak boleh diubah)
	user.Username = newUsername
	user.Email = newEmail
	user.PP = fileURL

	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "member_edit.html", gin.H{"error": "Gagal mengubah password", "user": user})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "member_edit.html", gin.H{"error": "Gagal memperbarui data", "user": user})
		return
	}

	// Redirect ke dashboard setelah update sukses
	c.Redirect(http.StatusFound, "/dashboard")
}

// CommonFeature dapat diakses oleh kedua role
func CommonFeature(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	var user models.User
	id, _ := strconv.Atoi(userID)
	if err := database.DB.First(&user, id).Error; err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "common_feature.html", gin.H{
		"username": user.Username,
		"role":     user.Role,
	})
}

// Logout menghapus sesi login
func Logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
