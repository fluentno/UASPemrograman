package models

import "time"

type User struct {
	ID       int
	Username string
	Password string
}

type Employee struct {
	ID           int
	Nama         string
	Posisi       string
	Email        string
	Telepon      string
	TanggalMasuk string // Menggunakan string untuk kemudahan format
	CreatedAt    time.Time
}
