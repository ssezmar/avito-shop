package handlers

import (
    "encoding/json"
    "net/http"
    "avito-shop/internal/models"
    "database/sql"
    "github.com/gorilla/context"
)

func GetMerchList(w http.ResponseWriter, r *http.Request) {
    merchList := models.GetMerchList()
    json.NewEncoder(w).Encode(merchList)
}

func BuyMerch(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var request struct {
            MerchName string `json:"merch_name"`
        }
        if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        email := context.Get(r, "user").(string)
        user, err := models.GetByEmail(db, email)
        if err != nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        merch, err := models.GetMerchByName(request.MerchName)
        if err != nil {
            http.Error(w, "Merch not found", http.StatusNotFound)
            return
        }

        if user.Balance < merch.Price {
            http.Error(w, "Insufficient balance", http.StatusBadRequest)
            return
        }

        user.Balance -= merch.Price
        if err := user.Update(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        transaction := models.Transaction{
            UserID: user.ID,
            Type:   "purchase",
            Amount: merch.Price,
            Note:   "Purchased " + merch.Name,
        }
        if err := transaction.Create(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(user)
    }
}