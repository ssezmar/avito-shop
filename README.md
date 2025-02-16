# Avito Shop

## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/ssezmar/avito-shop.git
    cd avito-shop
    ```

## Running the Project
1. Start the application using Docker Compose:
    ```bash
    docker-compose up --build
    ```

2. Open your browser and navigate to `http://localhost:8080`

## Running Tests
1. Run tests:
    ```bash
    docker-compose exec avito-shop-service 
    go test ./...
    ```

