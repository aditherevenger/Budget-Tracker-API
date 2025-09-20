# Budget Tracker API

A comprehensive RESTful API for personal budget tracking built with Go, Gin framework, PostgreSQL, and JWT authentication.

---

## Features

- **User Management:** Registration, login, and profile management  
- **Category Management:** Create, read, update, and delete expense/income categories  
- **Transaction Management:** Track income and expenses with detailed information  
- **Financial Reporting:** Get summaries and insights about your financial data  
- **JWT Authentication:** Secure API endpoints with JSON Web Tokens  
- **PostgreSQL Database:** Reliable data storage with GORM ORM  

---

## Tech Stack

- **Language:** Go 1.23+  
- **Web Framework:** Gin Gonic  
- **Database:** PostgreSQL  
- **ORM:** GORM  
- **Authentication:** JWT (JSON Web Tokens)  
- **Password Hashing:** bcrypt  

---

## Project Structure

```plaintext
budget-tracker-api/
├── main.go                          # Application entry point
├── config/
│   └── config.go                    # Configuration management
├── database/
│   └── connection.go                # Database connection and migration
├── models/
│   ├── user.go                      # User model and DTOs
│   ├── category.go                  # Category model and DTOs
│   └── transaction.go               # Transaction model and DTOs
├── repository/
│   ├── user_repository.go           # User data access layer
│   ├── category_repository.go       # Category data access layer
│   └── transaction_repository.go    # Transaction data access layer
├── services/
│   ├── auth_service.go              # Authentication business logic
│   ├── category_service.go          # Category business logic
│   └── transaction_service.go       # Transaction business logic
├── controllers/
│   ├── auth_controller.go           # Authentication HTTP handlers
│   ├── category_controller.go       # Category HTTP handlers
│   └── transaction_controller.go    # Transaction HTTP handlers
├── middleware/
│   └── auth_middleware.go           # JWT authentication middleware
├── utils/
│   ├── jwt.go                       # JWT token utilities
│   └── password.go                  # Password hashing utilities
├── .env                             # Environment variables
├── go.mod                           # Go module dependencies
├── go.sum                           # Go module checksums
├── README.md                        # Project documentation
```

## Prerequisites

- Go 1.23 or higher
- PostgreSQL 12 or higher
- Git

---

## Installation & Setup

1. Clone the Repository

```bash
git clone <repository-url>
cd budget-tracker-api
```

2. Install Dependencies

```bash
go mod tidy
```

3. Setup PostgreSQL

```bash
CREATE DATABASE budget_tracker;
-- If you want a custom user (optional)
CREATE USER postgres WITH PASSWORD 'AdiKal@1505';
ALTER ROLE postgres SUPERUSER;
GRANT ALL PRIVILEGES ON DATABASE budget_tracker TO postgres;
```

4. Configure Environment Variables
```plaintext
Create .env in the project root:
```

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Password@123
DB_NAME=budget_tracker
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug
```

5. Run the Application
   
```bash
go run main.go
```
The server will start on:

👉 http://localhost:8080

---

## API Endpoints
Authentication
- POST /api/auth/register → Register a new user
-  POST /api/auth/login → Login user
- GET /api/profile → Get user profile (protected)

Categories
- GET /api/categories → Get all categories (protected)
- POST /api/categories → Create a new category (protected)
- GET /api/categories/:id → Get category by ID (protected)
- PUT /api/categories/:id → Update category (protected)
- DELETE /api/categories/:id → Delete category (protected)

Transactions
- GET /api/transactions → Get all transactions with filters (protected)
- POST /api/transactions → Create a new transaction (protected)
- GET /api/transactions/:id → Get transaction by ID (protected)
- PUT /api/transactions/:id → Update transaction (protected)
- DELETE /api/transactions/:id → Delete transaction (protected)
- GET /api/transactions/summary → Get financial summary (protected)

Health Check
- GET /health → Health check endpoint

---

## Auth & Profile

Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","first_name":"John","last_name":"Doe"}'
```

Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'
```

Get Profile (Protected)
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Health Check
```bash
curl http://localhost:8080/health
```

---

## Categories

| Field       | Type   | Description              |
|-------------|--------|--------------------------|
| name        | string | Category name            |
| description | string | Optional description     |
| color       | string | Hex color code           |

List Categories
```bash
curl -X GET http://localhost:8080/api/categories \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Create Category
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Food","description":"Daily meals","color":"#00C853"}'
```

Get Category by ID
```bash
curl -X GET http://localhost:8080/api/categories/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Update Category
```bash
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Groceries","description":"Supermarket & essentials","color":"#FF6D00"}'
```

Delete Category
```bash
curl -X DELETE http://localhost:8080/api/categories/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

## Transactions

| Field       | Type    | Description                     |
|-------------|---------|---------------------------------|
| category_id | integer | ID of category                  |
| amount      | float   | Transaction amount              |
| type        | string  | "income" or "expense"           |
| description | string  | Optional description            |
| date        | string  | ISO 8601 datetime format        |

List Transactions
```bash
curl -X GET http://localhost:8080/api/transactions/ \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Add Transaction
```bash
curl -X POST http://localhost:8080/api/transactions/ \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"amount":149.99,"type":"expense","description":"Supermarket run","date":"2025-09-20T10:30:00Z"}'
```

Get Transaction by ID
```bash
curl -X PUT http://localhost:8080/api/transactions/1 \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"amount":129.49,"description":"Supermarket run (discount applied)"}'
```

Delete Transaction
```bash
curl -X DELETE http://localhost:8080/api/transactions/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

Financial Summary
```bash
curl -X GET "http://localhost:8080/api/transactions/summary?start_date=2025-09-01T00:00:00Z&end_date=2025-09-30T23:59:59Z" \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

---

# Database Schema

## Users Table
- **id** (Primary Key)  
- **email** (Unique)  
- **password** (Hashed)  
- **first_name**  
- **last_name**  
- **created_at**  
- **updated_at**  
- **deleted_at**  

## Categories Table
- **id** (Primary Key)  
- **user_id** (Foreign Key)  
- **name**  
- **description**  
- **color**  
- **created_at**  
- **updated_at**  
- **deleted_at**  

## Transactions Table
- **id** (Primary Key)  
- **user_id** (Foreign Key)  
- **category_id** (Foreign Key)  
- **amount**  
- **type** (income/expense)  
- **description**  
- **date**  
- **created_at**  
- **updated_at**  
- **deleted_at**

---

## License

*This project is licensed under the MIT License.*

