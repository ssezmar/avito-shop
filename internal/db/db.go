package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
)

func InitDB(host, port, user, password, dbname string) (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}