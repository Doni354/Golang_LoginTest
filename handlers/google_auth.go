package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"Golang_LoginTest/database"
	"Golang_LoginTest/models"
	"strconv"
	"github.com/joho/godotenv"
)

var (
	// Pastikan variabel ini diinisialisasi setelah load .env
	googleOauthConfig *oauth2.Config
	// Untuk demo, kita gunakan state string statis. Di produksi, gunakan mekanisme random dan simpan di session.
	oauthStateString = "pseudo-random-state"
)

func init() {
	// Jika kamu menggunakan godotenv, muat file .env
	godotenv.Load()
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleLogin mengarahkan pengguna ke halaman login Google
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback menangani callback dari Google setelah login
func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		fmt.Println("invalid oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("code exchange failed:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// Ambil data user dari endpoint userinfo Google
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("failed getting user info:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("failed reading response body:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
		Name          string `json:"name"`
	}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		fmt.Println("failed to unmarshal user info:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// Cek apakah user dengan email tersebut sudah ada
	var user models.User
	if err := database.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		// Jika belum ada, buat user baru dengan role "member"
		user = models.User{
			Username: userInfo.Name,
			Email:    userInfo.Email,
			Password: "", // Password kosong, karena login via Google
			Role:     "member",
			PP:       userInfo.Picture,
		}
		database.DB.Create(&user)
	}

	// Set cookie agar user dianggap sudah login
	c.SetCookie("user_id", strconv.Itoa(int(user.ID)), 3600, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}
