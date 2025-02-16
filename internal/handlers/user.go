package handlers

import (
    "encoding/json"
    "net/http"
    "avito-shop/internal/models"
    "avito-shop/internal/auth"
    "github.com/gorilla/mux"
    "database/sql"
    "strconv"
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

func CreateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := user.Create(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

func GetUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        user, err := models.GetByID(db, id)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "User not found", http.StatusNotFound)
            } else {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }

        json.NewEncoder(w).Encode(user)
    }
}

func GetUsers(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        users, err := models.GetAll(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(users)
    }
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        user.ID = id

        if err := user.Update(db); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(user)
    }
}

func DeleteUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        if err := models.DeleteByID(db, id); err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "User not found", http.StatusNotFound)
            } else {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}