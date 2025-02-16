package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
    "avito-shop/internal/models"
    "avito-shop/internal/auth"
    _ "github.com/lib/pq"
    "github.com/gorilla/context"
)

func TestGetMerchList(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/merch", nil)
    rr := httptest.NewRecorder()

    handler := http.HandlerFunc(GetMerchList)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var merchList []models.Merch
    if err := json.NewDecoder(rr.Body).Decode(&merchList); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if len(merchList) != 10 {
        t.Errorf("Expected 10 merch items, got %d", len(merchList))
    }
}

func TestBuyMerch(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    user := models.User{
        Name:     "Test User",
        Email:    "test.user@example.com",
        Password: "password",
        Balance:  1000,
    }

    if err := user.Create(db); err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }

    jwtSecret := "test_secret"
    token, err := auth.GenerateJWT(user.Email, jwtSecret)
    if err != nil {
        t.Fatalf("Failed to generate JWT: %v", err)
    }

    handler := BuyMerch(db)

    reqBody := `{"merch_name": "t-shirt"}`
    req := httptest.NewRequest("POST", "/api/merch/buy", bytes.NewBufferString(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    context.Set(req, "user", user.Email)

    rr := httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var updatedUser models.User
    if err := json.NewDecoder(rr.Body).Decode(&updatedUser); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if updatedUser.Balance != 920 {
        t.Errorf("Expected balance 920, got %d", updatedUser.Balance)
    }
}