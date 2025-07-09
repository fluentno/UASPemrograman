package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// Format: username:password@tcp(host:port)/dbname
	dsn := "root:@tcp(127.0.0.1:3306)/db_karyawan?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}

	fmt.Println("Berhasil terhubung ke database!")
}