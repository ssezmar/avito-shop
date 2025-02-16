FROM golang:1.22

WORKDIR /app
COPY . .

RUN go mod download

# Установка утилиты godotenv для работы с .env файлами
RUN go get github.com/joho/godotenv

# Запуск тестов
CMD ["go", "test", "./..."]