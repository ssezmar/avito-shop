package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"
    _ "github.com/lib/pq"
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

    // Выполняем миграции
    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return db, nil
}

func runMigrations(db *sql.DB) error {
    migrationDir := "migrations"
    files, err := os.ReadDir(migrationDir)
    if err != nil {
        return err
    }

    for _, file := range files {
        if filepath.Ext(file.Name()) == ".sql" {
            path := filepath.Join(migrationDir, file.Name())
            content, err := os.ReadFile(path)
            if err != nil {
                return err
            }

            if _, err := db.Exec(string(content)); err != nil {
                return err
            }
            log.Printf("Migration %s applied", file.Name())
        }
    }
    return nil
}