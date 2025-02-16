package models

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestCreateTransaction(t *testing.T) {
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

	transaction := Transaction{
		UserID: user.ID,
		Type:   "purchase",
		Amount: 100,
		Note:   "Purchased item",
	}

	if err := transaction.Create(db); err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}
}

func TestGetTransactionsByUserID(t *testing.T) {
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

	transaction := Transaction{
		UserID: user.ID,
		Type:   "purchase",
		Amount: 100,
		Note:   "Purchased item",
	}

	if err := transaction.Create(db); err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	transactions, err := GetTransactionsByUserID(db, user.ID)
	if err != nil {
		t.Fatalf("Failed to get transactions by user ID: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(transactions))
	}

	if transactions[0].Note != transaction.Note {
		t.Errorf("Expected note %s, got %s", transaction.Note, transactions[0].Note)
	}
}
