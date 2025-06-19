package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" 
	// "github.com/lib/pq"
	_ "github.com/lib/pq" 
)

// Database instance

var DB *sql.DB
func InitDB(){
	var err error
	
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v. Pastikan file .env ada di root proyek.", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Gagal membuka koneksi database PostgreSQL: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database PostgreSQL: %v", err)
	}

	fmt.Println("Koneksi database PostgreSQL berhasil!")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Error saat menutup koneksi database: %v", err)
		}
		fmt.Println("Koneksi database ditutup.")
	}
}