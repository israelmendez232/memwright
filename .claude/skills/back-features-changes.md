---
name: back-features-changes
description: Work on back-end Go features and API changes based on Jira tickets
---

# Back-End Features and Changes Skill

When the user invokes `/back-features-changes`, provide context about the back-end architecture and then work on the specified Jira ticket.

## Back-End Context

### Stack
- **Language**: Go
- **Router**: net/http or Chi router
- **Database**: PostgreSQL
- **Location**: `/api` directory

### Directory Structure
```
api/
├── cmd/server/       # Application entrypoint
├── internal/
│   ├── handler/      # HTTP handlers (controllers)
│   ├── service/      # Business logic layer
│   ├── repository/   # Database queries and data access
│   ├── model/        # Domain types and entities
│   └── srs/          # Spaced repetition algorithms (SM-2, FSRS)
├── migrations/       # SQL migrations
└── go.mod
```

### Architecture Layers
1. **Handler**: Receives HTTP requests, validates input, calls services, returns responses
2. **Service**: Contains business logic, orchestrates repositories, applies domain rules
3. **Repository**: Database operations, SQL queries, data mapping
4. **Model**: Domain entities and DTOs

### API Design
- RESTful API at `/api/v1/`
- JWT authentication (access + refresh tokens)
- Structured JSON responses
- Rate limiting enabled

### Key Endpoints
```
/api/v1/auth/*           # Authentication (register, login, refresh, logout)
/api/v1/decks/*          # Deck CRUD and subdecks
/api/v1/decks/:id/cards  # Cards within a deck
/api/v1/cards/:id        # Individual card operations
/api/v1/cards/:id/review # Submit review ratings
/api/v1/stats/*          # Dashboard and analytics
/api/v1/settings         # User and deck settings
```

### Core Domain Entities
- **User**: Account with global SRS settings
- **Deck**: Card collection with tree hierarchy (subdecks)
- **Card**: Flashcard with multiple types (basic, cloze, MCQ, image occlusion, audio, reverse)
- **CardSchedule**: Per-user scheduling state (new, learning, review, relearning, mastered)
- **ReviewLog**: Review history for analytics
- **StudySession**: Aggregated session metrics

### SRS Algorithms
- **SM-2**: Default classic algorithm (in `internal/srs/`)
- **FSRS**: Modern evidence-based scheduler
- Configurable per-deck with user defaults

## Workflow

1. **Create a new branch**: Create a new git branch for this work (e.g., `feature/TICKET-123-description`)

2. **Receive the Jira ticket**: Ask the user to provide the Jira ticket details (ticket ID, description, tasks, acceptance criteria)

3. **Analyze the ticket**: Review the requirements and identify:
   - Which handlers need to be created or modified
   - Required service layer changes
   - Database schema or query changes
   - New migrations needed

4. **Explore the codebase**: Before implementing, explore the existing back-end code to:
   - Understand current patterns and conventions
   - Find similar handlers/services to reference
   - Review existing repository patterns
   - Check model definitions

5. **Plan the implementation**: Create a task list based on the ticket tasks, mapping each to specific files and changes

6. **Implement the changes**:
   - Follow the handler → service → repository pattern
   - Write idiomatic Go code
   - Include proper error handling
   - Add input validation in handlers
   - Write SQL migrations for schema changes

7. **Verify acceptance criteria**: Check each acceptance criterion from the ticket

8. **Stop for review**: Once all changes are complete, stop and present the changes to the user for review

**IMPORTANT**: Do NOT commit, push, or create pull requests. Only create the branch and make code changes, then stop for user review.

## Rules

1. **Always explore first**: Before writing code, use the Explore agent to understand the current state of the back-end codebase

2. **Layer separation**: Keep handlers thin, business logic in services, data access in repositories

3. **Error handling**: Return appropriate HTTP status codes and structured error responses

4. **Database migrations**: Create migration files for any schema changes (up and down migrations)

5. **Input validation**: Validate all incoming request data in handlers before processing

6. **Security**:
   - Use parameterized queries (no SQL injection)
   - Hash passwords with bcrypt
   - Validate JWT tokens on protected routes
   - Apply rate limiting where appropriate

7. **Testing**: Write unit tests for services and integration tests for handlers

8. **Logging**: Use structured logging (slog) for observability

9. **Code style**: Follow Go conventions (gofmt, effective Go patterns)

## Example Interaction

**User**: /back-features-changes

**Assistant**: I'll help you work on a back-end feature or change. Please provide the Jira ticket details including:
- Ticket ID and title
- Description
- Tasks
- Acceptance criteria

Once you share the ticket, I'll analyze the requirements, explore the relevant parts of the codebase, and implement the changes following the project's layered architecture and Go conventions.
