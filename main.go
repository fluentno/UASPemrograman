package main

import (
	"log"
	"net/http"
	"proyek-karyawan/config"
	"proyek-karyawan/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	router := gin.Default()

	// PERBAIKAN: Muat app.html (untuk layout utama) dan login.html (untuk halaman login)
	router.LoadHTMLFiles("templates/app.html", "templates/login.html")

	router.Static("/static", "./static")

	// --- Rute Publik ---
	router.GET("/login", handlers.ShowLoginPage)
	router.POST("/login", handlers.Login)

	// --- Rute Terotentikasi ---
	authorized := router.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.GET("/", handlers.GetEmployees)
		authorized.GET("/logout", handlers.Logout)
		authorized.GET("/employee/add", handlers.ShowAddEmployeeForm)
		authorized.POST("/employee/add", handlers.CreateEmployee)
		authorized.GET("/employee/edit/:id", handlers.ShowEditEmployeeForm)
		authorized.POST("/employee/edit/:id", handlers.UpdateEmployee)
		authorized.POST("/employee/delete/:id", handlers.DeleteEmployee)
		authorized.GET("/report/pdf", handlers.GeneratePDF)
		authorized.GET("/report/excel", handlers.GenerateExcel)
	}

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	log.Println("Server dimulai pada http://localhost:8080")
	router.Run(":8080")
}
