# Product Management Backend

A Go (Golang) backend service for a Product Management module, built with Clean Architecture and Domain-Driven Design (DDD) principles.

## Features

- Full CRUD for products.
- PostgreSQL persistence.
- Layered architecture (Handler, Usecase, Repository, Entity).
- Docker & Docker Compose support.
- Configuration via Environment Variables.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose

## Quick Start (Docker)

1.  Clone the repository.
2.  Run the application and database:
    ```bash
    docker-compose up --build
    ```
3.  The API will be available at `http://localhost:8080`.

## Quick Start (Local)

1.  Start a PostgreSQL database (e.g., using the `db` service in docker-compose).
2.  Set environment variables in `.env` or export them:
    ```bash
    export APP_PORT=8080
    export DB_HOST=localhost
    export DB_PORT=5433
    export DB_USER=nofu
    export DB_PASSWORD=nofu2025
    export DB_NAME=nofuproductdb
    ```
3.  Run the application:
    ```bash
    go run cmd/main.go
    ```

## API Endpoints

| Method | Endpoint               | Description           |
| :----- | :--------------------- | :-------------------- |
| POST   | `/api/v1/products`     | Create a new product. |
| GET    | `/api/v1/products`     | Get all products.     |
| GET    | `/api/v1/products/:id` | Get a product by ID.  |
| PUT    | `/api/v1/products/:id` | Update a product.     |
| DELETE | `/api/v1/products/:id` | Delete a product.     |

### Example Request (Create Product)

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Super Widget",
    "description": "The best widget ever created.",
    "price": 19.99
  }'
```

## Project Structure

```
.
├── cmd
│   └── main.go           # Application entrypoint
├── internal
│   ├── config            # Configuration loading
│   ├── dto               # Data Transfer Objects
│   ├── entity            # Domain entities
│   ├── handler           # HTTP handlers (Gin)
│   ├── repository        # Data access layer (PostgreSQL)
│   ├── server            # Server setup
│   └── usecase           # Business logic
├── Dockerfile            # Docker build instructions
├── docker-compose.yml    # Docker Compose configuration
└── go.mod                # Go module file
```
