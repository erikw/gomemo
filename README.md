# gomemo

[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue)](/LICENSE)

A lightweight note-taking API server written in Go with a clean, modular architecture.

## Features

- **RESTful API** for note management built with [chi](https://github.com/go-chi/chi)
- **Clean Architecture** with separation of concerns (handlers, services, storage)
- **Pluggable Storage** with in-memory storage engine (ready for extensibility)
- **Environment Configuration** with support for dev and production modes
- **Structured Logging** using Go's standard `slog` package
- **Development Tools** with live reload via [air](https://github.com/air-verse/air)

## Quick Start

### Prerequisites

- Go 1.26.3 or later
- GNU Make (optional, for convenience)

### Build

```bash
make build
```

This produces the `gomemo` binary with version information embedded from git tags.

### Run

```bash
./gomemo
```

Set custom port:

```bash
PORT=8080 ./gomemo
```

### Development

Run the development server with live reload:

```bash
go tool air
```

Or run directly:

```bash
go run ./cmd/gomemo
```

Enable debug logging:

```bash
go run ./cmd/gomemo -debug
```

### Testing

```bash
make test
```

### Manual API Testing (`.http` files)

Shared `.http` request collections are available under `api/`, similar to a Postman collection but versioned in git.

1. Install the HTTP client tool:

```bash
npm install
```

2. Start the API server (recommended in dev mode so fixture notes are available):

```bash
ENV=dev go run ./cmd/gomemo
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

## Installation

Install globally:

```bash
make install
```

## Available Commands

- `make build` - Build the binary
- `make install` - Install to `$GOPATH/bin`
- `make run` - Run with optional args: `make run ARGS="-debug"`
- `make test` - Run tests
- `make clean` - Remove build artifacts
- `make release` - Create a git tag and release: `make release VERSION=v0.1.0`

## Configuration

The application supports the following environment variables:

- `PORT` - Server port (default: 3000)
- `HOST` - Server address (default: 127.0.0.1)
- `ENV` - Environment mode: `dev` or `prod` (default: prod)

In development mode (`ENV=dev`), the database is automatically seeded with fixtures from `data/dev.yaml`.

## Architecture

The project follows a layered architecture:

- **Handlers** - HTTP request/response handling
- **Services** - Business logic
- **Storage** - Data persistence abstraction
- **Models** - Domain entities

Handlers register routes via the `RouteRegistrar` interface, enabling clean, modular route registration.
