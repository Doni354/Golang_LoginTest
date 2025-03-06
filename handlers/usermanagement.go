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


// ManageUsers menampilkan tabel daftar user untuk admin
func ManageUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving users")
		return
	}
	c.HTML(http.StatusOK, "admin_user_management.html", gin.H{
		"users": users,
	})
}

// ShowCreateUserPage menampilkan form pembuatan user baru oleh admin
func ShowCreateUserPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_user.html", gin.H{})
}

// ProcessCreateUser memproses pembuatan user baru oleh admin
func ProcessCreateUser(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	role := c.PostForm("role") // Pilihan: admin atau member

	// Default foto profile
	defaultPP := "https://static.vecteezy.com/system/resources/previews/005/544/718/non_2x/profile-icon-design-free-vector.jpg"
	var fileURL string

	// Coba ambil file upload (field "pp")
	file, err := c.FormFile("pp")
	if err != nil {
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
			c.HTML(http.StatusInternalServerError, "admin_create_user.html", gin.H{"error": "Failed to save profile picture"})
			return
		}
		fileURL = "/uploads/" + newFileName
	}

	if username == "" || email == "" || password == "" || role == "" {
		c.HTML(http.StatusBadRequest, "admin_create_user.html", gin.H{"error": "Semua field harus diisi"})
		return
	}

	// Cek apakah username atau email sudah ada
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "admin_create_user.html", gin.H{"error": "Username atau email sudah terdaftar"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_user.html", gin.H{"error": "Gagal membuat password"})
		return
	}

	newUser := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		PP:       fileURL,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_user.html", gin.H{"error": "Gagal membuat user"})
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// ShowEditUserPage menampilkan form edit user
func ShowEditUserPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}
	c.HTML(http.StatusOK, "admin_edit_user.html", gin.H{
		"user": user,
	})
}

// ProcessEditUser memproses update data user
func ProcessEditUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	username := c.PostForm("username")
	email := c.PostForm("email")
	role := c.PostForm("role")

	// Update password jika diisi
	password := c.PostForm("password")
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "admin_edit_user.html", gin.H{"error": "Gagal hash password", "user": user})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Proses file foto profile jika diupload
	file, err := c.FormFile("pp")
	if err == nil {
		uploadPath := "./uploads"
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s_%d%s", username, time.Now().Unix(), ext)
		savePath := filepath.Join(uploadPath, newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.HTML(http.StatusInternalServerError, "admin_edit_user.html", gin.H{"error": "Gagal menyimpan foto profile", "user": user})
			return
		}
		user.PP = "/uploads/" + newFileName
	}

	user.Username = username
	user.Email = email
	user.Role = role

	if err := database.DB.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_user.html", gin.H{"error": "Gagal update user", "user": user})
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}

// DeleteUser menghapus user
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Gagal menghapus user")
		return
	}
	c.Redirect(http.StatusFound, "/admin/users")
}
