# Memwright API

The Memwright API is a Go-based REST API server for the Memwright spaced repetition learning platform.

## Requirements

- Go 1.16 or higher
- PostgreSQL 14+

## Quick Start

### 1. Configure Environment

Copy the sample environment file and update the values:

```bash
cp .env.sample .env
```

Edit `.env` with your configuration:

```bash
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=memwright
DATABASE_PASSWORD=your_password_here
DATABASE_NAME=memwright
DATABASE_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here_min_32_characters
JWT_EXPIRATION_HOURS=24

# Log level (debug, info, warn, error)
LOG_LEVEL=debug
```

### 2. Run the Server

```bash
go run ./cmd/server
```

The server will start on `http://localhost:8080` by default.

### 3. Verify Installation

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok","environment":"development"}
```

## Building

```bash
# Build the binary
go build -o bin/server ./cmd/server

# Run the binary
./bin/server
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run only unit tests
go test ./tests/unit/...

# Run only integration tests
go test ./tests/integration/...
```

## Configuration Reference

The API is configured via environment variables. Create a `.env` file in the `api/` directory or set them directly in your environment.

### Server

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `ENVIRONMENT` | Runtime environment (`development`, `staging`, `production`) | `development` |
| `LOG_LEVEL` | Minimum log level (`debug`, `info`, `warn`, `error`) | `info` |

### Database

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | PostgreSQL host | `localhost` |
| `DATABASE_PORT` | PostgreSQL port | `5432` |
| `DATABASE_USER` | Database username | `memwright` |
| `DATABASE_PASSWORD` | Database password | (empty) |
| `DATABASE_NAME` | Database name | `memwright` |
| `DATABASE_SSLMODE` | SSL mode (`disable`, `require`, `verify-full`) | `disable` |

### Authentication

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for JWT token signing (min 32 characters) | (empty) |
| `JWT_EXPIRATION_HOURS` | JWT token expiration time in hours | `24` |

## Project Structure

```
api/
├── cmd/server/               # Application entrypoint
│   └── main.go
├── internal/                 # Private application code
│   ├── config/               # Configuration management
│   ├── handler/              # HTTP handlers (controllers)
│   ├── model/                # Domain types and entities
│   ├── repository/           # Database queries and data access
│   ├── service/              # Business logic layer
│   └── srs/                  # Spaced repetition algorithms (SM-2, FSRS)
├── migrations/               # SQL database migrations
├── pkg/                      # Shared libraries
│   ├── env/                  # Environment variable loader
│   └── logger/               # Structured logging
├── tests/
│   ├── integration/          # Integration tests
│   └── unit/                 # Unit tests
├── .env                      # Environment variables (gitignored)
├── .env.sample               # Sample environment file
├── go.mod
└── README.md
```

### Architecture Layers

1. **Handler**: Receives HTTP requests, validates input, calls services, returns responses
2. **Service**: Contains business logic, orchestrates repositories, applies domain rules
3. **Repository**: Database operations, SQL queries, data mapping
4. **Model**: Domain entities and DTOs

## API Endpoints

### Health Check

```
GET /health
```

Returns the server health status. Used for container orchestration readiness checks.

**Response:**
```json
{
  "status": "ok",
  "environment": "development"
}
```

## Development

### Adding a New Endpoint

1. Create or update the model in `internal/model/`
2. Add repository methods in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Create the HTTP handler in `internal/handler/`
5. Register the route in `cmd/server/main.go`
6. Add tests in `tests/unit/` and `tests/integration/`

### Logger Usage

The logger supports four levels: `debug`, `info`, `warn`, `error`.

```go
import "memwright/api/pkg/logger"

log := logger.New(
    logger.WithDebug(true),
    logger.WithEnvironment("development"),
)

log.Debug("detailed information for debugging")
log.Info("general information")
log.Warn("warning message")
log.Error("error occurred: %v", err)
```

### Environment Variables

The `pkg/env` package automatically loads variables from `.env`:

```go
import "memwright/api/pkg/env"

env.Load()  // Loads from .env in current directory
```

Existing environment variables are not overwritten.

## License

See the [LICENSE](../LICENSE) file in the root directory.

