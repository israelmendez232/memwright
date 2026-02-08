# Memwright MVP - Jira Tickets

Total: 19 tickets | 152 story points

---

# [MEM-001] Project Scaffolding and Docker Setup

## Priority
High

## Story Points
8

## Description

The memwright project needs initial scaffolding before any feature development can begin. Currently, the repository only contains documentation (README.md, CLAUDE.md) with no actual source code. This ticket establishes the foundational project structure including the monorepo layout, build tooling, and local development environment.

The implementation will create the `/api` and `/web` directories with their respective build configurations. Docker Compose will orchestrate PostgreSQL, the Go API, and the React frontend for local development. A Makefile will provide standardized commands for common development tasks. This work unblocks all subsequent development tickets.

## Tasks

- [ ] Initialize Go module in /api with go.mod (1 pt)
- [ ] Create cmd/server/main.go with basic HTTP server and health endpoint (2 pts)
- [ ] Create internal/config/config.go for environment variable loading (1 pt)
- [ ] Initialize Vite React TypeScript project in /web (1 pt)
- [ ] Create docker-compose.yml with PostgreSQL, API, and web services (2 pts)
- [ ] Create Makefile with dev, test, build, and migrate targets (1 pt)

## Acceptance Criteria

- [ ] Running `make dev` starts all three services (db, api, web)
- [ ] API health endpoint responds at http://localhost:8080/health
- [ ] Web dev server runs at http://localhost:5173
- [ ] PostgreSQL is accessible at localhost:5432
- [ ] Go module compiles without errors
- [ ] React app renders without errors

---

# [MEM-002] Database Schema and Migrations

## Priority
High

## Story Points
10

## Description

The application requires a PostgreSQL database schema to persist users, decks, cards, and review data. This ticket implements the database layer including connection pooling, migration infrastructure, and all Phase 1 table definitions. The schema design supports future features like deck sharing through composite keys on card_schedules.

The implementation uses golang-migrate for version-controlled migrations and pgx v5 for high-performance PostgreSQL connectivity. Each entity (users, decks, cards, card_schedules, review_logs) gets separate up/down migration files. The schema includes proper indexes for common query patterns and triggers for automatic updated_at timestamps.

## Tasks

- [ ] Create internal/database/postgres.go with connection pool setup (2 pts)
- [ ] Set up golang-migrate infrastructure (1 pt)
- [ ] Create migration 001: users table with indexes (1 pt)
- [ ] Create migration 002: decks table with self-referential parent_deck_id (2 pts)
- [ ] Create migration 003: cards table (1 pt)
- [ ] Create migration 004: card_schedules table with composite PK (1 pt)
- [ ] Create migration 005: review_logs table with indexes (1 pt)
- [ ] Create updated_at trigger function and apply to all tables (1 pt)

## Acceptance Criteria

- [ ] Running `make migrate` applies all migrations successfully
- [ ] Running `make migrate-down` rolls back migrations correctly
- [ ] All foreign key constraints are properly defined
- [ ] Indexes exist for user_id, deck_id, due_at, and reviewed_at columns
- [ ] UUID extension is enabled and used for primary keys
- [ ] Decks table supports unlimited subdeck nesting via parent_deck_id

---

# [MEM-003] Backend Domain Models and Repository Layer

## Priority
High

## Story Points
12

## Description

The backend requires domain models and a repository layer to interact with the database. This ticket creates Go structs for all entities and implements CRUD operations using the repository pattern. The repository layer abstracts database queries from business logic, making the codebase testable and maintainable.

The implementation follows idiomatic Go patterns with interfaces for testability. Each entity gets its own repository file with standard operations (Create, GetByID, List, Update, Delete). The deck repository includes hierarchy-aware queries for building subdeck trees. All repositories use pgx for prepared statements and proper error handling.

## Tasks

- [ ] Create internal/model/ with User, Deck, Card, CardSchedule, ReviewLog structs (2 pts)
- [ ] Create internal/model/errors.go with domain error types (1 pt)
- [ ] Implement internal/repository/user.go with CRUD operations (2 pts)
- [ ] Implement internal/repository/deck.go with hierarchy queries (3 pts)
- [ ] Implement internal/repository/card.go with CRUD operations (2 pts)
- [ ] Implement internal/repository/card_schedule.go with due card queries (2 pts)

## Acceptance Criteria

- [ ] All model structs have proper JSON and database tags
- [ ] Repository methods return domain errors (ErrNotFound, ErrDuplicateEmail)
- [ ] Deck repository can fetch all subdecks for a parent deck
- [ ] Card schedule repository can query cards due before a given time
- [ ] All queries use parameterized statements to prevent SQL injection
- [ ] Repository interface is defined for dependency injection

---

# [MEM-004] SM-2 Spaced Repetition Algorithm

## Priority
High

## Story Points
6

## Description

The SM-2 algorithm is the core scheduling engine that determines when cards should be reviewed. This ticket implements the classic SuperMemo 2 algorithm that adjusts card intervals based on user ratings (Wrong, Correct, Easy). The algorithm calculates new due dates, updates ease factors, and tracks learning state transitions.

The implementation follows the original SM-2 specification with configurable parameters. Cards transition through states (new, learning, review, relearning, mastered) based on performance. Wrong ratings reset intervals and decrease ease factor. Correct ratings apply the current ease factor. Easy ratings provide bonus interval multipliers. Comprehensive unit tests validate all edge cases.

## Tasks

- [ ] Create internal/srs/algorithm.go with Algorithm interface (1 pt)
- [ ] Implement SM-2 core logic in internal/srs/sm2.go (3 pts)
- [ ] Write comprehensive unit tests in internal/srs/sm2_test.go (2 pts)

## Acceptance Criteria

- [ ] Wrong rating resets interval to 1 day and decreases ease factor by 0.2
- [ ] Correct rating multiplies interval by ease factor
- [ ] Easy rating applies 1.3x bonus multiplier and increases ease factor
- [ ] Ease factor never drops below 1.3
- [ ] New cards graduate to learning state after first correct rating
- [ ] All state transitions are correctly implemented
- [ ] Unit tests cover all rating scenarios and edge cases

---

# [MEM-005] Authentication Service and JWT Handling

## Priority
High

## Story Points
10

## Description

User authentication is required before any personalized features can function. This ticket implements the authentication service with user registration, login, and JWT token management. The system uses short-lived access tokens paired with longer-lived refresh tokens for security.

The implementation uses bcrypt for password hashing and golang-jwt for token generation. Access tokens expire after 15 minutes and are stored in memory on the client. Refresh tokens expire after 7 days and are sent as httpOnly cookies. The service includes input validation, duplicate email detection, and proper error responses.

## Tasks

- [ ] Create internal/service/auth.go with Register method (2 pts)
- [ ] Implement Login method with credential verification (2 pts)
- [ ] Implement JWT access token generation with claims (2 pts)
- [ ] Implement refresh token generation and validation (2 pts)
- [ ] Create internal/dto/auth.go with request/response types (1 pt)
- [ ] Add password strength validation (1 pt)

## Acceptance Criteria

- [ ] Passwords are hashed with bcrypt cost factor of 12
- [ ] Access tokens contain user ID and expiration claims
- [ ] Refresh tokens are cryptographically random and stored securely
- [ ] Duplicate email registration returns appropriate error
- [ ] Invalid credentials return 401 Unauthorized
- [ ] Token validation rejects expired or malformed tokens

---

# [MEM-006] Backend Middleware Layer

## Priority
High

## Story Points
6

## Description

The API requires middleware for cross-cutting concerns including authentication, CORS, logging, and panic recovery. This ticket implements the middleware stack that processes every request before reaching handlers. Proper middleware ensures security, observability, and reliability.

The implementation uses Chi's middleware pattern for composability. The auth middleware extracts and validates JWT tokens, injecting user context for downstream handlers. CORS middleware enables frontend communication. Logging middleware records request details using structured logging (slog). Recovery middleware catches panics and returns proper error responses.

## Tasks

- [ ] Create internal/middleware/auth.go for JWT validation and context injection (2 pts)
- [ ] Create internal/middleware/cors.go with configurable origins (1 pt)
- [ ] Create internal/middleware/logging.go with structured request logging (2 pts)
- [ ] Create internal/middleware/recovery.go for panic handling (1 pt)

## Acceptance Criteria

- [ ] Auth middleware extracts token from Authorization Bearer header
- [ ] Protected routes return 401 when token is missing or invalid
- [ ] User ID is available in request context after auth middleware
- [ ] CORS allows configured origins and required methods
- [ ] All requests are logged with method, path, status, and duration
- [ ] Panics are caught and return 500 Internal Server Error

---

# [MEM-007] Backend HTTP Handlers - Auth Endpoints

## Priority
High

## Story Points
6

## Description

The authentication API endpoints expose registration, login, token refresh, and logout functionality to clients. This ticket implements the HTTP handlers that parse requests, call the auth service, and format responses. These endpoints are public (no auth required) except for logout.

The implementation follows RESTful conventions with proper HTTP status codes. Request bodies are validated before processing. Responses use consistent JSON structure with error details when applicable. The refresh endpoint reads the httpOnly cookie and issues new token pairs.

## Tasks

- [ ] Create internal/handler/handler.go with dependency injection struct (1 pt)
- [ ] Implement POST /auth/register handler (2 pts)
- [ ] Implement POST /auth/login handler (1 pt)
- [ ] Implement POST /auth/refresh handler (1 pt)
- [ ] Implement POST /auth/logout handler (1 pt)

## Acceptance Criteria

- [ ] Register returns 201 Created with user data (excluding password)
- [ ] Login returns 200 OK with access token in body and refresh token in cookie
- [ ] Refresh returns new token pair when valid refresh token provided
- [ ] Logout clears the refresh token cookie
- [ ] Validation errors return 400 Bad Request with field details
- [ ] All responses follow consistent JSON structure

---

# [MEM-008] Deck Service and HTTP Handlers

## Priority
High

## Story Points
10

## Description

Decks are the primary organizational unit for flashcards. This ticket implements the deck service layer with business logic and the HTTP handlers for deck CRUD operations. The service handles subdeck hierarchy, position ordering, and ownership validation.

The implementation includes methods for creating decks with optional parent assignment, listing user decks as a tree structure, updating deck properties, and cascading deletes. The handlers map HTTP requests to service calls with proper authorization checks ensuring users can only access their own decks.

## Tasks

- [ ] Create internal/service/deck.go with Create, List, GetByID methods (3 pts)
- [ ] Implement Update and Delete methods with ownership validation (2 pts)
- [ ] Implement GetSubdecks and tree building logic (2 pts)
- [ ] Create internal/handler/deck.go with all CRUD handlers (2 pts)
- [ ] Create internal/dto/deck.go with request/response types (1 pt)

## Acceptance Criteria

- [ ] Users can create top-level decks and subdecks
- [ ] List endpoint returns decks with nested subdeck structure
- [ ] Users cannot access or modify other users decks
- [ ] Deleting a deck cascades to subdecks and cards
- [ ] Position field allows custom ordering of decks
- [ ] Archived decks are excluded from default list queries

---

# [MEM-009] Card Service and HTTP Handlers

## Priority
High

## Story Points
8

## Description

Cards contain the front and back content that users review. This ticket implements card CRUD operations including automatic card schedule creation. When a card is created, a corresponding card_schedule record initializes the spaced repetition state for that user.

The implementation validates deck ownership before card operations and supports basic card type (front/back text). The service creates card schedules with default values (new state, initial ease factor 2.5). Handlers provide endpoints for listing cards within a deck and individual card management.

## Tasks

- [ ] Create internal/service/card.go with Create method and schedule initialization (2 pts)
- [ ] Implement List, GetByID, Update, Delete methods (2 pts)
- [ ] Create internal/handler/card.go with all CRUD handlers (2 pts)
- [ ] Create internal/dto/card.go with request/response types (1 pt)
- [ ] Add deck ownership validation for all card operations (1 pt)

## Acceptance Criteria

- [ ] Creating a card automatically creates a card_schedule in new state
- [ ] Cards can only be created in decks owned by the authenticated user
- [ ] List endpoint returns all non-suspended cards in a deck
- [ ] Updating card content does not affect scheduling state
- [ ] Deleting a card cascades to schedule and review logs
- [ ] Card responses include current schedule state and due date

---

# [MEM-010] Review Service and HTTP Handlers

## Priority
High

## Story Points
10

## Description

The review system orchestrates study sessions by selecting due cards and processing user ratings. This ticket implements the review service that queries cards ready for review, applies the SM-2 algorithm based on ratings, and logs review history. This is the core learning loop of the application.

The implementation fetches cards where due_at is in the past plus new cards up to the daily limit. When a rating is submitted, the service applies SM-2 to update the schedule, creates a review log entry, and returns the updated card state. The handler exposes endpoints for getting the next review batch and submitting ratings.

## Tasks

- [ ] Create internal/service/review.go with GetNextCards method (3 pts)
- [ ] Implement SubmitReview method with SM-2 application (3 pts)
- [ ] Create internal/repository/review_log.go for history tracking (1 pt)
- [ ] Create internal/handler/review.go with endpoints (2 pts)
- [ ] Create internal/dto/review.go with request/response types (1 pt)

## Acceptance Criteria

- [ ] GetNextCards returns due cards ordered by due_at ascending
- [ ] New cards are included up to deck daily_new_limit
- [ ] Review cards are included up to deck daily_review_limit
- [ ] SubmitReview updates card_schedule with new interval and due_at
- [ ] Each review creates a review_log entry with before/after values
- [ ] Review duration can be optionally tracked in milliseconds

---

# [MEM-011] Backend Router Setup and Integration

## Priority
Medium

## Story Points
4

## Description

All handlers and middleware need to be wired together into a cohesive router configuration. This ticket creates the main router setup that applies middleware in the correct order and maps all endpoints to their handlers. The router defines public vs protected route groups.

The implementation uses Chi router with middleware applied globally (CORS, logging, recovery) and selectively (auth for protected routes). All API routes are prefixed with /api/v1 for versioning. The router configuration is centralized for maintainability.

## Tasks

- [ ] Create router setup function with middleware chain (2 pts)
- [ ] Configure public routes (/auth/*) without auth middleware (1 pt)
- [ ] Configure protected routes with auth middleware (1 pt)

## Acceptance Criteria

- [ ] All routes are prefixed with /api/v1
- [ ] Auth endpoints are accessible without token
- [ ] All other endpoints require valid JWT token
- [ ] CORS headers are present on all responses
- [ ] Request logging works for all endpoints
- [ ] Health endpoint remains at /health (no prefix)

---

# [MEM-012] Frontend Scaffolding and UI Foundation

## Priority
High

## Story Points
10

## Description

The frontend requires initial setup with React, TypeScript, Tailwind CSS, and shadcn/ui components. This ticket establishes the build tooling, styling foundation, and core UI components needed for all pages. The dark theme is the default with carefully chosen color variables.

The implementation uses Vite for fast development builds, Tailwind for utility-first styling, and shadcn/ui for accessible, customizable components. The dark theme CSS variables follow shadcn/ui conventions. Layout components (header, sidebar, main layout) provide consistent structure across pages.

## Tasks

- [ ] Configure Tailwind CSS with dark theme as default (2 pts)
- [ ] Initialize shadcn/ui and install button, input, card, dialog, toast components (2 pts)
- [ ] Create src/styles/globals.css with dark theme color variables (1 pt)
- [ ] Create src/components/layout/header.tsx with navigation (2 pts)
- [ ] Create src/components/layout/main-layout.tsx wrapper (1 pt)
- [ ] Create src/components/layout/auth-layout.tsx for login/register (1 pt)
- [ ] Set up lucide-react for icons (1 pt)

## Acceptance Criteria

- [ ] Application renders with dark background and light text
- [ ] All shadcn/ui components use the dark theme variables
- [ ] Header displays app name and navigation links
- [ ] Main layout includes sidebar area for deck navigation
- [ ] Auth layout is centered with max-width container
- [ ] Icons render correctly throughout the application

---

# [MEM-013] Frontend Routing and Auth Context

## Priority
High

## Story Points
8

## Description

The frontend needs client-side routing and authentication state management. This ticket implements React Router for navigation and an Auth context provider that manages user session state. Protected routes redirect unauthenticated users to login.

The implementation uses React Router v6 with nested routes for layout sharing. The Auth context stores user data and tokens, provides login/logout/register functions, and persists session across page refreshes using localStorage. An Axios interceptor automatically refreshes expired tokens.

## Tasks

- [ ] Set up React Router with route definitions in App.tsx (2 pts)
- [ ] Create src/context/auth-context.tsx with AuthProvider (3 pts)
- [ ] Create src/components/auth/protected-route.tsx (1 pt)
- [ ] Create src/lib/api/client.ts with Axios instance and interceptors (2 pts)

## Acceptance Criteria

- [ ] Unauthenticated users are redirected to /login
- [ ] Authenticated users can access protected routes
- [ ] Auth state persists across page refreshes
- [ ] Expired access tokens trigger automatic refresh
- [ ] Failed refresh redirects to login
- [ ] Logout clears all auth state and redirects to login

---

# [MEM-014] Frontend API Layer and Types

## Priority
High

## Story Points
6

## Description

The frontend needs TypeScript types matching the backend DTOs and API functions for all endpoints. This ticket creates the type definitions and API client functions that components will use for data fetching. Consistent types ensure type safety across the application.

The implementation defines interfaces for User, Deck, Card, CardSchedule, and API responses. API functions wrap Axios calls with proper typing. React Query will use these functions for caching and synchronization. Error responses are typed for consistent error handling.

## Tasks

- [ ] Create src/types/ with user.ts, deck.ts, card.ts, review.ts (2 pts)
- [ ] Create src/lib/api/auth.ts with register, login, refresh, logout functions (1 pt)
- [ ] Create src/lib/api/decks.ts with CRUD functions (1 pt)
- [ ] Create src/lib/api/cards.ts with CRUD functions (1 pt)
- [ ] Create src/lib/api/reviews.ts with getNext and submitReview functions (1 pt)

## Acceptance Criteria

- [ ] All API response types match backend DTO structures
- [ ] API functions return properly typed promises
- [ ] Error responses include message and optional field errors
- [ ] All functions use the configured Axios client with interceptors
- [ ] TypeScript compilation succeeds with strict mode

---

# [MEM-015] Frontend Auth Pages

## Priority
High

## Story Points
6

## Description

Users need login and registration pages to access the application. This ticket implements the auth forms with validation, error handling, and loading states. The forms use shadcn/ui components and integrate with the auth context.

The implementation includes email and password fields with client-side validation. Registration adds a display name field. Form submission shows loading state and handles errors gracefully. Successful login redirects to the dashboard.

## Tasks

- [ ] Create src/components/auth/login-form.tsx with validation (2 pts)
- [ ] Create src/components/auth/register-form.tsx with validation (2 pts)
- [ ] Create src/pages/login.tsx using auth layout (1 pt)
- [ ] Create src/pages/register.tsx using auth layout (1 pt)

## Acceptance Criteria

- [ ] Login form validates email format and password presence
- [ ] Register form validates matching passwords and display name
- [ ] Form errors display inline next to fields
- [ ] API errors display as toast notifications
- [ ] Submit button shows loading spinner during request
- [ ] Successful auth redirects to dashboard

---

# [MEM-016] Frontend Dashboard and Deck Components

## Priority
High

## Story Points
12

## Description

The dashboard is the main landing page showing all user decks with their subdeck hierarchy. This ticket implements the dashboard page and deck-related components including the deck list, individual deck cards, and create/edit deck dialog. Users can see due card counts and navigate to individual decks.

The implementation uses React Query for data fetching with automatic refetching. The deck list renders as a tree structure with collapsible subdecks. Each deck card shows title, description preview, and due/new card counts. The deck form dialog handles both creation and editing.

## Tasks

- [ ] Create src/hooks/use-decks.ts with React Query hooks (2 pts)
- [ ] Create src/components/deck/deck-list.tsx container (2 pts)
- [ ] Create src/components/deck/deck-item.tsx with stats display (2 pts)
- [ ] Create src/components/deck/deck-tree.tsx for recursive rendering (2 pts)
- [ ] Create src/components/deck/deck-form.tsx dialog for create/edit (2 pts)
- [ ] Create src/pages/dashboard.tsx with deck list and create button (2 pts)

## Acceptance Criteria

- [ ] Dashboard displays all user decks in tree structure
- [ ] Subdecks are indented and collapsible under parent decks
- [ ] Each deck shows title, card count, and due count
- [ ] Clicking a deck navigates to deck detail page
- [ ] Create deck button opens form dialog
- [ ] New decks appear in list after creation without page refresh

---

# [MEM-017] Frontend Deck Detail and Card Management

## Priority
High

## Story Points
12

## Description

The deck detail page shows a single deck with its cards and subdecks. Users can add, edit, and delete cards from this page. A study button initiates the review session for due cards. This ticket implements the complete card management interface.

The implementation displays deck metadata in a header section with edit and delete options. Cards are listed in a table or grid with front content preview. Card creation and editing use a form with front and back text areas. Bulk selection enables multi-card operations.

## Tasks

- [ ] Create src/hooks/use-cards.ts with React Query hooks (2 pts)
- [ ] Create src/components/card/card-list.tsx table/grid (2 pts)
- [ ] Create src/components/card/card-item.tsx preview component (1 pt)
- [ ] Create src/components/card/card-form.tsx for create/edit (2 pts)
- [ ] Create src/components/card/card-editor.tsx with front/back inputs (2 pts)
- [ ] Create src/pages/deck-detail.tsx with cards and study button (3 pts)

## Acceptance Criteria

- [ ] Deck detail shows title, description, and deck settings
- [ ] Cards list displays front content preview for each card
- [ ] Add card button opens form with front and back editors
- [ ] Edit card opens form with existing content
- [ ] Delete card shows confirmation dialog
- [ ] Study button navigates to review page (only enabled if cards due)

---

# [MEM-018] Frontend Review Session

## Priority
High

## Story Points
12

## Description

The review session is the core learning experience where users study flashcards. This ticket implements the review page with card flip animation, rating buttons, and session progress tracking. The interface is clean and focused to minimize distractions during study.

The implementation fetches due cards on mount and presents them one at a time. Clicking the card or pressing space flips to reveal the back. Rating buttons (Wrong, Correct, Easy) submit the review and advance to the next card. Session stats show cards reviewed, accuracy, and remaining count. Keyboard shortcuts enhance desktop experience.

## Tasks

- [ ] Create src/hooks/use-review.ts with session state management (3 pts)
- [ ] Create src/components/review/flashcard.tsx with flip animation (3 pts)
- [ ] Create src/components/review/rating-buttons.tsx with Wrong/Correct/Easy (2 pts)
- [ ] Create src/components/review/session-stats.tsx progress display (1 pt)
- [ ] Create src/pages/review.tsx combining all components (2 pts)
- [ ] Add keyboard shortcuts (space to flip, 1/2/3 for ratings) (1 pt)

## Acceptance Criteria

- [ ] Card displays front content initially
- [ ] Clicking card or pressing space reveals back content with animation
- [ ] Rating buttons appear only after card is flipped
- [ ] Submitting rating advances to next card
- [ ] Session stats update after each review
- [ ] Session complete screen shows when all cards reviewed
- [ ] Keyboard shortcuts work on desktop browsers

---

# [MEM-019] Frontend-Backend Integration Testing

## Priority
Medium

## Story Points
6

## Description

After all features are implemented, the complete system needs end-to-end verification. This ticket covers manual integration testing of all user flows and fixing any issues discovered. The goal is to ensure frontend and backend work together seamlessly.

The testing covers the complete user journey: registration, login, deck creation, card creation, review session, and logout. Each API endpoint is exercised through the UI. Error scenarios are tested including invalid inputs, network failures, and expired tokens.

## Tasks

- [ ] Test complete auth flow (register, login, refresh, logout) (1 pt)
- [ ] Test deck CRUD with subdeck hierarchy (1 pt)
- [ ] Test card CRUD and deck ownership validation (1 pt)
- [ ] Test review session with all rating types (2 pts)
- [ ] Fix any integration issues discovered (1 pt)

## Acceptance Criteria

- [ ] New users can register and immediately log in
- [ ] Decks and subdecks display correctly after creation
- [ ] Cards appear in their assigned decks
- [ ] Review session correctly schedules cards based on ratings
- [ ] Session persists across page refreshes
- [ ] All API error responses display user-friendly messages
