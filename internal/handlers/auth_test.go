package handlers

import (
	"database/sql"
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"encoding/json"
	"avito-shop/internal/models"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"log"
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

func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	jwtSecret := "test_secret"

	handler := Register(db, jwtSecret)

	reqBody := `{"name": "Test User", "email": "test.user@example.com", "password": "password"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := resp["token"]; !ok {
		t.Errorf("Expected token in response")
	}
}

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	jwtSecret := "test_secret"

	user := models.User{
		Name:     "Test User",
		Email:    "test.user@example.com",
		Password: "password",
		Balance:  1000,
	}

	if err := user.Create(db); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	handler := Login(db, jwtSecret)

	reqBody := `{"email": "test.user@example.com", "password": "password"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := resp["token"]; !ok {
		t.Errorf("Expected token in response")
	}
}