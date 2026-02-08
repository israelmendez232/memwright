---
name: back-features-changes
description: Work on back-end Go features and API changes based on Jira tickets
args: ticket-id
---

# Back-End Features and Changes Skill

When the user invokes `/back-features-changes --ticket-id <ID>`, fetch the Jira ticket details via MCP and work on the back-end implementation.

## Input

- `--ticket-id`: The Jira ticket ID (e.g., `MEM-1` or just `1`). If only a number is provided, prefix with `MEM-`.

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

1. **Fetch the Jira ticket**: Use the Atlassian MCP tool to fetch the ticket details:
   - Call `mcp__atlassian__getJiraIssue` with `cloudId: "45d5c0c7-c3c3-468a-a8fb-7f19720b0424"` and `issueIdOrKey: "<ticket-id>"`
   - Extract: title (summary), description, tasks, acceptance criteria
   - If the ticket cannot be found, inform the user and stop

2. **ALWAYS sync with main branch first**: Before creating a new branch, you MUST ensure you're on main with the latest code:
   - Run: `git checkout main`
   - Run: `git pull`
   - **NEVER create a branch from another feature branch**

3. **Create a feature branch**: Based on the ticket title, create a branch name:
   - Format: `<ticket-id-lowercase>-<title-slug>`
   - Convert the title to lowercase, replace spaces with hyphens, remove special characters
   - Example: Ticket `MEM-1` with title "Project Scaffolding and Docker Setup" → branch `mem-1-project-scaffolding-and-docker-setup`
   - Run: `git checkout -b <branch-name>`

4. **Analyze the ticket**: Review the requirements and identify:
   - Which handlers need to be created or modified
   - Required service layer changes
   - Database schema or query changes
   - New migrations needed

5. **Explore the codebase**: Before implementing, explore the existing back-end code to:
   - Understand current patterns and conventions
   - Find similar handlers/services to reference
   - Review existing repository patterns
   - Check model definitions

6. **Plan the implementation**: Create a task list based on the ticket tasks, mapping each to specific files and changes

7. **Implement the changes**:
   - Follow the handler → service → repository pattern
   - Write idiomatic Go code
   - Include proper error handling
   - Add input validation in handlers
   - Write SQL migrations for schema changes

8. **Verify acceptance criteria**: Check each acceptance criterion from the ticket

9. **Stop for review**: Present a summary of changes to the user

**IMPORTANT**: Do NOT run any git commands after making code changes. No `git add`, no `git commit`, no `git push`, no staging, no pull requests. Only create the branch at the start, then make code changes and stop.

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

10. **Go standards and best practices**:
    - Follow idiomatic Go patterns (accept interfaces, return structs)
    - Use meaningful variable names that are concise but descriptive
    - Prefer short variable names in small scopes (e.g., `r` for request, `w` for writer)
    - Handle errors explicitly, don't ignore them
    - Use `context.Context` for cancellation and timeouts
    - Avoid global state, prefer dependency injection

11. **Go module path**:
    - Use a platform-agnostic module path: `memwright/api`
    - Do NOT hardcode hosting platforms in the module path (e.g., avoid `github.com/...` or `gitlab.com/...`)
    - All internal imports should use the module path as prefix: `memwright/api/internal/...`
    - This ensures portability if the repository moves between hosting providers

12. **Avoid unnecessary comments**:
    - Don't add comments that just restate what the code does
    - Only add comments for non-obvious business logic or complex algorithms
    - Let the code be self-documenting through clear naming
    - No redundant godoc comments on unexported functions

## Example Interaction

**User**: /back-features-changes --ticket-id MEM-1

**Assistant**:
1. Fetches ticket MEM-1 from Jira via MCP
2. Displays the ticket summary:
   - Title: "Project Scaffolding and Docker Setup"
   - Description: Initialize Go project structure...
   - Tasks: Set up directory structure, Configure Docker Compose...
3. Syncs with main: `git checkout main && git pull`
4. Creates branch: `git checkout -b mem-1-project-scaffolding-and-docker-setup`
5. Explores the codebase and plans implementation
6. Implements changes following the layered architecture
7. Stops for user review
