BINARY=gomemo
CMD=./cmd/gomemo
MIGRATIONS_DIR=./db/migrations
COMPOSE_CMD=$(shell if docker compose version >/dev/null 2>&1; then echo "docker compose"; elif command -v docker-compose >/dev/null 2>&1; then echo "docker-compose"; else echo ""; fi)

VERSION=$(shell git describe --tags --dirty --always)
PKG=github.com/erikw/gomemo

LDFLAGS=-ldflags "-X $(PKG)/internal/version.Version=$(VERSION)"

.PHONY: all build clean run install test release check-compose db-up db-down db-reset migrate-up migrate-down migrate-version migrate-create seed

all: build test

build:
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

clean:
	$(RM) $(BINARY)

# Usage: make run ARGS="-h"
run:
	go run $(LDFLAGS) $(CMD) $(ARGS)

install:
	go install $(LDFLAGS) $(CMD)

test:
	go test ./...

check-compose:
	@if [ -z "$(COMPOSE_CMD)" ]; then echo "Docker Compose is required. Install either 'docker compose' (v2 plugin) or 'docker-compose'."; exit 1; fi

db-up:
	@$(MAKE) check-compose
	$(COMPOSE_CMD) up -d postgres

db-down:
	@$(MAKE) check-compose
	$(COMPOSE_CMD) down

db-reset:
	@$(MAKE) check-compose
	$(COMPOSE_CMD) down -v
	$(COMPOSE_CMD) up -d postgres

migrate-up:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

migrate-version:
	@if [ -z "$(DATABASE_URL)" ]; then echo "DATABASE_URL is required"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

# Usage: make migrate-create NAME=create_notes_table
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "NAME is required"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

seed:
	go run $(LDFLAGS) $(CMD) -d seed

# Usage: make release VERSION=v0.2.0
release:
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
