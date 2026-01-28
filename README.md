# Kasir API 2 - POS System API

A comprehensive Point of Sale (POS) API built with Go, featuring product and category management with full CRUD operations, input validation, and PostgreSQL database integration.

## 📋 Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Server](#running-the-server)
- [API Documentation](#api-documentation)
- [Testing the API](#testing-the-api)
- [Project Structure](#project-structure)
- [Error Handling](#error-handling)
- [Validation Rules](#validation-rules)

## ✨ Features

- **Product Management**: Create, read, update, and delete products
- **Category Management**: Manage product categories
- **Input Validation**: Comprehensive validation on all endpoints
- **PostgreSQL Database**: Reliable data persistence
- **RESTful API**: Standard HTTP methods and status codes
- **Health Check**: Monitor API status
- **Error Handling**: Detailed error messages and appropriate HTTP status codes

## 📦 Prerequisites

- Go 1.25.5 or higher
- PostgreSQL 10 or higher
- Git

## 🚀 Installation

### 1. Clone the Repository
```bash
git clone <repository-url>
cd kasir-api-2
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
DB_CONN=postgres://username:password@localhost:5432/kasir_db?sslmode=disable
```

**Available Configuration:**
- `PORT`: Server port (default: 8080)
- `DB_CONN`: PostgreSQL connection string

## 🗄️ Database Setup

### 1. Create Database
```sql
CREATE DATABASE kasir_db;
```

### 2. Create Tables

Connect to your database and run:

```sql
-- Categories Table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products Table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_categories_name ON categories(name);
```

## ▶️ Running the Server

### Development Mode
```bash
go run .
```

### Build and Run
```bash
go build -o kasir-api
./kasir-api
```

Expected output:
```
Server running di localhost:8080
Database connected successfully
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080
```

### Health Check
Check if the API is running:

**Request:**
```bash
GET /health
```

**Response (200 OK):**
```json
{
  "status": "OK",
  "message": "API Running"
}
```

---

## 🏷️ Category Endpoints

### 1. Get All Categories
**Request:**
```bash
GET /api/categories
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets"
  },
  {
    "id": 2,
    "name": "Books",
    "description": "Physical and digital books"
  }
]
```

---

### 2. Get Category by ID
**Request:**
```bash
GET /api/categories/1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

**Response (404 Not Found):**
```json
Category not found
```

---

### 3. Create Category
**Request:**
```bash
POST /api/categories
Content-Type: application/json

{
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

**Validation Rules:**
- `name`: Required, 3-100 characters
- `description`: Optional, max 500 characters

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

**Response (400 Bad Request):**
```json
Category name is required
```

---

### 4. Update Category
**Request:**
```bash
PUT /api/categories/1
Content-Type: application/json

{
  "name": "Electronic Devices",
  "description": "Updated description"
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Electronic Devices",
  "description": "Updated description"
}
```

---

### 5. Delete Category
**Request:**
```bash
DELETE /api/categories/1
```

**Response (200 OK):**
```json
{
  "message": "Category deleted successfully"
}
```

---

## 📦 Product Endpoints

### 1. Get All Products
**Request:**
```bash
GET /api/products
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Laptop",
    "price": 15000000,
    "stock": 10,
    "category_id": 1,
    "category_name": "Electronics"
  },
  {
    "id": 2,
    "name": "Mouse",
    "price": 500000,
    "stock": 50,
    "category_id": 1,
    "category_name": "Electronics"
  }
]
```

---

### 2. Get Product by ID
**Request:**
```bash
GET /api/products/1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Laptop",
  "price": 15000000,
  "stock": 10,
  "category_id": 1,
  "category_name": "Electronics"
}
```

**Response (404 Not Found):**
```json
Product not found
```

---

### 3. Create Product
**Request:**
```bash
POST /api/products
Content-Type: application/json

{
  "name": "Laptop",
  "price": 15000000,
  "stock": 10,
  "category_id": 1
}
```

**Validation Rules:**
- `name`: Required, 3-255 characters
- `price`: Required, cannot be negative
- `stock`: Required, cannot be negative
- `category_id`: Required, must be greater than 0

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Laptop",
  "price": 15000000,
  "stock": 10,
  "category_id": 1,
  "category_name": "Electronics"
}
```

**Response (400 Bad Request):**
```json
Product name is required
```

---

### 4. Update Product
**Request:**
```bash
PUT /api/products/1
Content-Type: application/json

{
  "name": "Gaming Laptop",
  "price": 18000000,
  "stock": 8,
  "category_id": 1
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Gaming Laptop",
  "price": 18000000,
  "stock": 8,
  "category_id": 1,
  "category_name": "Electronics"
}
```

---

### 5. Delete Product
**Request:**
```bash
DELETE /api/products/1
```

**Response (200 OK):**
```json
{
  "message": "Product deleted successfully"
}
```

---

## 🧪 Testing the API

### Using cURL

#### Test Health Endpoint
```bash
curl -X GET http://localhost:8080/health
```

#### Create a Category
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices"
  }'
```

#### Get All Categories
```bash
curl -X GET http://localhost:8080/api/categories
```

#### Get Category by ID
```bash
curl -X GET http://localhost:8080/api/categories/1
```

#### Update Category
```bash
curl -X PUT http://localhost:8080/api/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Electronics",
    "description": "Updated description"
  }'
```

#### Delete Category
```bash
curl -X DELETE http://localhost:8080/api/categories/1
```

#### Create a Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 15000000,
    "stock": 10,
    "category_id": 1
  }'
```

#### Get All Products
```bash
curl -X GET http://localhost:8080/api/products
```

#### Get Product by ID
```bash
curl -X GET http://localhost:8080/api/products/1
```

#### Update Product
```bash
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptop",
    "price": 18000000,
    "stock": 8,
    "category_id": 1
  }'
```

#### Delete Product
```bash
curl -X DELETE http://localhost:8080/api/products/1
```

---

### Using Postman

1. **Import Collection:**
   - Open Postman
   - Create a new collection "Kasir API"
   - Add the endpoints below

2. **Set Base URL Variable:**
   - Go to Variables
   - Set `base_url` = `http://localhost:8080`

3. **Create Requests:**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `{{base_url}}/health` | Health check |
| GET | `{{base_url}}/api/categories` | Get all categories |
| GET | `{{base_url}}/api/categories/1` | Get category by ID |
| POST | `{{base_url}}/api/categories` | Create category |
| PUT | `{{base_url}}/api/categories/1` | Update category |
| DELETE | `{{base_url}}/api/categories/1` | Delete category |
| GET | `{{base_url}}/api/products` | Get all products |
| GET | `{{base_url}}/api/products/1` | Get product by ID |
| POST | `{{base_url}}/api/products` | Create product |
| PUT | `{{base_url}}/api/products/1` | Update product |
| DELETE | `{{base_url}}/api/products/1` | Delete product |

---

## 📁 Project Structure

```
kasir-api-2/
├── main.go                      # Application entry point
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── README.md                    # This file
├── .env                         # Environment variables (create this)
├── database/
│   └── database.go             # Database connection setup
├── models/
│   ├── category.go             # Category model
│   └── product.go              # Product model
├── repositories/
│   ├── category_repository.go  # Category data access layer
│   └── product_repository.go   # Product data access layer
├── services/
│   ├── category_service.go     # Category business logic
│   └── product_service.go      # Product business logic
└── handlers/
    ├── category_handler.go     # Category HTTP handlers
    └── product_handler.go      # Product HTTP handlers
```

### Architecture Overview

```
HTTP Request
    ↓
Handler (Validation)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Database (PostgreSQL)
```

---

## ✅ Validation Rules

### Category Validation

| Field | Required | Min Length | Max Length | Rules |
|-------|----------|-----------|-----------|-------|
| name | Yes | 3 | 100 | Must be unique |
| description | No | - | 500 | - |

### Product Validation

| Field | Required | Min Value | Max Value | Rules |
|-------|----------|-----------|-----------|-------|
| name | Yes | 3 chars | 255 chars | - |
| price | Yes | 0 | No limit | Cannot be negative |
| stock | Yes | 0 | No limit | Cannot be negative |
| category_id | Yes | 1 | No limit | Must exist in categories |

---

## 🔴 Error Handling

### HTTP Status Codes

| Status Code | Meaning | Example |
|------------|---------|---------|
| 200 | OK | Successful GET, PUT, DELETE |
| 201 | Created | Successful POST |
| 400 | Bad Request | Invalid input, validation error |
| 404 | Not Found | Resource not found |
| 405 | Method Not Allowed | Wrong HTTP method |
| 500 | Internal Server Error | Database error |

### Error Response Format

```json
{
  "error": "Error message description"
}
```

### Common Validation Errors

**Category Errors:**
```
"Category name is required"
"Category name must be at least 3 characters"
"Category name must not exceed 100 characters"
"Category description must not exceed 500 characters"
```

**Product Errors:**
```
"Product name is required"
"Product name must be at least 3 characters"
"Product name must not exceed 100 characters"
"Product price cannot be negative"
"Product stock cannot be negative"
"Category ID must be greater than 0"
```

---

## 🔧 Development

### Adding New Endpoints

1. Create a model in `models/`
2. Create a repository in `repositories/`
3. Create a service in `services/`
4. Create handlers in `handlers/`
5. Register routes in `main.go`

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./handlers
```

---

## 📝 Example Workflow

### 1. Create a Category
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics", "description": "Electronic devices"}'
```

Response:
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices"
}
```

### 2. Create a Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 15000000,
    "stock": 10,
    "category_id": 1
  }'
```

Response:
```json
{
  "id": 1,
  "name": "Laptop",
  "price": 15000000,
  "stock": 10,
  "category_id": 1,
  "category_name": "Electronics"
}
```

### 3. Get All Products
```bash
curl -X GET http://localhost:8080/api/products
```

### 4. Update Product
```bash
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptop",
    "price": 18000000,
    "stock": 8,
    "category_id": 1
  }'
```

### 5. Delete Product
```bash
curl -X DELETE http://localhost:8080/api/products/1
```

---

## 🐛 Troubleshooting

### Database Connection Error
```
Failed to initialize database: ...
```
- Verify PostgreSQL is running
- Check `DB_CONN` in `.env` file
- Ensure database exists

### Port Already in Use
```
address already in use
```
- Change `PORT` in `.env`
- Or kill the process using port 8080

### Validation Errors
Check the error message returned and verify:
- All required fields are provided
- String lengths meet minimum and maximum requirements
- Numbers are non-negative where required
- Foreign keys reference existing records

---

## 📞 Support

For issues or questions, please check:
1. The error message returned by the API
2. Database connection settings
3. Input validation rules
4. PostgreSQL logs

---

## 📄 License

This project is part of the Golang Training Course by Aria.

---

## ✍️ Author

Muflikhan Dwiprayogi

---

**Last Updated:** January 2026
