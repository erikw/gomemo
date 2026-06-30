# gomemo

[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue)](/LICENSE)

A lightweight note-taking API server written in Go with a clean, modular architecture.

## Features

- RESTful notes API built with [chi](https://github.com/go-chi/chi)
- Pluggable storage backends:
  - `memory` (default)
  - `postgres`
- PostgreSQL schema migrations with [golang-migrate](https://github.com/golang-migrate/migrate)
- Development workflow with [air](https://github.com/air-verse/air)

## API Versioning

The API uses semantic versioning via URL prefixes:

- Business endpoints are versioned under `/api/v1`
- Health endpoints remain unversioned (`/`, `/health`)

### API v1 endpoints

- `GET /api/v1/notes`
- `POST /api/v1/notes`
- `GET /api/v1/notes/{noteID}`
- `PATCH /api/v1/notes/{noteID}`
- `DELETE /api/v1/notes/{noteID}`

## Prerequisites

- Go 1.26.3+
- GNU Make (optional)
- Docker + Docker Compose (recommended for PostgreSQL dev setup)
- `migrate` CLI for migrations

### Install `migrate` CLI

On macOS:

```bash
brew install golang-migrate
```

## Quick start (memory storage)

Build:

```bash
make build
```

Run:

```bash
./gomemo serve
```
or 

```bash
make run
```


Memory mode does not require PostgreSQL.

## Quick start (PostgreSQL with Docker Compose)

1. Copy dev env file:

```bash
cp .env.example .env
```

2. Start PostgreSQL:

```bash
make db-up
```

3. Export environment from `.env`:

```bash
set -a
source .env
set +a
```

4. Run migrations:

```bash
make migrate-up
```

5. Seed fixture notes:

```bash
./gomemo seed
```

6. Start server:

```bash
./gomemo serve
```

## Development (air)

The default `.air.toml` workflow expects PostgreSQL mode and runs migrations + seed before starting.

```bash
go tool air
```

Ensure `.env` exists and PostgreSQL is running (`make db-up`) before starting Air.

## Configuration

Environment variables:

- `HOST` (default: `127.0.0.1`)
- `PORT` (default: `8080`)
- `ENV` (default: `prod`)
- `STORAGE_TYPE` (default: `memory`, allowed: `memory`, `postgres`)
- `DATABASE_URL` (required when `STORAGE_TYPE=postgres`)

Example:

```bash
ENV=dev STORAGE_TYPE=postgres DATABASE_URL=postgres://gomemo:gomemo@127.0.0.1:5432/gomemo?sslmode=disable ./gomemo serve
```

## Database migrations

Migration files are in `db/migrations`.

- Apply all pending migrations:

```bash
make migrate-up
```

- Roll back last migration:

```bash
make migrate-down
```

- Show migration version:

```bash
make migrate-version
```

- Create a new migration:

```bash
make migrate-create NAME=add_note_tags
```

All migration targets require `DATABASE_URL` to be set.

## Make targets

- `make build`
- `make test`
- `make run ARGS="serve"`
- `make db-up`
- `make db-down`
- `make db-reset`
- `make migrate-up`
- `make migrate-down`
- `make migrate-version`
- `make migrate-create NAME=...`

## Seeding behavior

`gomemo seed` loads fixtures from `data/dev.yaml` and clears existing notes first, so seeding is repeatable in development.

## Testing

```bash
make test
```

### Manual API testing (`.http` files)

Shared `.http` request collections are available under `api/`, similar to a Postman collection but versioned in git.

1. Install the HTTP client tool:

```bash
npm install
```

2. Start the API server:

```bash
ENV=dev go run ./cmd/gomemo serve
```

3. Run a single request file:

```bash
npx httpyac --env dev send -a api/notes/get-notes.http
```

4. Run all request files:

```bash
npx httpyac --env dev send -a "api/**/*.http"
```

`BASE_URL`, `NOTE_ID`, and `MISSING_NOTE_ID` are configured in `.httpyac.js`.
Some requests are intentionally negative test cases and are expected to return `400`/`404`.

## Architecture

- Handlers: HTTP request/response layer
- Services: business logic
- Storage: persistence abstraction
- Models: domain entities
