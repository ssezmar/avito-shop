package main

import (
    "log"
    "net/http"
    "os"
    "avito-shop/internal/handlers"
    "avito-shop/internal/db"
    "avito-shop/internal/auth"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Читаем переменные окружения
    dbHost := os.Getenv("DATABASE_HOST")
    dbPort := os.Getenv("DATABASE_PORT")
    dbUser := os.Getenv("DATABASE_USER")
    dbPassword := os.Getenv("DATABASE_PASSWORD")
    dbName := os.Getenv("DATABASE_NAME")
    jwtSecret := os.Getenv("JWT_SECRET")

    // Инициализация базы данных
    dbConn, err := db.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }
    defer dbConn.Close()

    // Создаем маршрутизатор
    r := mux.NewRouter()

    // Регистрация обработчиков
    r.HandleFunc("/register", handlers.Register(dbConn, jwtSecret)).Methods("POST")
    r.HandleFunc("/login", handlers.Login(dbConn, jwtSecret)).Methods("POST")

    // Маршруты с защитой JWT
    api := r.PathPrefix("/api").Subrouter()
    api.Use(auth.JWTMiddleware(jwtSecret))
    api.HandleFunc("/merch", handlers.GetMerchList).Methods("GET")
    api.HandleFunc("/merch/buy", handlers.BuyMerch(dbConn)).Methods("POST")
    api.HandleFunc("/transactions", handlers.GetTransactions(dbConn)).Methods("GET")
    api.HandleFunc("/transactions/transfer", handlers.TransferCoins(dbConn)).Methods("POST")

    // Запускаем HTTP сервер
    log.Println("Запуск сервера на порту 8080...")
    log.Fatal(http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), r))
}