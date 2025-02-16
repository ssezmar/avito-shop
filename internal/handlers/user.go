package handlers

import (
    "encoding/json"
    "net/http"
    "avito-shop/internal/models"
    "database/sql"
    "github.com/gorilla/context"
)

func GetTransactions(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        email := context.Get(r, "user").(string)
        user, err := models.GetByEmail(db, email)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        transactions, err := models.GetTransactionsByUserID(db, user.ID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(transactions)
    }
}

func TransferCoins(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var request struct {
            RecipientEmail string `json:"recipient_email"`
            Amount         int    `json:"amount"`
        }
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        senderEmail := context.Get(r, "user").(string)
        sender, err := models.GetByEmail(db, senderEmail)
        if err != nil {
            http.Error(w, "Sender not found", http.StatusNotFound)
            return
        }

        recipient, err := models.GetByEmail(db, request.RecipientEmail)
        if err != nil {
            http.Error(w, "Recipient not found", http.StatusNotFound)
            return
        }

        if sender.Balance < request.Amount {
            http.Error(w, "Insufficient balance", http.StatusBadRequest)
            return
        }

        sender.Balance -= request.Amount
        recipient.Balance += request.Amount

        if err := sender.Update(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        if err := recipient.Update(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        senderTransaction := models.Transaction{
            UserID: sender.ID,
            Type:   "transfer",
            Amount: -request.Amount,
            Note:   "Transferred to " + recipient.Email,
        }
        if err := senderTransaction.Create(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        recipientTransaction := models.Transaction{
            UserID: recipient.ID,
            Type:   "transfer",
            Amount: request.Amount,
            Note:   "Received from " + sender.Email,
        }
        if err := recipientTransaction.Create(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(sender)
    }
}