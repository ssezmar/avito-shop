package models

import (
    "database/sql"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password,omitempty"`
    Balance  int    `json:"balance"`
}

func (u *User) Create(db *sql.DB) error {
    query := `INSERT INTO users (name, email, password, balance) VALUES ($1, $2, $3, $4) RETURNING id`
    return db.QueryRow(query, u.Name, u.Email, u.Password, u.Balance).Scan(&u.ID)
}

func GetByID(db *sql.DB, id int) (*User, error) {
    var user User
    query := `SELECT id, name, email, balance FROM users WHERE id = $1`
    err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Balance)
    return &user, err
}

func GetByEmail(db *sql.DB, email string) (*User, error) {
    var user User
    query := `SELECT id, name, email, password, balance FROM users WHERE email = $1`
    err := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Balance)
    return &user, err
}

func (u *User) Update(db *sql.DB) error {
    query := `UPDATE users SET name = $1, email = $2, password = $3, balance = $4 WHERE id = $5`
    _, err := db.Exec(query, u.Name, u.Email, u.Password, u.Balance, u.ID)
    return err
}