# CLAUDE.md

This repository is a production-ready Go Fiber REST API template.

## Commands

```bash
make build
make test
make test-coverage
make fmt
make vet
make lint
make swagger
make migrate
make migrate-fresh
make migrate-status
make seed
```

## AI Assistant Rules

Allowed:

- Read and analyze files.
- Edit code and docs.
- Run `make fmt`, `make vet`, `make test`, `make build`, and `make swagger`.

Not allowed:

- Do not run the application with `make run`, `go run ./cmd/api`, or similar.
- Do not start hot reload with `make dev` or `air`.
- Do not start Docker containers with `make docker-dev`, `make docker-up`, or `docker compose up`.
- Do not run long-running server processes.

The user controls application execution and manual testing.

## Architecture

- DTOs live in `internal/dto` and own input validation.
- Handlers live in `internal/handlers`, are struct-based, and receive service interfaces by constructor.
- Services live in `internal/services`, expose interfaces, and contain business logic.
- Models live in `internal/models` and represent database entities only.
- Routes in `internal/routes` are the composition root.
- SQL migrations in `assets/migrations` are the source of truth for schema.
- Shared utilities live in `pkg/utils` and `pkg/jwt`.

## Development Rules

- Use SQL migrations for schema changes.
- Keep GORM models synchronized with migrations.
- Use standard response helpers from `pkg/utils/response.go`.
- Use `utils.LogCtx(c.UserContext(), "Module")` in handlers.
- Map service sentinel errors to HTTP status in handlers.
- Never expose raw internal errors in `500` responses.
- Protected endpoints must use `middleware.AuthMiddleware()` and Swagger `@Security BearerAuth`.

See `API_GUIDELINE.md` for the full project guideline.
