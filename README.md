# Inventory API - Modular Monolith

A modular monolith implementation of an Inventory Management System with POS and Merchant modules.

## Architecture

This project follows a **Modular Monolith** architecture pattern with clear module boundaries:

### Modules:
1. **Inventory Module** - Manages products and stock
2. **POS Module** - Handles sales transactions
3. **Merchant Module** - Manages stores and merchants

### Structure:
```
inventory-api/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Config
│   ├── database/         # Database
│   ├── server/           # Server and routes
│   └── modules/          # Business modules
│       ├── inventory/    # Product management
│       ├── pos/          # Point of Sale
│       └── merchant/     # Store management
```

### Layer Architecture (per module):
- **Domain**: Business entities and DTOs
- **Repository**: Data access layer
- **Service**: Business logic
- **Handler**: HTTP handlers
- **Routes**: Route registration

## Setup

1. Manage dependencies:
```bash
go mod tidy
```

2. Set up MongoDB:
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

3. Configure environment variables in ```.env```

4. Run the application:
```bash
go run cmd/api/main.go 
```
 or 
```bash
go run ./cmd/api .
```
## API Endpoints

### Merchant/Store Management
- `POST /api/v1/merchants/stores` - Create a new store
- `GET /api/v1/merchants/stores` - Get all stores
- `GET /api/v1/merchants/stores/:id` - Get store by ID

### Inventory Management
- `GET /api/v1/inventory/products` - List of products
- `POST /api/v1/inventory/products` - Create a new product
- `GET /api/v1/inventory/products/:id` - Get product by ID
- `GET /api/v1/inventory/stores/:storeId/products` - Get all products for a store

### POS (Point of Sale)
- `POST /api/v1/pos/sales` - Create a new sale transaction
- `GET /api/v1/pos/stores/:storeId/sales` - Get all sales for a store

## Example Requests

### Create Store
```bash
curl -X POST http://localhost:3000/api/v1/merchants/stores \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Main Store",
    "email": "main@store.com",
    "phone": "+1234567890",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zip_code": "10001",
      "country": "USA"
    }
  }'
```

### Create Product
```bash
curl -X POST http://localhost:3000/api/v1/inventory/products \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": "65f1234567890abcdef12345",
    "sku": "PROD-001",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1299.99,
    "cost": 899.99,
    "quantity": 50,
    "category": "Electronics"
  }'
```

### Create Sale
```bash
curl -X POST http://localhost:3000/api/v1/pos/sales \
  -H "Content-Type: application/json" \
  -d '{
    "store_id": "65f1234567890abcdef12345",
    "payment_type": "cash",
    "items": [
      {
        "product_id": "65f1234567890abcdef12346",
        "quantity": 2
      }
    ]
  }'
```

## Key Features

✅ **Modular Architecture**: Clear separation of concerns with independent modules
✅ **Domain-Driven Design**: Each module has its own domain, repository, service, and handler layers
✅ **MongoDB Integration**: Using official MongoDB Go driver
✅ **Clean Architecture**: Dependencies point inward (handler → service → repository)
✅ **Inter-module Communication**: POS module communicates with Inventory module through repositories
✅ **Graceful Shutdown**: Proper cleanup of resources
✅ **Stock Management**: Automatic inventory deduction on sales
✅ **Transaction Safety**: Sales validate stock before completing

## Module Communication

Modules can communicate through:
1. **Repository Layer**: Direct repository calls (e.g., POS accessing Product repository)
2. **Service Layer**: Service-to-service calls for complex operations
3. **Events**: (Can be added for async operations)

Example: When a sale is created, the POS module:
1. Validates products exist (via Inventory repository)
2. Checks stock availability
3. Creates the sale
4. Updates product quantities

## Future Updates

- Add unit tests
- Add authentication
- Add report module
- Add Event-driven communication between modules
- Implement CQRS pattern for complex queries
- Add caching (Redis or In-Memory Caching)