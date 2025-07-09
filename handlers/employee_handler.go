package handlers

import (
	"fmt"
	"log"
	"net/http"
	"proyek-karyawan/config"
	"proyek-karyawan/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// GetEmployees menampilkan semua karyawan di dasbor.
func GetEmployees(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, posisi, email, telepon, DATE_FORMAT(tanggal_masuk, '%Y-%m-%d') FROM employees ORDER BY id DESC")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: Gagal mengambil data karyawan.")
		return
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.Nama, &emp.Posisi, &emp.Email, &emp.Telepon, &emp.TanggalMasuk); err != nil {
			log.Println("Gagal memindai data karyawan:", err)
			continue
		}
		employees = append(employees, emp)
	}

	session, _ := Store.Get(c.Request, "session-karyawan") // Menggunakan Store
	username := session.Values["username"]

	c.HTML(http.StatusOK, "app.html", gin.H{
		"content_name": "dashboard",
		"employees":    employees,
		"username":     username,
	})
}

// ShowAddEmployeeForm menampilkan formulir untuk menambah karyawan baru.
func ShowAddEmployeeForm(c *gin.Context) {
	session, _ := Store.Get(c.Request, "session-karyawan") // Menggunakan Store
	username := session.Values["username"]

	c.HTML(http.StatusOK, "app.html", gin.H{
		"content_name": "add_employee",
		"username":     username,
	})
}

// CreateEmployee menyimpan data karyawan baru ke database.
func CreateEmployee(c *gin.Context) {
	nama := c.PostForm("nama")
	posisi := c.PostForm("posisi")
	email := c.PostForm("email")
	telepon := c.PostForm("telepon")
	tanggal_masuk := c.PostForm("tanggal_masuk")
	_, err := config.DB.Exec("INSERT INTO employees (nama, posisi, email, telepon, tanggal_masuk) VALUES (?, ?, ?, ?, ?)",
		nama, posisi, email, telepon, tanggal_masuk)
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal menyimpan data karyawan baru.")
		return
	}
	c.Redirect(http.StatusFound, "/")
}

// ShowEditEmployeeForm menampilkan formulir untuk mengedit data karyawan.
func ShowEditEmployeeForm(c *gin.Context) {
	id := c.Param("id")
	log.Println("Akses form edit untuk ID:", id)

	if id == "" {
		c.String(http.StatusBadRequest, "ID tidak valid.")
		return
	}

	var emp models.Employee
	err := config.DB.QueryRow("SELECT id, nama, posisi, email, telepon, DATE_FORMAT(tanggal_masuk, '%Y-%m-%d') FROM employees WHERE id = ?", id).
		Scan(&emp.ID, &emp.Nama, &emp.Posisi, &emp.Email, &emp.Telepon, &emp.TanggalMasuk)

	if err != nil {
		log.Println("Gagal mengambil data untuk ID:", id, "Error:", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	session, _ := Store.Get(c.Request, "session-karyawan") // Menggunakan Store
	username := session.Values["username"]

	c.HTML(http.StatusOK, "app.html", gin.H{
		"content_name": "edit_employee",
		"employee":     emp,
		"username":     username,
	})
}

// UpdateEmployee memperbarui data karyawan di database.
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	log.Println("Proses update untuk ID:", id)

	if id == "" {
		c.String(http.StatusBadRequest, "ID tidak valid.")
		return
	}

	nama := c.PostForm("nama")
	posisi := c.PostForm("posisi")
	email := c.PostForm("email")
	telepon := c.PostForm("telepon")
	tanggal_masuk := c.PostForm("tanggal_masuk")

	_, err := config.DB.Exec("UPDATE employees SET nama=?, posisi=?, email=?, telepon=?, tanggal_masuk=? WHERE id=?",
		nama, posisi, email, telepon, tanggal_masuk, id)
	if err != nil {
		log.Println("Gagal mengupdate data untuk ID:", id, "Error:", err)
		c.String(http.StatusInternalServerError, "Gagal mengupdate data karyawan.")
		return
	}
	c.Redirect(http.StatusFound, "/")
}

// DeleteEmployee menghapus data karyawan dari database.
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	_, err := config.DB.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal menghapus data.")
		return
	}
	c.Redirect(http.StatusFound, "/")
}

// GeneratePDF membuat laporan data karyawan dalam format PDF.
func GeneratePDF(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, posisi, email, telepon FROM employees ORDER BY id ASC")
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal query data: %v", err)
		return
	}
	defer rows.Close()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Laporan Data Karyawan")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"ID", "Nama", "Posisi", "Email", "Telepon"}
	widths := []float64{10, 50, 40, 60, 30}
	for i, header := range headers {
		pdf.CellFormat(widths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.Nama, &emp.Posisi, &emp.Email, &emp.Telepon); err != nil {
			log.Println(err)
			continue
		}
		pdf.CellFormat(widths[0], 6, strconv.Itoa(emp.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 6, emp.Nama, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[2], 6, emp.Posisi, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[3], 6, emp.Email, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[4], 6, emp.Telepon, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=laporan_karyawan.pdf")
	if err := pdf.Output(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Gagal generate PDF")
	}
}

// GenerateExcel membuat laporan data karyawan dalam format Excel.
func GenerateExcel(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, nama, posisi, email, telepon, tanggal_masuk FROM employees ORDER BY id ASC")
	if err != nil {
		c.String(http.StatusInternalServerError, "Gagal query data: %v", err)
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	sheetName := "Sheet1"

	headers := []string{"ID", "Nama", "Posisi", "Email", "Telepon", "Tanggal Masuk"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	rowNum := 2
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.Nama, &emp.Posisi, &emp.Email, &emp.Telepon, &emp.TanggalMasuk); err != nil {
			log.Println(err)
			continue
		}
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), emp.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), emp.Nama)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), emp.Posisi)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), emp.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), emp.Telepon)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), emp.TanggalMasuk)
		rowNum++
	}

	c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=laporan_karyawan.xlsx")
	if err := f.Write(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Gagal generate Excel")
	}
}
