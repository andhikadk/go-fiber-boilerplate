# Go Fiber Boilerplate

A production-ready boilerplate for building REST APIs with **Fiber**, a fast and lightweight Go web framework inspired by Express.js.

## 🚀 Features

- **Fiber Web Framework** - Fast, minimalist web framework
- **JWT Authentication** - Secure token-based authentication with flexible header support
- **Swagger/OpenAPI Docs** - Auto-generated interactive API documentation
- **DTO Layer** - Self-validating Data Transfer Objects with clean separation
- **GORM ORM** - Database abstraction layer with explicit dependency injection
- **PostgreSQL & SQLite** - Multiple database support
- **Testing Infrastructure** - Test utilities, fixtures, and custom assertions
- **Structured Logging** - File and console logging with rotation
- **Concurrent Programming Examples** - 7 production-ready concurrency patterns
- **Middleware Stack** - CORS, Logger, Recovery, Helmet, Auth
- **Request Validation** - Self-validating DTOs with detailed error messages
- **Error Handling** - Centralized error management
- **Database Migrations** - Dual-system (AutoMigrate for dev, SQL for prod)
- **Unit Tests** - Testing setup with in-memory SQLite
- **Docker Support** - Containerized deployment with hot reload option
- **Hot Reload** - Development mode with Air
- **Environment Management** - .env configuration with validation

## 📋 Project Structure

```
go-fiber-boilerplate/
├── cmd/
│   └── api/
│       └── main.go                # Application entry point
├── assets/
│   ├── embed.go                   # Embedded migrations
│   └── migrations/                # SQL migration files
│       ├── 001_initial_schema.sql
│       ├── 002_add_indexes.sql
│       └── seeds/                 # Seed data files
├── config/
│   ├── config.go                  # Configuration management
│   └── database.go                # Database connection setup
├── docs/                          # Swagger/OpenAPI docs (auto-generated)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── dto/                       # Data Transfer Objects (NEW)
│   │   ├── auth_dto.go            # Auth DTOs with validation
│   │   ├── user_dto.go            # User DTOs
│   │   └── book_dto.go            # Book DTOs
│   ├── handlers/                  # HTTP handlers (with Swagger annotations)
│   │   ├── auth.go                # Authentication endpoints
│   │   ├── books.go               # Books CRUD endpoints
│   │   ├── concurrent.go.bak      # Concurrent patterns (disabled)
│   │   ├── health.go              # Health checks
│   │   └── user.go                # User management endpoints
│   ├── services/                  # Business logic (explicit DI)
│   │   ├── auth_service.go        # Auth business logic
│   │   ├── book_service.go        # Books business logic
│   │   └── concurrent_service.go  # Concurrent patterns
│   ├── models/                    # Domain models (database entities)
│   │   ├── user.go                # User model
│   │   ├── book.go                # Book model
│   │   └── response.go            # API response models
│   ├── middleware/                # Custom middlewares
│   │   ├── auth.go                # JWT authentication
│   │   └── error.go               # Error handling
│   ├── database/                  # Database layer
│   │   ├── db.go                  # Connection & migrations
│   │   ├── migrator.go            # SQL migration runner
│   │   └── seeder.go              # Database seeder
│   ├── routes/                    # Route definitions
│   │   └── routes.go              # All routes with Swagger endpoint
│   ├── testutil/                  # Testing utilities (NEW)
│   │   ├── db.go                  # Test database setup
│   │   ├── fixtures.go            # Test data fixtures
│   │   └── assert.go              # Custom assertions
│   └── utils/                     # Internal utilities (NEW)
│       └── logger.go              # Structured logging
├── pkg/                           # Reusable utilities
│   ├── utils/                     # Response formatting
│   │   └── response.go
│   └── jwt/                       # JWT token management
│       └── jwt.go
├── logs/                          # Application log files
│   └── app.log
├── .env.example                   # Environment template
├── go.mod & go.sum                # Dependencies
├── Dockerfile                     # Multi-stage Docker build (production)
├── Dockerfile.dev                 # Development Docker with hot reload
├── docker-compose.yml             # Production Docker Compose
├── docker-compose.dev.yml         # Development Docker Compose
├── .air.toml                      # Air hot reload configuration
├── .dockerignore                  # Docker build ignore rules
├── Makefile                       # Build and development commands
├── README.md                      # This file
└── CLAUDE.md                      # Architecture & Claude Code instructions
```

## 🛠️ Tech Stack

- **Framework:** Fiber v2
- **Documentation:** Swagger/OpenAPI (swaggo/swag)
- **Database:** GORM, PostgreSQL, SQLite
- **Authentication:** JWT (golang-jwt) with flexible header support
- **Security:** bcrypt (golang.org/x/crypto)
- **Validation:** Self-validating DTOs + go-playground/validator
- **Testing:** In-memory SQLite, testutil package, standard library
- **Logging:** Structured logging (file + console)
- **Environment:** godotenv with validation
- **Middleware:** Fiber built-in + custom (Auth, Error, CORS, etc.)

## 📦 Dependencies

### Core Dependencies
```
github.com/gofiber/fiber/v2 v2.52.5
github.com/gofiber/contrib/jwt v1.0.10
github.com/golang-jwt/jwt/v5 v5.2.1
gorm.io/gorm v1.31.0
gorm.io/driver/postgres v1.6.0
gorm.io/driver/sqlite v1.6.0
golang.org/x/crypto v0.43.0
github.com/go-playground/validator/v10 v10.28.0
github.com/joho/godotenv v1.5.1
```

## ⚡ Quick Start

### Prerequisites (Choose ONE based on your preferred setup)

**Option A: Docker + Hot Reload (Recommended)** ⭐
- Docker Desktop or Docker Engine
- Docker Compose v1.29+
- (No need to install Go, PostgreSQL, or Make locally!)

**Option B: Production-like (Docker)**
- Docker Desktop or Docker Engine
- Docker Compose v1.29+
- (No need to install Go, PostgreSQL, or Make locally!)

**Option C: Local Development (No Docker)**
- Go 1.25 or higher
- PostgreSQL 12+ (or SQLite)
- Make (optional, for using Makefile)
- git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/go-fiber-boilerplate.git
cd go-fiber-boilerplate
```

2. **Install dependencies**
```bash
make install-deps
# or
go mod download && go mod tidy
```

3. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Setup database and run application**

**Option A: Using Docker Compose with Hot Reload (Recommended for Development)** ⭐
```bash
make docker-dev
```
This will:
- Start PostgreSQL database (accessible at `localhost:6543`)
- Run application with hot reload (instant reload on code changes)
- Run migrations automatically
- Start the API on `http://localhost:4000`

> 💡 **Tip:** This is perfect for development! Edit your code and see changes instantly without rebuilding.

View logs:
```bash
make docker-dev-logs
```

**Option B: Using Docker Compose (Production-like)**
```bash
make docker-up
```
This will:
- Start PostgreSQL database (accessible at `localhost:6543`)
- Run migrations automatically
- Start the API on `http://localhost:4000`

> ⚠️ **Important:** The application is already running inside Docker! You do NOT need to run `make run` after this. The API is ready at `http://localhost:4000`

**Option C: Local Development (without Docker)**
```bash
# Create PostgreSQL database manually
createdb fiber_boilerplate

# Run migrations
make migrate

# Run application
make run
```
The API will be available at `http://localhost:4000`

## 🎯 For New Developers - Getting Started

Welcome! This guide will help you customize this boilerplate for your own project.

### Step 1: Initial Setup (5 minutes)

1. **Clone and setup**
   ```bash
   git clone <your-repo-url>
   cd go-fiber-boilerplate
   make install-deps        # Install Go dependencies
   cp .env.example .env     # Create environment file
   ```

2. **Start development environment**
   ```bash
   make docker-dev          # PostgreSQL + API with hot reload
   ```

3. **Verify it works**
   - API Health: http://localhost:4000/health
   - Swagger Docs: http://localhost:4000/swagger/
   - Expected response: `{"status":"ok"}`

### Step 2: Customize for Your Project (10 minutes)

#### 2.1 Change Module Name

```bash
# 1. Update go.mod
go mod edit -module github.com/yourname/your-project

# 2. Update all imports (find & replace in your IDE)
# From: go-fiber-boilerplate
# To: github.com/yourname/your-project

# 3. Tidy dependencies
go mod tidy
```

#### 2.2 Update Project Info

- **Swagger Info**: Edit `cmd/api/main.go` lines 26-41
  ```go
  //	@title       Your Project Name
  //	@description Your project description
  //	@contact.name    Your Name
  //	@contact.email   your@email.com
  //	@host localhost:4000
  ```

- **Environment**: Edit `.env`
  - Change `JWT_SECRET` to a strong random value
  - Update `DB_*` credentials if needed
  - Set `APP_NAME` to your project name

- **README**: Update this file with your project name and description

- **Regenerate Swagger**: `make swagger`

#### 2.3 Remove Example Features (Optional)

If you don't need the Books CRUD example:
```bash
rm internal/handlers/books.go
rm internal/services/book_service.go
rm internal/dto/book_dto.go
# Then remove Book model from internal/models/book.go
# And remove book routes from internal/routes/routes.go
make swagger              # Regenerate docs
```

### Step 3: Add Your First Feature (15 minutes)

**Example: Add "Products" CRUD**

#### 3.1 Create Model

```go
// internal/models/product.go
package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"default:0" json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
```

#### 3.2 Create DTOs with Validation

```go
// internal/dto/product_dto.go
package dto

import (
	"errors"
	"strings"
)

type CreateProductRequest struct {
	Name        string  `json:"name" example:"Laptop"`
	Description string  `json:"description" example:"High-performance laptop"`
	Price       float64 `json:"price" example:"999.99"`
	Stock       int     `json:"stock" example:"50"`
}

func (r *CreateProductRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if r.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if r.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	return nil
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" example:"Updated Laptop"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" example:"899.99"`
	Stock       *int     `json:"stock" example:"45"`
}

func (r *UpdateProductRequest) Validate() error {
	if r.Price != nil && *r.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if r.Stock != nil && *r.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	return nil
}
```

#### 3.3 Create Service

```go
// internal/services/product_service.go
package services

import (
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) Create(req *dto.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := s.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) Update(id uint, req *dto.UpdateProductRequest) (*models.Product, error) {
	product, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := s.db.Save(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}
```

#### 3.4 Create Handler with Swagger Annotations

```go
// internal/handlers/products.go
package handlers

import (
	"strconv"
	"go-fiber-boilerplate/internal/database"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/services"
	"go-fiber-boilerplate/internal/utils"
	pkgUtils "go-fiber-boilerplate/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// CreateProduct godoc
// @Summary      Create product
// @Description  Create a new product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProductRequest true "Product data"
// @Success      201 {object} models.APIResponse{data=models.Product}
// @Failure      400 {object} models.APIResponse
// @Failure      500 {object} models.APIResponse
// @Router       /api/products [post]
// @Security     BearerAuth
func CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		utils.ErrorLogger.Printf("[CreateProduct] Failed to parse: %v", err)
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	if err := req.Validate(); err != nil {
		utils.ErrorLogger.Printf("[CreateProduct] Validation failed: %v", err)
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	service := services.NewProductService(database.GetDB())
	product, err := service.Create(&req)
	if err != nil {
		utils.ErrorLogger.Printf("[CreateProduct] Failed to create: %v", err)
		return pkgUtils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	utils.InfoLogger.Printf("[CreateProduct] Product created: ID=%d", product.ID)
	return pkgUtils.CreatedResponse(c, "Product created successfully", product)
}

// GetProducts godoc
// @Summary      Get all products
// @Description  Retrieve list of all products
// @Tags         Products
// @Produce      json
// @Success      200 {object} models.APIResponse{data=[]models.Product}
// @Failure      500 {object} models.APIResponse
// @Router       /api/products [get]
// @Security     BearerAuth
func GetProducts(c *fiber.Ctx) error {
	service := services.NewProductService(database.GetDB())
	products, err := service.GetAll()
	if err != nil {
		utils.ErrorLogger.Printf("[GetProducts] Failed to fetch: %v", err)
		return pkgUtils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch products")
	}

	return pkgUtils.SuccessResponse(c, "Products retrieved successfully", products)
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Retrieve a specific product
// @Tags         Products
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} models.APIResponse{data=models.Product}
// @Failure      400 {object} models.APIResponse
// @Failure      404 {object} models.APIResponse
// @Router       /api/products/{id} [get]
// @Security     BearerAuth
func GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return pkgUtils.BadRequestResponse(c, "Invalid product ID")
	}

	service := services.NewProductService(database.GetDB())
	product, err := service.GetByID(uint(id))
	if err != nil {
		return pkgUtils.NotFoundResponse(c, "Product not found")
	}

	return pkgUtils.SuccessResponse(c, "Product retrieved successfully", product)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update an existing product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Param        request body dto.UpdateProductRequest true "Product update data"
// @Success      200 {object} models.APIResponse{data=models.Product}
// @Failure      400 {object} models.APIResponse
// @Failure      404 {object} models.APIResponse
// @Router       /api/products/{id} [put]
// @Security     BearerAuth
func UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return pkgUtils.BadRequestResponse(c, "Invalid product ID")
	}

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgUtils.BadRequestResponse(c, "Invalid request body")
	}

	if err := req.Validate(); err != nil {
		return pkgUtils.BadRequestResponse(c, err.Error())
	}

	service := services.NewProductService(database.GetDB())
	product, err := service.Update(uint(id), &req)
	if err != nil {
		return pkgUtils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update product")
	}

	return pkgUtils.SuccessResponse(c, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Delete a product
// @Tags         Products
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} models.APIResponse
// @Failure      400 {object} models.APIResponse
// @Router       /api/products/{id} [delete]
// @Security     BearerAuth
func DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return pkgUtils.BadRequestResponse(c, "Invalid product ID")
	}

	service := services.NewProductService(database.GetDB())
	if err := service.Delete(uint(id)); err != nil {
		return pkgUtils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete product")
	}

	return pkgUtils.SuccessResponse(c, "Product deleted successfully", nil)
}
```

#### 3.5 Add Routes

```go
// internal/routes/routes.go
// Add inside apiGroup

productsGroup := apiGroup.Group("/products")
{
	productsGroup.Post("/", handlers.CreateProduct)
	productsGroup.Get("/", handlers.GetProducts)
	productsGroup.Get("/:id", handlers.GetProduct)
	productsGroup.Put("/:id", handlers.UpdateProduct)
	productsGroup.Delete("/:id", handlers.DeleteProduct)
}
```

#### 3.6 Add to Database Migration

```go
// internal/database/db.go
// Add Product to AutoMigrate
if err := db.AutoMigrate(
	&models.User{},
	&models.Book{},
	&models.Product{}, // Add this line
); err != nil {
```

#### 3.7 Generate Swagger & Test

```bash
make swagger              # Generate Swagger docs
make migrate              # Run migrations (AutoMigrate)
make dev                  # Start with hot reload

# Visit: http://localhost:4000/swagger/
# Test the new /api/products endpoints
```

### Step 4: Common Development Tasks

#### Add Environment Variable

1. Add to `.env.example`:
   ```
   STRIPE_API_KEY=sk_test_xxxxx
   ```

2. Add to `config/config.go`:
   ```go
   type Config struct {
       // ... existing fields
       StripeAPIKey string `mapstructure:"STRIPE_API_KEY"`
   }
   ```

3. Use in code:
   ```go
   apiKey := config.AppConfig.StripeAPIKey
   ```

#### Add Database Migration (SQL)

```bash
# Create new migration file
touch assets/migrations/003_add_products.sql

# Write migration
-- 003_add_products.sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);

# Run migration
make migrate-sql
```

#### Add Middleware

```go
// internal/middleware/rate_limit.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

func RateLimitMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP
		},
	})
}
```

Apply in routes:
```go
apiGroup.Use(middleware.RateLimitMiddleware())
```

#### Write Tests

```go
// internal/services/product_service_test.go
package services

import (
	"testing"
	"go-fiber-boilerplate/internal/dto"
	"go-fiber-boilerplate/internal/models"
	"go-fiber-boilerplate/internal/testutil"
)

func TestCreateProduct(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	service := NewProductService(db)

	// Create test data
	req := &dto.CreateProductRequest{
		Name:  "Test Product",
		Price: 99.99,
		Stock: 10,
	}

	// Execute
	product, err := service.Create(req)

	// Assert
	testutil.AssertNoError(t, err, "Should create product without error")
	testutil.AssertNotNil(t, product, "Product should not be nil")
	testutil.AssertEqual(t, "Test Product", product.Name, "Product name should match")
	testutil.AssertEqual(t, 99.99, product.Price, "Product price should match")
}

func TestGetAllProducts(t *testing.T) {
	db := testutil.SetupTestDB(t)
	service := NewProductService(db)

	// Create test products
	db.Create(&models.Product{Name: "Product 1", Price: 10.0, Stock: 5})
	db.Create(&models.Product{Name: "Product 2", Price: 20.0, Stock: 3})

	// Execute
	products, err := service.GetAll()

	// Assert
	testutil.AssertNoError(t, err)
	testutil.AssertLen(t, products, 2, "Should return 2 products")
}
```

Run tests:
```bash
make test
```

### Step 5: Deploy to Production

#### Production Checklist

- [ ] Change `JWT_SECRET` to cryptographically secure random value
- [ ] Set `ENV=production` in `.env`
- [ ] Update database credentials (use strong passwords)
- [ ] Configure `CORS_ALLOWED_ORIGINS` to your frontend domain
- [ ] Review and enable security middleware (helmet, rate limiting)
- [ ] Setup SSL/TLS certificates (Let's Encrypt, Cloudflare, etc.)
- [ ] Configure proper logging (set `LOG_LEVEL=error` for production)
- [ ] Setup health check monitoring
- [ ] Configure database backups
- [ ] Test with `make docker-up` (production mode)
- [ ] Setup CI/CD pipeline (GitHub Actions, GitLab CI, etc.)

#### Build & Deploy

**Option 1: Docker (Recommended)**
```bash
# Build production image
docker compose build --no-cache

# Start production containers
docker compose up -d

# Check logs
docker compose logs -f fiber_app

# Check health
curl http://your-domain.com/health
```

**Option 2: Binary Deployment**
```bash
# Build binary
make build

# Deploy binary to server
scp ./bin/go-fiber-boilerplate user@server:/opt/myapp/

# Run on server
./go-fiber-boilerplate
```

**Option 3: Kubernetes**
```bash
# Create Docker image
docker build -t your-registry/app:v1.0.0 .

# Push to registry
docker push your-registry/app:v1.0.0

# Deploy to K8s
kubectl apply -f k8s/deployment.yaml
```

### Useful Commands Cheatsheet

```bash
# Development
make dev                 # Hot reload (local)
make docker-dev         # Docker with hot reload
make run                # Run without hot reload

# Database
make migrate            # AutoMigrate (dev)
make migrate-sql        # SQL migrations (prod)
make migrate-status     # Show migration status
make seed               # Seed sample data

# Code Quality
make test               # Run all tests
make test-coverage      # Coverage report (opens browser)
make lint               # Lint code (requires golangci-lint)
make fmt                # Format code
make vet                # Run go vet

# Swagger
make swagger            # Generate API docs
make swagger-fmt        # Format Swagger comments

# Docker
make docker-dev         # Dev with hot reload
make docker-dev-logs    # View dev logs
make docker-dev-down    # Stop dev containers
make docker-dev-reset   # Reset dev (removes data!)
make docker-up          # Production mode
make docker-logs        # View prod logs
make docker-down        # Stop prod containers
make docker-reset       # Reset prod (removes data!)

# Build
make build              # Build binary
make clean              # Clean build artifacts
make install-deps       # Download dependencies
```

### Learning Resources

- **Architecture Details**: See `CLAUDE.md` for in-depth patterns and best practices
- **Swagger Docs**: http://localhost:4000/swagger/ (auto-generated from code)
- **Concurrency Patterns**: Check `internal/services/concurrent_service.go` for 7 production-ready patterns
- **Testing Examples**: See `internal/testutil/` for test utilities and patterns

### Need Help?

- Check `CLAUDE.md` for detailed architecture explanations
- Review existing handlers in `internal/handlers/` for patterns
- Look at `internal/services/` for business logic examples
- Explore `internal/dto/` for DTO validation examples

## 🚀 Usage

### Run Application
```bash
make run
# or
go run cmd/api/main.go
```

### Build
```bash
make build
```
Binary will be created at `./bin/go-fiber-boilerplate`

### Development with Hot Reload
```bash
make dev
# Requires: go install github.com/cosmtrek/air@latest
```

### Testing
```bash
make test                # Run all tests
make test-coverage       # Run tests with coverage report
```

### Database
```bash
make migrate            # Run migrations
make seed               # Seed sample data
```

### Code Quality
```bash
make fmt                # Format code
make lint               # Run linter
make vet                # Run go vet
```

### Docker
```bash
make docker-build       # Build Docker image
make docker-up          # Start containers
make docker-down        # Stop containers
make docker-logs        # View logs
```

## 🔐 Authentication

This boilerplate uses JWT (JSON Web Tokens) for authentication:

1. **Register** - POST `/auth/register`
2. **Login** - POST `/auth/login` (returns JWT token)
3. **Protected Routes** - Add `Authorization: Bearer <token>` header

Tokens expire after 15 minutes by default. Adjust in `.env` with `JWT_EXPIRY`.

## 📚 API Endpoints (Example)

### Health Check
```
GET /health
```

### Authentication
```
POST /auth/register      # Register new user
POST /auth/login         # Login and get JWT token
POST /auth/refresh       # Refresh JWT token
```

### Books (Protected)
```
GET    /api/books        # List all books (requires auth)
GET    /api/books/:id    # Get book by ID (requires auth)
POST   /api/books        # Create book (requires auth)
PUT    /api/books/:id    # Update book (requires auth)
DELETE /api/books/:id    # Delete book (requires auth)
```

### Concurrent Patterns (Protected) 🆕
```
GET    /api/concurrent                    # Overview of all patterns
GET    /api/concurrent/parallel           # Parallel processing with goroutines
GET    /api/concurrent/worker-pool        # Worker pool pattern
GET    /api/concurrent/fan-out-fan-in     # Fan-out/fan-in pattern
GET    /api/concurrent/pipeline           # Pipeline pattern
POST   /api/concurrent/bulk-create        # Semaphore (rate limiting)
GET    /api/concurrent/timeout/:id        # Timeout pattern
GET    /api/concurrent/monitor/:id        # Select with multiple channels
```

## ⚡ Concurrent Programming Patterns

This boilerplate includes **7 production-ready concurrent programming patterns** to help developers understand and implement Go's concurrency features.

### 🎯 Why Learn Concurrency?

Go's concurrency model (goroutines and channels) is one of its most powerful features. Understanding these patterns will help you:

- **Build faster applications** - Process multiple tasks simultaneously
- **Handle high traffic** - Scale your API to serve thousands of requests
- **Implement background jobs** - Run tasks asynchronously without blocking
- **Control resource usage** - Prevent overwhelming your database or external APIs

### 📚 Available Patterns

| Pattern | Use Case | Example |
|---------|----------|---------|
| **Basic Goroutines + WaitGroup** | Parallel data fetching | Fetch multiple books simultaneously |
| **Worker Pool** | Rate limiting, job queues | Process tasks with limited workers |
| **Fan-Out/Fan-In** | Multi-source aggregation | Search across multiple fields in parallel |
| **Pipeline** | Multi-stage processing | ETL operations with stages |
| **Semaphore (Rate Limiting)** | API rate limiting | Limit concurrent database writes |
| **Timeout** | External API calls | Cancel slow operations |
| **Select with Multiple Channels** | Event handling | Monitor changes in real-time |

### 🚀 Quick Start

1. **Start the application:**
   ```bash
   make docker-dev
   ```

2. **Login to get JWT token:**
   ```bash
   curl -X POST http://localhost:4000/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@example.com","password":"admin123"}'
   ```

3. **View all available patterns:**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:4000/api/concurrent
   ```

4. **Test a pattern (Worker Pool example):**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://localhost:4000/api/concurrent/worker-pool?ids=1,2,3,4,5&workers=3"
   ```

### 📖 Pattern Details

Each pattern includes:
- **Production-ready code** with comprehensive error handling
- **Context-based cancellation** for graceful shutdown
- **Detailed comments** explaining each step
- **Real-world use cases** demonstrated via API endpoints

All patterns are fully functional and can be tested immediately via the API endpoints above.

### 💡 Key Concepts

**Goroutines** - Lightweight threads managed by Go runtime
```go
go func() {
    // This runs concurrently
}()
```

**Channels** - Communication between goroutines
```go
ch := make(chan int)
ch <- 42        // Send
value := <-ch   // Receive
```

**Select** - Handle multiple channels
```go
select {
case msg := <-ch1:
    // Handle ch1
case msg := <-ch2:
    // Handle ch2
case <-time.After(5*time.Second):
    // Timeout
}
```

**Context** - Cancellation and timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### 🎓 Learning Path

1. **Start with Pattern 1** (Basic Goroutines) - Understand fundamentals
2. **Try Worker Pool** - Learn resource control
3. **Explore Fan-Out/Fan-In** - Master result aggregation
4. **Practice with examples** - Test all 7 patterns via API
5. **Read source code** - Study implementation details
6. **Apply to your project** - Use patterns in real scenarios

### 📁 Source Files

- **Service Layer:** `internal/services/concurrent_service.go` - All pattern implementations
- **Handler Layer:** `internal/handlers/concurrent.go` - API endpoints for each pattern
- **Routes:** `internal/routes/routes.go` - Route definitions

## 📝 Configuration

All configuration is managed through `.env` file. See `.env.example` for all available options.

### Key Configuration
- `PORT` - Server port (default: 4000)
- `ENV` - Environment (development/production)
- `DB_HOST` - Database host (localhost for local, postgres for Docker)
- `DB_PORT` - Database port (5432 for local/Docker internal, 6543 for Docker host access)
- `DB_DRIVER` - Database driver (postgres/sqlite)
- `JWT_SECRET` - Secret key for JWT signing
- `CORS_ALLOWED_ORIGINS` - Allowed origins for CORS
- `LOG_LEVEL` - Logging level (info/debug/error)

## 🧪 Testing

Run all tests:
```bash
go test -v ./...
```

Run specific test:
```bash
go test -v ./tests -run TestName
```

With coverage:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🐳 Docker Deployment

### Development Setup with Hot Reload (Recommended)
```bash
make docker-dev
# or
docker-compose -f docker-compose.dev.yml up -d
```

This will:
- Start PostgreSQL database (port 6543 from host, 5432 internal)
- Run application with hot reload enabled (air watches for code changes)
- Run database migrations automatically
- Expose API on port 4000 (`http://localhost:4000`)

> 💡 **Benefits:** Changes to Go files are automatically detected and compiled. Just save your file and refresh the browser - no rebuild needed!

### Production Setup (No Hot Reload)
```bash
make docker-up
# or
docker-compose up -d
```

This will:
- Build the Fiber application into a compiled binary
- Start PostgreSQL database (port 6543 from host, 5432 internal)
- Run database migrations automatically
- Expose API on port 4000 (`http://localhost:4000`)

### View Logs
```bash
# For development setup
make docker-dev-logs
# or
docker-compose -f docker-compose.dev.yml logs -f

# For production setup
make docker-logs
# or
docker-compose logs -f

# View specific service logs
docker-compose logs -f fiber_app    # Just app logs
docker-compose logs -f postgres     # Just database logs
```

### Check Container Status
```bash
docker-compose ps
```

Expected output (both containers should be "Up"):
```
NAME                    STATUS
fiber_boilerplate_app   Up (healthy)
fiber_boilerplate_db    Up (healthy)
```

### Verify Application is Running
```bash
# Check if app is responding
curl http://localhost:4000/health

# Expected response: {"status":"ok"}
```

### Stop

**Development setup:**
```bash
make docker-dev-down
# or
docker-compose -f docker-compose.dev.yml down
```

**Production setup:**
```bash
make docker-down
# or
docker-compose down
```

### Reset Database (Remove all data and volumes)

**Development setup:**
```bash
make docker-dev-reset
# This removes development containers, networks, AND database volumes
# ⚠️ Warning: All data will be deleted!
```

**Production setup:**
```bash
make docker-reset
# This removes containers, networks, AND database volumes
# ⚠️ Warning: All data will be deleted!
```

## 📖 Project Structure Details

### `main.go`
Application entry point. Initializes config, database, and starts the server.

### `embed.go`
Embedded file system for database migrations. Uses Go's `embed` package to bundle migration files into the binary.

### `Dockerfile`
Production-ready multi-stage Docker build. Compiles Go code into an optimized binary with minimal dependencies.

### `Dockerfile.dev`
Development Docker image with hot reload support using air. Includes Go compiler and air tool for watching code changes.

### `.air.toml`
Configuration file for air hot reload tool. Specifies which files to watch, build commands, and reload behavior.

### `docker-compose.yml`
Production Docker Compose configuration. Runs compiled binary in isolated containers with PostgreSQL.

### `docker-compose.dev.yml`
Development Docker Compose configuration. Runs application with air hot reload and mounts source code as volume for instant updates.

### `config/`
Configuration management and database setup.

### `internal/handlers/`
HTTP request handlers for different routes.

### `internal/models/`
Data structures for the application.

### `internal/services/`
Business logic layer.

### `internal/middleware/`
Custom middleware for authentication, error handling, etc.

### `internal/database/`
Database connection and initialization.

### `internal/routes/`
Route definitions and grouping.

### `pkg/utils/`
Utility functions (response formatting, validation, etc).

### `pkg/jwt/`
JWT token generation and validation.

## 🔄 Development Workflow

> ⚠️ **Important:** Choose ONE path below. Do NOT run multiple paths at the same time - they will conflict on port 4000.

### With Docker + Hot Reload ⭐ (Recommended)
```bash
make docker-dev         # Start database and app with hot reload
make docker-dev-logs    # View logs (in another terminal)
make test              # Run tests in another terminal
make docker-dev-down   # Stop containers when done
```

**Features:**
- ✅ Instant reload on code changes (no rebuild needed)
- ✅ Database in Docker (simple setup)
- ✅ Production-like environment
- 💡 Perfect for rapid development

**When code changes:**
1. Edit file (e.g., `internal/handlers/books.go`)
2. Save file
3. Air automatically detects changes and rebuilds
4. Refresh `http://localhost:4000` → see your changes instantly!

### With Docker (Production-like, No Hot Reload)
```bash
make docker-up         # Start database and production-like app
make docker-logs       # View logs (in another terminal)
make test             # Run tests in another terminal
make docker-down      # Stop containers when done
```

**Features:**
- ✅ Production-ready build
- ✅ Uses compiled binary
- ❌ No hot reload (need rebuild on changes)

### Local Development (without Docker)
```bash
make run               # Start application with compiled binary
make test              # Run tests in another terminal
make dev               # OR run with hot reload (requires air installed locally)
make migrate           # Run migrations
make seed              # Seed sample data
```

**Features:**
- ✅ Hot reload with local air (`make dev`)
- ❌ Need to setup PostgreSQL locally
- 💡 Lightest setup, full control

**Note:** Make sure PostgreSQL is running locally before `make run` or `make dev`.

### Development Steps
1. **Create models** in `internal/models/`
2. **Create handlers** in `internal/handlers/`
3. **Add business logic** in `internal/services/`
4. **Define routes** in `internal/routes/`
5. **Write tests** in `tests/`
6. **Run and test** with `make run` and `make test`

## 📚 Learning Resources

### Framework & Libraries
- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Guide](https://gorm.io/docs/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [Go Best Practices](https://golang.org/doc/effective_go)

### Concurrency (included in this boilerplate)
- **Source Code**: `internal/services/concurrent_service.go` - 7 production-ready patterns
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)
- [Go Concurrency Patterns (Video)](https://www.youtube.com/watch?v=f6kdp27TYZs)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the MIT License.

## 🙏 Acknowledgments

- Fiber team for the amazing framework
- GORM team for the powerful ORM
- Go community for best practices

---

**Happy coding! 🚀**

For issues and questions, please open an issue on GitHub.
