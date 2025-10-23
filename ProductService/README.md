# Product Management Microservice

A clean architecture Go microservice for managing products using Gin framework and SQLite.

## Features

- **Clean Architecture**: Separated layers (handlers, service, repository, models)
- **JWT Authentication**: Role-based access control (Admin/Super Admin)
- **SQLite Database**: Lightweight, embedded database
- **RESTful API**: Standard HTTP methods and JSON responses
- **Public & Protected Routes**: Public product viewing, Admin-only modifications

## Architecture

```
ProductService/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── models/
│   │   └── product.go       # Product model and DTOs
│   ├── repo/
│   │   └── product_repo.go  # Data access layer
│   ├── service/
│   │   └── product_service.go # Business logic layer
│   ├── handlers/
│   │   └── product_handler.go # HTTP handlers
│   ├── middleware/
│   │   └── auth.go          # JWT authentication
│   └── db/
│       └── db.go            # Database initialization
├── Dockerfile
├── go.mod
└── README.md
```

## Environment Variables

```bash
JWT_SECRET=your-jwt-secret-key  # Required: Same secret as AuthService
PORT=8002                        # Optional: Default 8002
```

## API Endpoints

### Public Routes (No Authentication)

#### List All Products
```http
GET /products
```

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Product Name",
      "description": "Product description",
      "price": 99.99,
      "quantity": 50,
      "created_at": "2025-10-18T12:00:00Z",
      "updated_at": "2025-10-18T12:00:00Z"
    }
  ]
}
```

#### Get Single Product
```http
GET /products/:id
```

**Response:**
```json
{
  "product": {
    "id": 1,
    "name": "Product Name",
    "description": "Product description",
    "price": 99.99,
    "quantity": 50,
    "created_at": "2025-10-18T12:00:00Z",
    "updated_at": "2025-10-18T12:00:00Z"
  }
}
```

### Protected Routes (Admin/Super Admin Only)

**Authorization Header Required:**
```
Authorization: Bearer <jwt-token>
```

#### Create Product
```http
POST /products
Content-Type: application/json

{
  "name": "New Product",
  "description": "Product description",
  "price": 99.99,
  "quantity": 100
}
```

**Response:**
```json
{
  "message": "product created successfully",
  "product": { ... }
}
```

#### Update Product
```http
PATCH /products/:id
Content-Type: application/json

{
  "name": "Updated Name",
  "price": 89.99,
  "quantity": 75
}
```

**Response:**
```json
{
  "message": "product updated successfully",
  "product": { ... }
}
```

#### Update Stock Only
```http
PATCH /products/:id/stock
Content-Type: application/json

{
  "quantity": 200
}
```

**Response:**
```json
{
  "message": "stock updated successfully",
  "product": { ... }
}
```

#### Delete Product
```http
DELETE /products/:id
```

**Response:**
```json
{
  "message": "product deleted successfully"
}
```

## Quick Start (Windows PowerShell)

### 1. Set Environment Variables

```powershell
$env:JWT_SECRET = "change-this-secret"  # Use same secret as AuthService
$env:PORT = "8002"
```

### 2. Install Dependencies

```powershell
cd d:\projects\Ecommerce\ProductService
go mod tidy
```

### 3. Run the Service

```powershell
go run cmd/main.go
```

Or build and run:

```powershell
go build -o product-service.exe cmd/main.go
.\product-service.exe
```

## Testing with Postman

### 1. Get JWT Token from AuthService

First, login to AuthService to get a token:

```http
POST http://localhost:8001/login
Content-Type: application/json

{
  "email": "root@root.com",
  "password": "root123"
}
```

Copy the `token` from response.

### 2. Create a Product (Admin Only)

```http
POST http://localhost:8002/products
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 1299.99,
  "quantity": 50
}
```

### 3. List All Products (Public)

```http
GET http://localhost:8002/products
```

### 4. Update Stock (Admin Only)

```http
PATCH http://localhost:8002/products/1/stock
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "quantity": 75
}
```

## Role-Based Access

- **Public Access**: `GET /products`, `GET /products/:id`
- **Admin/Super Admin Only**: All other endpoints (CREATE, UPDATE, DELETE)
- JWT token must contain `role` claim with value `"saler"` or `"superadmin"`

## Database Schema

**Table: products**

| Column      | Type      | Constraints                |
|-------------|-----------|----------------------------|
| id          | INTEGER   | PRIMARY KEY, AUTOINCREMENT |
| name        | TEXT      | NOT NULL                   |
| description | TEXT      |                            |
| price       | REAL      | NOT NULL                   |
| quantity    | INTEGER   | NOT NULL, DEFAULT 0        |
| created_at  | TIMESTAMP | AUTO                       |
| updated_at  | TIMESTAMP | AUTO                       |

## Docker

### Build Image

```powershell
docker build -t product-service .
```

### Run Container

```powershell
docker run -d -p 8002:8002 `
  -e JWT_SECRET=your-secret `
  -v ${PWD}/product.db:/root/product.db `
  --name product-service `
  product-service
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "invalid product id"
}
```

### 401 Unauthorized
```json
{
  "error": "missing authorization header"
}
```

### 403 Forbidden
```json
{
  "error": "admin access required"
}
```

### 404 Not Found
```json
{
  "error": "product not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "failed to create product"
}
```

## Development

### Project Structure Principles

- **cmd/**: Application entry points
- **internal/**: Private application code
  - **models/**: Data structures and DTOs
  - **repo/**: Database access layer (interfaces + implementations)
  - **service/**: Business logic layer
  - **handlers/**: HTTP request handlers
  - **middleware/**: HTTP middleware (auth, logging, etc.)
  - **db/**: Database initialization and migrations

### Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM library
- **glebarez/sqlite**: Pure Go SQLite driver
- **golang-jwt/jwt**: JWT implementation

## Notes

- The service uses the same JWT_SECRET as AuthService for token validation
- SQLite database file `product.db` is created automatically on first run
- All timestamps are stored in UTC
- Stock updates use the dedicated `/stock` endpoint to ensure clean separation
- Partial updates are supported via PATCH endpoints (only send fields to update)

## Next Steps

- Add pagination for product listing
- Add search and filtering capabilities
- Add product categories
- Add product images support
- Add audit logging
- Add rate limiting
- Add caching layer (Redis)
