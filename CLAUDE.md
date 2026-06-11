# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
DB_USER=root DB_PASSWORD=secret DB_NAME=accounting go run ./cmd/api

# Build
go build ./cmd/api

# Apply migration
mysql -u root -p accounting < migrations/001_create_accounts.sql

# Add/tidy dependencies
go mod tidy
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | `` | MySQL password |
| `DB_NAME` | `accounting` | Database name |
| `PORT` | `8080` | HTTP server port |

## Architecture

Standard layered Go architecture: **handler → service → repository → MySQL**.

Each domain (currently: `accounts`) has one file per layer. When adding a new domain, create a file in each of the four packages and wire it up in `cmd/api/main.go`.

```
cmd/api/main.go          — wires DB, repo, service, handler, starts gin server
internal/model/          — plain structs and request types (no logic)
internal/repository/     — database/sql queries against MySQL; interface + postgres impl in same file
internal/service/        — business logic; depends on repository interface
internal/handler/        — gin handlers; router.go registers all routes under /api
migrations/              — raw SQL files, applied manually
```

### Key conventions

- MySQL placeholders are `?` (not `$1`). MySQL has no `RETURNING`; use `LastInsertId()` + `GetByID` after INSERT, and `RowsAffected()` after UPDATE to detect 404.
- DSN uses `parseTime=true` so `time.Time` fields scan correctly from MySQL `DATETIME` columns.
- `GetAll` returns `[]model.Account{}` (not nil) so the JSON response is always an array.
- `updated_at` is managed by MySQL `ON UPDATE CURRENT_TIMESTAMP` — no need to set it manually in queries.
- Gin handlers use `c.ShouldBindJSON` for request decoding and `c.JSON` / `c.Status` for responses.
