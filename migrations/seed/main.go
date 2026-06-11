package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("veritabanına bağlanılamadı: %v", err)
	}

	// Тестовый пользователь
	email    := "admin@hisartour.az"
	password := "admin123"
	name     := "Admin"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		INSERT INTO users (name, email, password)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE password = VALUES(password), name = VALUES(name)
	`, name, email, string(hash))
	if err != nil {
		log.Fatalf("kullanıcı oluşturulamadı: %v", err)
	}

	fmt.Println("✓ Test kullanıcısı oluşturuldu:")
	fmt.Printf("  E-posta : %s\n", email)
	fmt.Printf("  Şifre   : %s\n", password)
}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		envOr("DB_USER", "root"),
		envOr("DB_PASSWORD", ""),
		envOr("DB_HOST", "127.0.0.1"),
		envOr("DB_PORT", "3306"),
		envOr("DB_NAME", "accounting"),
	)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
