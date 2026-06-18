# Migrations

SQL files in this directory are the source of truth for runtime schema.

## Commands

```bash
make migrate
make migrate-fresh
make migrate-status
make seed
```

## Baseline

- `001_initial_schema.sql`: users, user profiles, password resets, resources.
- `002_add_indexes.sql`: indexes for auth, reset tokens, and resources.

Seed files live in `assets/migrations/seeds`.

## Rules

- Add new migrations with sequential numbers.
- Keep GORM models synchronized with SQL schema.
- Do not use AutoMigrate for runtime schema.
- Use `make migrate-fresh` only in development.
- Add indexes for foreign keys and frequently queried columns.
