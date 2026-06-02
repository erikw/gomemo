# Gomemo Implementation Plan (v1)

## Philosophy

Optimize for **learning progression**, not for building features as quickly as possible.

Introduce one new concept at a time while keeping the project runnable from day one.

---

# Phase 0: Project Setup

## Goal

Have a compilable application with a clear structure.

## Initial Structure

```text
gomemo/
├── cmd/
│   └── gomemo/
│       └── main.go
├── internal/
│   ├── config/
│   ├── notes/
│   └── storage/
├── go.mod
└── README.md
```

## Tasks

Initialize the module:

```bash
go mod init github.com/you/gomemo
```

Add dependencies:

```bash
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
```

## Learn

* Go modules
* Package organization
* Application entrypoint

## Deliverable

```bash
go run ./cmd/gomemo
```

prints:

```text
Starting Gomemo...
```

---

# Phase 1: Basic HTTP Server

## Goal

Learn `net/http` and Chi before introducing databases.

## Tasks

Create:

```http
GET /
```

Returns:

```json
{
  "status": "ok"
}
```

Add a Chi router.

Create:

```go
func NewRouter() http.Handler
```

## Learn

* Handlers
* Routing
* `http.ResponseWriter`
* `*http.Request`

## Deliverable

```bash
curl localhost:8080/
```

returns a valid response.

---

# Phase 2: In-Memory Notes Store

## Goal

Learn API design without database complexity.

## Note Model

```go
type Note struct {
    ID        int64
    Title     string
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Storage

Use an in-memory store:

```go
map[int64]Note
```

or

```go
[]Note
```

## Endpoints

```http
POST   /notes
GET    /notes
GET    /notes/{id}
DELETE /notes/{id}
```

## Learn

* JSON encoding
* JSON decoding
* Route parameters
* Status codes
* Request validation

## Deliverable

A fully working CRUD API backed only by memory.

No database yet.

---

# Phase 3: Proper Error Handling

## Goal

Make API responses consistent and production-like.

## Tasks

Return structured errors:

```json
{
  "error": "note not found"
}
```

Handle:

* Invalid JSON
* Missing required fields
* Unknown note IDs

Use wrapped errors:

```go
fmt.Errorf("create note: %w", err)
```

## Learn

* `errors.Is`
* `errors.As`
* Error wrapping
* Mapping domain errors to HTTP responses

## Deliverable

Consistent error responses across the API.

---

# Phase 4: Structured Logging

## Goal

Introduce observability basics.

## Tasks

Add `slog`.

Log:

* Startup
* Requests
* Failures

Example:

```text
INFO request
method=POST
path=/notes
duration=4ms
```

## Learn

* Structured logging
* Log levels
* Request metadata

## Deliverable

Readable, structured logs.

---

# Phase 5: Configuration

## Goal

Separate runtime configuration from code.

## Example Environment Variables

```env
PORT=8080
```

Later:

```env
DATABASE_URL=postgres://...
```

## Example

```go
type Config struct {
    Port string
}
```

## Learn

* Environment variables
* Startup configuration
* Configuration validation

## Deliverable

Server starts using configuration loaded from the environment.

---

# Phase 6: PostgreSQL

## Goal

Replace the in-memory store with a real database.

## Migration

```sql
CREATE TABLE notes (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
```

## Tasks

* Choose Goose or golang-migrate
* Add migrations
* Create pgx connection pool
* Implement CRUD with SQL

## Learn

* PostgreSQL
* SQL
* Migrations
* Connection pooling
* Schema design

## Deliverable

All CRUD operations backed by PostgreSQL.

---

# Phase 7: Repository Layer

## Goal

Introduce abstraction once there is a concrete need.

## Example

```go
type NoteRepository interface {
    Create(...)
    Get(...)
    List(...)
    Delete(...)
}
```

Implementation:

```go
type PostgresNoteRepository struct {
    ...
}
```

## Why Here?

You will better understand why the abstraction exists after building the database-backed implementation.

Avoid introducing interfaces prematurely.

## Learn

* Dependency injection
* Interfaces
* Testability

## Deliverable

Clean separation between storage and HTTP layers.

---

# Phase 8: Contexts

## Goal

Use request contexts properly throughout the application.

## Flow

```text
handler
  ↓
service
  ↓
repository
  ↓
pgx
```

Pass:

```go
r.Context()
```

through every layer.

Add timeouts where appropriate.

## Learn

* Context propagation
* Cancellation
* Deadlines
* Request-scoped operations

## Deliverable

Database operations respect request cancellation and timeouts.

---

# Phase 9: Testing

## Goal

Build confidence in the application through automated tests.

## Handler Tests

Using:

```go
httptest
```

Example:

```go
POST /notes
```

returns:

```http
201 Created
```

## Integration Tests

Use a real PostgreSQL database.

## Learn

* Table-driven tests
* Handler testing
* Integration testing
* Race detector

## Deliverable

Coverage for both success and failure paths.

---

# Phase 10: Graceful Shutdown

## Goal

Handle process termination correctly.

## Tasks

Handle:

```text
SIGINT
SIGTERM
```

Shutdown cleanly:

```go
server.Shutdown(...)
```

Close the database pool.

## Learn

* Signals
* Cleanup
* Production readiness

## Deliverable

Ctrl+C stops the service cleanly.

---

# Gomemo v1 Definition

## Included

* Create note
* List notes
* Get note
* Delete note
* PostgreSQL
* Migrations
* Structured logging
* Configuration management
* Context propagation
* Automated tests
* Graceful shutdown

## Explicitly Excluded

* Authentication
* Redis
* Docker
* Search
* Pagination
* Tags
* Metrics
* OpenTelemetry
* Background workers

These belong in future iterations after the core backend fundamentals are complete.

---

# Success Criteria

By the end of Gomemo v1, I should be comfortable with:

* Building a Go HTTP API from scratch
* Using Chi effectively
* Working with PostgreSQL via pgx
* Designing clean package structures
* Writing tests
* Using `context.Context` correctly
* Implementing structured logging
* Performing graceful shutdown
* Explaining architectural decisions in interviews

The project should remain intentionally small, focused, and complete.


