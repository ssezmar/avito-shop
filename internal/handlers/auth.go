package handlers

import (
    "encoding/json"
    "net/http"
    "avito-shop/internal/models"
    "avito-shop/internal/auth"
    "database/sql"
    "golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB, jwtSecret string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        user.Password = string(hashedPassword)
        user.Balance = 1000

        if err := user.Create(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        token, err := auth.GenerateJWT(user.Email, jwtSecret)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}

func Login(db *sql.DB, jwtSecret string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds models.User
        if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        user, err := models.GetByEmail(db, creds.Email)
        if err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        token, err := auth.GenerateJWT(user.Email, jwtSecret)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}