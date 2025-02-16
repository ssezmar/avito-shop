package models

import (
    "database/sql"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password,omitempty"`
}

func (u *User) Create(db *sql.DB) error {
    query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
    return db.QueryRow(query, u.Name, u.Email, u.Password).Scan(&u.ID)
}

func GetByID(db *sql.DB, id int) (*User, error) {
    var user User
    query := `SELECT id, name, email FROM users WHERE id = $1`
    err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}

func GetByEmail(db *sql.DB, email string) (*User, error) {
    var user User
    query := `SELECT id, name, email, password FROM users WHERE email = $1`
    err := db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
    return &user, err
}

func GetAll(db *sql.DB) ([]User, error) {
    query := `SELECT id, name, email FROM users`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

func (u *User) Update(db *sql.DB) error {
    query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
    _, err := db.Exec(query, u.Name, u.Email, u.ID)
    return err
}

func DeleteByID(db *sql.DB, id int) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := db.Exec(query, id)
    return err
}