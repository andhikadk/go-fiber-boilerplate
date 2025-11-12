# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go Fiber Boilerplate is a production-ready REST API template using the Fiber framework (Express.js-like for Go), GORM ORM, and PostgreSQL/SQLite. It includes JWT authentication, database migrations, Docker support, and hot reload for development.

## Common Development Commands

### Building & Running
```bash
make build              # Build binary to ./bin/go-fiber-boilerplate
make run               # Run application immediately
make dev               # Run with hot reload (requires air installed)
```

### Testing
```bash
make test              # Run all tests with verbose output
make test-coverage     # Generate HTML coverage report (coverage.html)
```

### Code Quality
```bash
make fmt               # Format code with go fmt
make vet               # Run go vet analysis
make lint              # Run golangci-lint (install via: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
```

### Database
```bash
make migrate           # Run migrations (AutoMigrate in dev, SQL in prod)
make migrate-sql       # Run SQL migrations from embedded files
make migrate-status    # Show applied migrations and seeds
make seed              # Seed database with sample data
```

### Docker Development (Recommended)
```bash
make docker-dev        # Start with hot reload (best for development)
make docker-dev-logs   # View logs
make docker-dev-down   # Stop containers
make docker-dev-reset  # Reset containers and database
```

### Docker Production
```bash
make docker-up         # Start production containers
make docker-down       # Stop containers
make docker-logs       # View logs
make docker-reset      # Reset containers and database
```

### Swagger Documentation
```bash
make swagger-install   # Install swag CLI tool
make swagger           # Generate Swagger documentation
make swagger-fmt       # Format Swagger comments
```

After running the application, access interactive API documentation at:
- Swagger UI: `http://localhost:4000/swagger/index.html`

### Utilities
```bash
make install-deps      # Download and tidy dependencies
make clean             # Remove build artifacts
make help              # Show all available commands
make all               # Clean, install, build, and test
```

## Project Architecture

### High-Level Structure

The application follows a **layered architecture** with clear separation of concerns:

```
main.go (entry point)
  ├── config/ (configuration management)
  ├── docs/ (Swagger/OpenAPI documentation - auto-generated)
  ├── internal/ (core application logic)
  │   ├── dto/ (Data Transfer Objects - API contracts)
  │   ├── handlers/ (HTTP request handlers)
  │   ├── services/ (business logic layer)
  │   ├── models/ (domain models - database entities)
  │   ├── middleware/ (request interceptors)
  │   ├── database/ (DB initialization & migrations)
  │   ├── routes/ (route definitions)
  │   ├── testutil/ (testing utilities & fixtures)
  │   └── utils/ (internal utilities - logging, etc.)
  ├── pkg/ (reusable utilities)
  ├── logs/ (application log files)
  └── migrations/ (SQL migration files)
```

### Key Architectural Patterns

#### 1. **DTO Layer (Data Transfer Objects)** ⭐ NEW
- Located in `internal/dto/`
- Separate API contracts from domain models
- Each DTO has a `Validate()` method for self-validation
- Examples: `dto.RegisterRequest`, `dto.CreateBookRequest`, `dto.LoginResponse`
- Benefits: Clean separation, explicit validation, better testability

Example:
```go
type RegisterRequest struct {
    Name     string `json:"name" example:"John Doe"`
    Email    string `json:"email" example:"john@example.com"`
    Password string `json:"password" example:"password123"`
}

func (r *RegisterRequest) Validate() error {
    if strings.TrimSpace(r.Name) == "" {
        return errors.New("name is required")
    }
    // ... more validations
    return nil
}
```

#### 2. **Handlers (HTTP Layer)**
- Located in `internal/handlers/`
- Convert HTTP requests to DTOs and call services
- Use `fiber.Ctx` for request/response handling
- Self-validate DTOs using `req.Validate()`
- Log operations with structured logging (`utils.InfoLogger`, `utils.ErrorLogger`)
- Return standardized responses via utility functions
- Include Swagger annotations for API documentation

Example flow: `Register` handler → parses to `dto.RegisterRequest` → validates → calls `AuthService.Register()` → returns JSON

#### 3. **Services (Business Logic Layer)**
- Located in `internal/services/`
- Contain all business logic
- Accept DTOs as input, return domain models
- Use **explicit dependency injection**: `NewAuthService(db *gorm.DB)`
- Access database through injected GORM instance
- Services are stateless and created per request

Example:
```go
authService := services.NewAuthService(database.GetDB())
user, err := authService.Register(&req) // req is dto.RegisterRequest
```

#### 4. **Models (Domain Models)**
- Located in `internal/models/`
- Pure domain entities representing database tables
- GORM tags for database mapping
- No request/response structs (moved to DTO layer)
- Examples: `User`, `Book`, `APIResponse`

#### 5. **Middleware Chain**
- Located in `internal/middleware/`
- Auth middleware extracts and validates JWT tokens
- Error handling middleware catches panics and formats errors
- Set in `main.go:setupMiddleware()` for global middleware
- Route-specific middleware applied via `group.Use(middleware.AuthMiddleware())`

#### 5. **Database Layer**
- Located in `internal/database/`
- `db.go`: Connection and migration orchestration
- `migrator.go`: SQL migration runner from embedded files
- `seeder.go`: Sample data seeding
- Uses `embed.FS` to bundle migration files into binary for production

#### 6. **Configuration Management**
- Located in `config/`
- `config.go`: Load from `.env` using `godotenv`, validate, expose via `config.AppConfig`
- Database connection created in `database.go` based on `DBDriver` (postgres/sqlite)
- JWT secret, timeouts, CORS settings all configured here

#### 7. **JWT Authentication**
- Located in `pkg/jwt/`
- Token generation and validation
- Access tokens (short-lived, 15min default) and refresh tokens (7d default)
- Claims include `UserID`, `Email`, `Role`
- Auth middleware validates token and extracts user info into context

#### 8. **Testing Infrastructure** ⭐ NEW
- Located in `internal/testutil/`
- `db.go`: In-memory SQLite database setup for tests
- `fixtures.go`: Reusable test data creation functions
- `assert.go`: Custom assertion helpers

Example usage:
```go
func TestCreateBook(t *testing.T) {
    db := testutil.SetupTestDB(t)
    defer testutil.CleanupTestDB(db)

    book := testutil.CreateBookFixture(db, "Test Book", "Author", "ISBN123", 2024)
    testutil.AssertEqual(t, "Test Book", book.Title)
}
```

Available utilities:
- `SetupTestDB(t)` - Create in-memory test database
- `CreateUserFixture()`, `CreateBookFixture()` - Generate test data
- `AssertEqual()`, `AssertStatusCode()`, `ParseJSONResponse()` - Test assertions

#### 9. **Structured Logging** ⭐ NEW
- Located in `internal/utils/logger.go`
- Two loggers: `InfoLogger` (general operations) and `ErrorLogger` (errors)
- Logs written to both file (`logs/app.log`) and console
- Format: `[LEVEL] timestamp [HandlerName] Message: details`

Example usage in handlers:
```go
utils.InfoLogger.Printf("[Register] User registered successfully: %s (ID: %d)", user.Email, user.ID)
utils.ErrorLogger.Printf("[Register] Validation failed: %v", err)
```

#### 10. **Swagger/OpenAPI Documentation** ⭐ NEW
- Auto-generated from code annotations
- Located in `docs/` directory (swagger.json, swagger.yaml)
- Access interactive UI at `/swagger/index.html`
- Annotations on handlers define endpoints, parameters, responses

Example annotation:
```go
// Register godoc
//  @Summary      Register a new user
//  @Description  Register a new user with name, email, and password
//  @Tags         Authentication
//  @Accept       json
//  @Produce      json
//  @Param        request body dto.RegisterRequest true "Registration request"
//  @Success      201 {object} models.APIResponse "User registered successfully"
//  @Router       /auth/register [post]
func Register(c *fiber.Ctx) error { ... }
```

### Request/Response Flow (Updated with DTO Layer)

1. **Request arrives** → Fiber routes to handler in `internal/handlers/`
2. **Handler parses** request body into DTO (e.g., `dto.RegisterRequest`)
3. **Handler validates** DTO using `req.Validate()` method
4. **Handler logs** operation start with `utils.InfoLogger`
5. **Handler creates** service with explicit DI: `services.NewAuthService(database.GetDB())`
6. **Handler calls** service method, passing DTO
7. **Service executes** business logic: database queries, validations, token generation
8. **Service returns** domain objects (models.User) or errors
9. **Handler logs** operation result (success with InfoLogger, error with ErrorLogger)
10. **Handler formats** response using utility functions (`utils.SuccessResponse()`, `utils.CreatedResponse()`, etc.)
11. **Response sent** as JSON with standardized structure

**Example Complete Flow:**
```go
// 1. Parse request to DTO
var req dto.RegisterRequest
c.BodyParser(&req)

// 2. Validate DTO
if err := req.Validate(); err != nil {
    utils.ErrorLogger.Printf("[Register] Validation failed: %v", err)
    return utils.BadRequestResponse(c, err.Error())
}

// 3. Call service with DI
authService := services.NewAuthService(database.GetDB())
user, err := authService.Register(&req)

// 4. Log and respond
utils.InfoLogger.Printf("[Register] User registered: %s (ID: %d)", user.Email, user.ID)
return utils.CreatedResponse(c, "User registered successfully", user.GetPublicUser())
```

### Database Migrations

The project uses a **dual-migration system**:

- **Development**: `AutoMigrate` in `internal/database/db.go` for fast iteration
- **Production**: SQL files in `migrations/` directory embedded via `embed.go`

Run migrations via:
```bash
go run main.go -migrate        # AutoMigrate (dev)
go run main.go -migrate=sql    # SQL migrations (prod)
```

Migrations track applied status in `schema_migrations` table.

### Configuration & Environment

- Load configuration in `main.go` via `config.LoadConfig()`
- Read from `.env` file (or system env vars)
- Copy `.env.example` to `.env` and customize
- Critical settings: `PORT`, `DB_DRIVER`, `DB_HOST`, `JWT_SECRET`, `DB_USER`, `DB_PASSWORD`
- Different port mapping for Docker: `5432` internal → `6543` host-accessible

## Adding New Features

### Adding a New Entity (e.g., "Product")

1. **Create DTOs** in `internal/dto/product_dto.go`
   - Define request DTOs: `CreateProductRequest`, `UpdateProductRequest`
   - Define response DTOs: `ProductResponse` (optional)
   - Add `Validate()` method to each request DTO
   - Include Swagger example tags (`example:"..."`)

2. **Create Model** in `internal/models/product.go`
   - Define GORM model struct with tags (domain entity)
   - Optionally add helper methods (e.g., `GetPublicProduct()`)

3. **Create Service** in `internal/services/product_service.go`
   - Implement business logic (CRUD operations)
   - Accept DTOs as parameters, return models
   - Use explicit DI: `NewProductService(db *gorm.DB)`

4. **Create Handler** in `internal/handlers/product.go`
   - Implement HTTP endpoints with Swagger annotations
   - Parse requests to DTOs, validate with `req.Validate()`
   - Add structured logging (`utils.InfoLogger`, `utils.ErrorLogger`)
   - Call service with: `services.NewProductService(database.GetDB())`
   - Format responses using utility functions

5. **Add Routes** in `internal/routes/routes.go`
   - Define new route groups (e.g., `/api/products`)
   - Add auth middleware if needed

6. **Generate Swagger** docs
   - Run `make swagger` to update API documentation

7. **Create Migration** in `migrations/` (optional for SQL migrations)
   - Or rely on `AutoMigrate` in development

8. **Write Tests** in `internal/handlers/product_test.go` (co-located)
   - Use `testutil.SetupTestDB()` for test database
   - Use `testutil.CreateProductFixture()` for test data
   - Use `testutil.AssertEqual()` and other assertion helpers

### Modifying Existing Features

- Keep changes isolated to relevant layers (model → service → handler)
- Update migrations if schema changes
- Add tests for behavioral changes
- Run `make test-coverage` to ensure coverage

## Testing

The project uses `testify` and Go's standard `testing` package.

- Run all tests: `make test`
- Run specific test: `go test -v ./tests -run TestName`
- Generate coverage: `make test-coverage`
- Coverage report opens as HTML

## Docker Development Workflow

**Recommended setup for active development:**

```bash
make docker-dev              # Terminal 1: Start containers with hot reload
make docker-dev-logs         # Terminal 2: Monitor logs
# Edit code in your editor
# Air automatically rebuilds on file save
# Refresh browser to see changes
```

Benefits:
- Database runs in container (no local PostgreSQL needed)
- Hot reload: changes instantly visible (no rebuild)
- Production-like environment
- Easy to reset with `make docker-dev-reset`

**For production testing:**
```bash
make docker-up               # Build production image, run compiled binary
```

## Important Files & Their Purposes

| File | Purpose |
|------|---------|
| `main.go` | Entry point; sets up logger, middleware, routes, migrations, Swagger |
| `embed.go` | Embeds migration files into binary |
| `config/config.go` | Load and validate environment configuration |
| `config/database.go` | GORM connection setup based on driver |
| **DTO Layer (NEW)** | |
| `internal/dto/auth_dto.go` | Authentication DTOs with self-validation (Register, Login, Refresh) |
| `internal/dto/user_dto.go` | User operation DTOs (UpdateProfile, ChangePassword) |
| `internal/dto/book_dto.go` | Book operation DTOs (Create, Update) |
| **Handlers** | |
| `internal/handlers/auth.go` | Auth endpoints with DTOs, logging, Swagger annotations |
| `internal/handlers/books.go` | Book endpoints with DTOs, logging, Swagger annotations |
| **Services** | |
| `internal/services/auth_service.go` | Auth business logic with explicit DI |
| `internal/services/book_service.go` | Book business logic with explicit DI |
| **Infrastructure** | |
| `internal/routes/routes.go` | Route definitions, Swagger endpoint |
| `internal/middleware/auth.go` | JWT validation and user context extraction |
| `internal/middleware/error.go` | Global error handling |
| `internal/utils/logger.go` | Structured logging (InfoLogger, ErrorLogger) |
| **Testing (NEW)** | |
| `internal/testutil/db.go` | Test database setup utilities |
| `internal/testutil/fixtures.go` | Test data creation functions |
| `internal/testutil/assert.go` | Custom test assertion helpers |
| **Utilities** | |
| `pkg/jwt/manager.go` | JWT token creation and validation |
| `pkg/utils/responses.go` | Standard response formatting |
| **Documentation (NEW)** | |
| `docs/swagger.json` | OpenAPI spec (auto-generated) |
| `docs/docs.go` | Embedded Swagger documentation |
| **Configuration** | |
| `.env.example` | Template for environment variables |
| `.air.toml` | Hot reload configuration |

## Key Dependencies

- **Fiber v2**: Web framework (Express.js-like for Go)
- **GORM**: ORM for database operations
- **golang-jwt**: JWT token handling
- **bcrypt**: Password hashing
- **godotenv**: .env file loading
- **swaggo/swag**: Swagger documentation generation ⭐ NEW
- **gofiber/swagger**: Swagger UI middleware for Fiber ⭐ NEW

## Code Organization Principles

1. **Single Responsibility**: Each file/function has one reason to change
2. **DTO Layer Separation**: API contracts (DTOs) separated from domain models
3. **Self-Validating DTOs**: DTOs contain their own validation logic
4. **Dependency Injection**: Services receive database explicitly, don't create it
5. **Error Handling**: Return errors explicitly, use middleware for global handling
6. **Structured Logging**: Consistent logging format with context (`[HandlerName] Message`)
7. **API Documentation**: Swagger annotations on all public endpoints
8. **Validation**: Validate input early in handlers using DTO.Validate(); return 400 for bad requests
9. **No Sensitive Data in Logs**: Password hashes are excluded from JSON response
10. **Standard Response Format**: All endpoints return consistent JSON structure

## Concurrent Programming Patterns

> **Note**: The concurrent programming patterns have been temporarily moved to separate documentation to keep the core boilerplate focused and clean. These educational examples demonstrate common Go concurrency patterns and will be available as optional examples or in separate documentation.

The patterns were previously available in `internal/services/concurrent_service.go` and `internal/handlers/concurrent.go`. If you need these examples, they can be found in the git history or will be documented separately.

### Patterns Included (For Reference)

These patterns demonstrate real-world Go concurrency use cases:

#### 1. **Basic Goroutines with WaitGroup**
**Endpoint:** `GET /api/concurrent/parallel?ids=1,2,3`

Process multiple items simultaneously using goroutines and `sync.WaitGroup`.

**Use Cases:**
- Parallel data fetching from database
- Batch processing of independent tasks
- Concurrent API calls

**Key Concepts:**
- `go func()` to launch goroutines
- `sync.WaitGroup` to wait for completion
- `sync.Mutex` to protect shared data

#### 2. **Worker Pool Pattern**
**Endpoint:** `GET /api/concurrent/worker-pool?ids=1,2,3&workers=3`

Limit concurrent operations using a fixed number of workers processing jobs from a queue.

**Use Cases:**
- Rate limiting external API calls
- Database connection pooling
- Background job processing

**Key Concepts:**
- Job queue (buffered channel)
- Fixed number of worker goroutines
- Results collection channel

#### 3. **Fan-Out/Fan-In Pattern**
**Endpoint:** `GET /api/concurrent/fan-out-fan-in?query=golang`

Split work across multiple goroutines, then merge results.

**Use Cases:**
- Multi-source data aggregation
- Parallel searches across different fields
- Distributed computation

**Key Concepts:**
- Multiple goroutines produce results (fan-out)
- Single goroutine collects results (fan-in)
- Result deduplication

#### 4. **Pipeline Pattern**
**Endpoint:** `GET /api/concurrent/pipeline`

Process data through multiple stages using channels.

**Use Cases:**
- Multi-stage data processing
- ETL (Extract, Transform, Load) pipelines
- Stream processing

**Key Concepts:**
- Chained channels between stages
- Each stage is a separate goroutine
- Data flows through pipeline

#### 5. **Semaphore Pattern (Rate Limiting)**
**Endpoint:** `POST /api/concurrent/bulk-create`

Control concurrency using a buffered channel as a semaphore.

**Use Cases:**
- API rate limiting
- Resource pooling (e.g., database connections)
- Controlled parallel writes

**Key Concepts:**
- Buffered channel as semaphore
- Acquire/release pattern
- Concurrent operation limiting

**Example Request:**
```json
{
  "books": [
    {"title": "Book 1", "author": "Author 1", "isbn": "ISBN1"},
    {"title": "Book 2", "author": "Author 2", "isbn": "ISBN2"}
  ],
  "max_concurrent": 3
}
```

#### 6. **Timeout Pattern**
**Endpoint:** `GET /api/concurrent/timeout/1?timeout=5`

Cancel operations that exceed a time limit using context.

**Use Cases:**
- External API calls with timeout
- Slow database queries
- User-facing operations requiring responsiveness

**Key Concepts:**
- `context.WithTimeout()`
- `select` statement with `<-ctx.Done()`
- Graceful timeout handling

#### 7. **Select with Multiple Channels**
**Endpoint:** `GET /api/concurrent/monitor/1?interval=2&duration=10`

Handle multiple channel operations simultaneously.

**Use Cases:**
- Event handling systems
- Real-time monitoring
- Pub/sub implementations

**Key Concepts:**
- `select` statement for channel multiplexing
- Ticker for periodic operations
- Context for cancellation

### Testing Concurrent Patterns

To test these patterns:

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

3. **Test a pattern (example: worker pool):**
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://localhost:4000/api/concurrent/worker-pool?ids=1,2,3,4,5&workers=3"
   ```

### Best Practices for Concurrent Code

1. **Always close channels** when done producing to avoid goroutine leaks
2. **Use context for cancellation** to handle timeouts and cleanup
3. **Protect shared data** with mutexes or use channels for communication
4. **Avoid goroutine leaks** by ensuring all goroutines can exit
5. **Handle errors properly** in goroutines (use error channels)
6. **Use `defer` for cleanup** (e.g., `defer wg.Done()`)
7. **Test with race detector** (`go test -race`)

### Common Pitfalls to Avoid

- **Not waiting for goroutines to finish** → incomplete work
- **Accessing shared memory without synchronization** → race conditions
- **Forgetting to close channels** → goroutine leaks
- **Deadlocks from improper channel usage** → program hangs
- **Creating too many goroutines** → resource exhaustion

### Additional Resources

- Source code: `internal/services/concurrent_service.go`
- Handler implementations: `internal/handlers/concurrent.go`
- Route definitions: `internal/routes/routes.go:47-61`

## Development Checklist for New Code

- [ ] Test locally: `make test`
- [ ] Format code: `make fmt`
- [ ] Check for issues: `make vet` and `make lint`
- [ ] Verify migrations: `make migrate-status`
- [ ] Run in Docker: `make docker-dev`
- [ ] Add/update tests for new functionality
- [ ] Update `.env.example` if new env vars are needed
