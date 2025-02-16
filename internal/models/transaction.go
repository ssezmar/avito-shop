package models

import (
    "database/sql"
    "time"
)

type Transaction struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Type      string    `json:"type"`
    Amount    int       `json:"amount"`
    Note      string    `json:"note"`
    Timestamp time.Time `json:"timestamp"`
}

func (t *Transaction) Create(db *sql.DB) error {
    query := `INSERT INTO transactions (user_id, type, amount, note, timestamp) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    return db.QueryRow(query, t.UserID, t.Type, t.Amount, t.Note, time.Now()).Scan(&t.ID)
}

func GetTransactionsByUserID(db *sql.DB, userID int) ([]Transaction, error) {
    query := `SELECT id, user_id, type, amount, note, timestamp FROM transactions WHERE user_id = $1`
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []Transaction
    for rows.Next() {
        var transaction Transaction
        if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Type, &transaction.Amount, &transaction.Note, &transaction.Timestamp); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }

    return transactions, nil
}