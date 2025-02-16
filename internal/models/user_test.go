package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	psqlInfo := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Очистка таблиц перед каждым тестом
	_, err = db.Exec(`TRUNCATE TABLE transactions, users RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}

	log.Println("Database tables truncated successfully")

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := User{
		Name:     "Test User",
		Email:    fmt.Sprintf("test.user+%d@example.com", time.Now().UnixNano()),
		Password: "password",
		Balance:  1000,
	}

	if err := user.Create(db); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	user := User{
		Name:     "Test User",
		Email:    fmt.Sprintf("test.user+%d@example.com", time.Now().UnixNano()),
		Password: "password",
		Balance:  1000,
	}

	if err := user.Create(db); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userFromDB, err := GetByEmail(db, user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}

	if userFromDB.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, userFromDB.Email)
	}
}
