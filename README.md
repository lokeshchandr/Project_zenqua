

## Architecture

```
Client → API Gateway (port 8000) → Backend Services
                                    ├─ Auth Service (8001)
                                    ├─ Product Service (8002)
                                    └─ Order Service (8003)
```

## Installation

```powershell
cd d:\projects\Ecommerce\ApiGateway

# Install dependencies
go mod tidy

```

## Running the Gateway

```powershell
# Development mode
go run cmd/main.go

# Or build and run
go build -o api-gateway.exe cmd/main.go
.\api-gateway.exe
```

## API Routes

### Public Routes (No Authentication)

| Method | Path | Backend | Description |
|--------|------|---------|-------------|
| POST | `/auth/register` | Auth | Register new user |
| POST | `/auth/login` | Auth | User login |
| GET | `/products` | Product | List all products |
| GET | `/products/:id` | Product | Get single product |
| GET | `/health` | Gateway | Health check |

### Protected Routes (Authentication Required)

| Method | Path | Backend | Role | Description |
|--------|------|---------|------|-------------|
| GET | `/notifications` | Auth | Any | Get user notifications |
| GET | `/orders` | Order | Any | Get user orders |
| POST | `/orders` | Order | Any | Create new order |
| GET | `/orders/:id` | Order | Any | Get specific order |
| GET | `/products/my-products` | Product | Seller | Get seller's products |

### Admin Only Routes

| Method | Path | Backend | Role | Description |
|--------|------|---------|------|-------------|
| PUT | `/admin/approve/:id` | Auth | Admin | Approve admin user |
| POST | `/products` | Product | Admin | Create product |
| PATCH | `/products/:id` | Product | Admin | Update product |
| PATCH | `/products/:id/stock` | Product | Admin | Update stock |
| DELETE | `/products/:id` | Product | Admin | Delete product |
| PATCH | `/orders/:id/status` | Order | Admin | Update order status |


### Super Admin password
```
Email: root@root.com
password: root123
```
### Postman Collection
```
https://.postman.co/workspace/My-Workspace~27f72f6c-384a-44f6-8ab5-df9444976f7a/collection/39904644-d6d7f568-8616-4c4e-9d43-0d3d77424152?action=share&creator=39904644
```
