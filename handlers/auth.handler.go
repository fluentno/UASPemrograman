package handlers

import (
	"database/sql"
	"net/http"
	"proyek-karyawan/config"
	"proyek-karyawan/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("super-secret-key"))

func ShowLoginPage(c *gin.Context) {
	// PERBAIKAN: Tampilkan halaman login.html secara langsung
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	err := config.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			// Kode ini sekarang akan bekerja karena login.html sudah dimuat
			c.HTML(http.StatusOK, "login.html", gin.H{"error": "Username tidak ditemukan"})
			return
		}
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Terjadi kesalahan server"})
		return
	}

	if user.Password != password {
		// Kode ini sekarang akan bekerja
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Password salah"})
		return
	}

	session, _ := Store.Get(c.Request, "session-karyawan")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

func Logout(c *gin.Context) {
	session, _ := Store.Get(c.Request, "session-karyawan")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/login")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := Store.Get(c.Request, "session-karyawan")
		if auth, ok := session.Values["user_id"].(int); !ok || auth == 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
