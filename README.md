# Go Fiber Boilerplate

A production-ready boilerplate for building REST APIs with **Fiber**, a fast and lightweight Go web framework inspired by Express.js.

## 🚀 Features

- **Fiber Web Framework** - Fast, minimalist web framework
- **JWT Authentication** - Secure token-based authentication
- **GORM ORM** - Database abstraction layer
- **PostgreSQL & SQLite** - Multiple database support
- **Middleware Stack** - CORS, Logger, Recovery, Helmet
- **Request Validation** - Struct-based validation
- **Error Handling** - Centralized error management
- **Database Migrations** - Schema versioning
- **Unit Tests** - Testing setup ready
- **Docker Support** - Containerized deployment
- **Hot Reload** - Development mode with air
- **Environment Management** - .env configuration

## 📋 Project Structure

```
go-fiber-boilerplate/
├── cmd/
│   └── main.go                 # Entry point
├── config/
│   ├── config.go               # Configuration management
│   └── database.go             # Database setup
├── internal/
│   ├── handlers/               # HTTP handlers
│   ├── models/                 # Data structures
│   ├── services/               # Business logic
│   ├── middleware/             # Custom middlewares
│   ├── database/               # Database layer
│   └── routes/                 # Route definitions
├── pkg/
│   ├── utils/                  # Utility functions
│   └── jwt/                    # JWT utilities
├── migrations/                 # Database migrations
├── tests/                      # Test files
├── .env.example                # Environment template
├── go.mod & go.sum             # Dependencies
├── Dockerfile                  # Container image
├── docker-compose.yml          # Compose configuration
├── Makefile                    # Build commands
└── README.md                   # This file
```

## 🛠️ Tech Stack

- **Framework:** Fiber v2
- **Database:** GORM, PostgreSQL, SQLite
- **Authentication:** JWT (golang-jwt)
- **Security:** bcrypt (golang.org/x/crypto)
- **Validation:** go-playground/validator
- **Testing:** testify, standard library
- **Environment:** godotenv
- **Middleware:** Fiber built-in + custom

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

### Prerequisites
- Go 1.25 or higher
- PostgreSQL 12+ (or SQLite for development)
- Make (optional, for using Makefile)

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

4. **Setup database**
```bash
# Using Docker
make docker-up

# Or manually create PostgreSQL database
createdb fiber_boilerplate
```

5. **Run application**
```bash
make run
# or
go run cmd/main.go
```

The API will be available at `http://localhost:3000`

## 🚀 Usage

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
make migrate             # Run migrations
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

## 📝 Configuration

All configuration is managed through `.env` file. See `.env.example` for all available options.

### Key Configuration
- `PORT` - Server port (default: 3000)
- `ENV` - Environment (development/production)
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

### Build and Run with Docker Compose
```bash
docker-compose up -d
```

This will:
- Build the Fiber application
- Start PostgreSQL database
- Create necessary tables
- Expose API on port 3000

### Stop
```bash
docker-compose down
```

## 📖 Project Structure Details

### `cmd/main.go`
Application entry point. Initializes config, database, and starts the server.

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

1. **Create models** in `internal/models/`
2. **Create handlers** in `internal/handlers/`
3. **Add business logic** in `internal/services/`
4. **Define routes** in `internal/routes/`
5. **Write tests** in `tests/`
6. **Run and test** with `make run` and `make test`

## 📚 Learning Resources

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Guide](https://gorm.io/docs/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)
- [Go Best Practices](https://golang.org/doc/effective_go)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the MIT License.

## 👨‍💻 Author

Your Name - [@yourtwitter](https://twitter.com/yourtwitter)

## 🙏 Acknowledgments

- Fiber team for the amazing framework
- GORM team for the powerful ORM
- Go community for best practices

---

**Happy coding! 🚀**

For issues and questions, please open an issue on GitHub.
