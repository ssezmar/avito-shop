package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"avito-shop/internal/models"
	"avito-shop/internal/auth"
	"log"
)

func Register(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := user.Create(db); err != nil {
			log.Printf("Failed to create user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateJWT(user.Email, jwtSecret)
		if err != nil {
			log.Printf("Failed to generate token: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func Login(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := models.GetByEmail(db, creds.Email)
		if err != nil {
			log.Printf("Failed to get user by email: %v", err)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if !auth.CheckPasswordHash(creds.Password, user.Password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateJWT(user.Email, jwtSecret)
		if err != nil {
			log.Printf("Failed to generate token: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}