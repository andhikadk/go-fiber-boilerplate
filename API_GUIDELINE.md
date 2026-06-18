# API Guideline

This document is the blueprint for building new REST APIs from this template.

## Base Stack

- Go
- Fiber v2
- GORM
- PostgreSQL
- SQL migration files
- JWT access and refresh tokens
- Swagger generation via `swaggo/swag`
- Scalar API Reference
- Structured JSON logging
- Optional Redis cache/rate-limit storage
- Docker Compose for development and production runtime

Optional integrations such as Redis, SMTP email, object storage, Sentry, and workers must be safe by default: an empty configuration means disabled/no-op, not application startup failure.

## Project Structure

```text
cmd/api
config
assets/migrations
assets/migrations/seeds
docs
internal/cache
internal/database
internal/dto
internal/handlers
internal/middleware
internal/models
internal/routes
internal/services
internal/testutil
internal/workers
pkg/jwt
pkg/utils
```

## Layering

Request flow:

```text
HTTP route -> handler -> DTO validation -> service -> GORM/repository -> model -> standard response
```

Rules:

- DTOs define API contracts and own `Validate()`.
- Handlers only parse input, call validation, read auth context, log, map errors, and return responses.
- Services contain business logic, transactions, and external integration orchestration.
- Models are database entities only.
- Routes are the composition root and create services/handlers once during startup.

## Naming

| Layer | Pattern |
| --- | --- |
| DTO | `*_dto.go` |
| Handler | `*_handler.go` |
| Service | `*_service.go` |
| Model | `{entity}.go` |

Handler:

```go
type Resource struct {
	resourceService services.ResourceService
}

func NewResource(resourceService services.ResourceService) *Resource {
	return &Resource{resourceService: resourceService}
}
```

Service:

```go
type ResourceService interface {
	CreateResource(userID uint, req *dto.CreateResourceRequest) (*dto.ResourceResponse, error)
}

type resourceService struct {
	db *gorm.DB
}

func NewResourceService(db *gorm.DB) ResourceService {
	return &resourceService{db: db}
}
```

## Required Core Services

- `AuthService`: register, login, refresh, forgot password, reset password.
- `UserService`: get profile, update profile, change password.
- `ResourceService`: generic CRUD sample for future feature reference.
- `EmailService`: interface + no-op implementation.
- `StorageService`: interface + no-op placeholder.
- `cache.Client`: Redis-backed cache with no-op fallback.

Do not copy product-specific modules into the template unless they are generic and disabled by default.

## Migrations

SQL migration files are the source of truth.

```text
assets/migrations/001_initial_schema.sql
assets/migrations/002_add_indexes.sql
```

Rules:

- Do not use AutoMigrate for runtime schema.
- Use `make migrate` for pending migrations.
- Use `make migrate-fresh` only in development.
- Keep models synchronized with SQL files.
- Add indexes for foreign keys and frequently filtered columns.

## Middleware

Global middleware:

- request ID
- request context propagation
- panic recovery
- CORS
- Helmet/security headers
- global limiter
- compression
- structured access log with redaction
- error handling middleware

Route middleware:

- auth limiter on `/api/auth`
- `AuthMiddleware()` on protected groups
- `RequireRoles()` when role-specific access is needed

## Response Format

Success:

```json
{
  "status": 200,
  "message": "Success",
  "data": {}
}
```

Error:

```json
{
  "status": 400,
  "message": "Invalid request body",
  "error": "Invalid request body"
}
```

Always use helpers in `pkg/utils/response.go`.

## Error Handling

- Define service sentinel errors for business conditions.
- Compare errors with `errors.Is`.
- Map errors to HTTP status codes in handlers.
- Return generic messages for unexpected `500` responses.
- Log details with `utils.LogCtx`.

## Feature Checklist

1. Add SQL migration.
2. Add or update the model.
3. Add request/response DTOs.
4. Add the service interface and implementation.
5. Add a struct-based handler.
6. Wire the service and handler in routes.
7. Add Swagger annotations.
8. Add tests.
9. Run `make fmt`, `make vet`, `make test`, and `make build`.

## AI Workflow

AI agents may build, test, format, vet, lint, and generate Swagger. AI agents must not run the app, hot reload server, Docker Compose, or any long-running server process.
