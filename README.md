# memwright

[![Frontend CI](https://github.com/israelmendez232/memwright/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/israelmendez232/memwright/actions/workflows/frontend-ci.yml)
[![Backend CI](https://github.com/israelmendez232/memwright/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/israelmendez232/memwright/actions/workflows/backend-ci.yml)

An open-source, web-first spaced repetition learning platform. Self-hostable alternative to Anki with a modern stack, clean dark UI, and extensible algorithm support.

> **WARNING: PERSONAL EXPERIMENT**
>
> This project is a personal learning experiment and is currently under active development.
>
> **DO NOT USE IN PRODUCTION.** The codebase is incomplete, APIs may change without notice, and there are no stability or security guarantees.

## Project Overview

An open-source, web-first spaced repetition learning platform. Self-hostable alternative to Anki with a modern stack, clean dark UI, and extensible algorithm support.

---

## Architecture

### Stack

|Layer|Technology|
|---|---|
|Front-end|TypeScript, React, Tailwind CSS, shadcn/ui|
|Back-end|Go (net/http or Chi router), PostgreSQL|
|Infrastructure|Terraform (state in AWS S3 + DynamoDB lock)|
|CI/CD|GitHub Actions|
|Deployment|Single VPS (e.g., Hetzner, DigitalOcean) via Docker Compose|

### Design Principles

- **Cheap to run**: single server, no managed services beyond S3 for state; target < $10/mo for a personal instance.
- **Self-host first**: Docker Compose as the primary deployment; all features work without external services.
- **Open source**: permissive license (MIT or Apache 2.0); community contributions welcome.
- **Monorepo**: one repository, distinct folders (`/web`, `/api`, `/infra`, `/docs`).

### Repository Structure

```
/
├── web/                  # Front-end (React + TypeScript)
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── lib/          # SRS algorithm implementations (client-side preview)
│   │   └── styles/
│   └── package.json
├── api/                  # Back-end (Go)
│   ├── cmd/server/       # Entrypoint
│   ├── internal/
│   │   ├── handler/      # HTTP handlers
│   │   ├── service/      # Business logic
│   │   ├── repository/   # DB queries
│   │   ├── model/        # Domain types
│   │   └── srs/          # Algorithm implementations (SM-2, FSRS, etc.)
│   ├── migrations/       # SQL migrations
│   └── go.mod
├── infra/                # Terraform
│   ├── modules/
│   └── environments/
├── docs/                 # Project documentation
├── docker-compose.yml
├── Makefile
└── .github/workflows/    # CI/CD
```

### Dev Environment (tmux sessions)

|Session|Panes|
|---|---|
|Front-end|Claude Code · NeoVim + terminal · LazyGit|
|Back-end + DB|Claude Code · NeoVim + terminal · LazyGit|
|Infrastructure|Claude Code · NeoVim + terminal · LazyGit|
|Project Mgmt|Claude Code (Jira MCP plugin)|

---

## Data Model (High Level)

### Core Entities

```
User
├── id, email, display_name, password_hash
├── timezone, locale, theme_preference
└── global_srs_settings (JSON)

Deck
├── id, user_id, parent_deck_id (nullable → subdecks)
├── title, description, cover_color, cover_image_url
├── tags[], privacy (public/private/shared)
├── srs_algorithm (sm2 | fsrs | custom), srs_config (JSON)
├── daily_new_limit, daily_review_limit
├── archived, position (ordering)
└── created_at, updated_at

Card
├── id, deck_id, card_type (basic | cloze | mcq | image_occlusion | audio | reverse)
├── front_content (rich text / JSON), back_content
├── tags[], custom_fields (JSONB)
├── media_attachments[]
├── suspended, flagged
└── created_at, updated_at

ReviewLog
├── id, card_id, user_id
├── rating (wrong | correct | easy)
├── interval_before, interval_after
├── ease_factor_before, ease_factor_after
├── review_duration_ms
└── reviewed_at

CardSchedule
├── card_id, user_id
├── state (new | learning | review | relearning | mastered)
├── due_at, interval_days, ease_factor
├── lapses, reps
└── updated_at

StudySession
├── id, user_id, deck_id
├── cards_reviewed, correct_count, wrong_count, easy_count
├── duration_seconds, streak_maintained
└── started_at, ended_at
```

### Relationships

- User → many Decks → many Cards
- Deck → optional parent Deck (tree structure for subdecks)
- Card → many ReviewLogs
- Card → one CardSchedule per user
- User → many StudySessions

---

## Features

### 1. Web Interface

- **Mobile-first responsive** layout
- **PWA** with service worker for offline review
- **Dark theme** default, light theme toggle
- Touch swipe gestures for card review on mobile
- Keyboard shortcuts on desktop (1/2/3 for rating, space to flip, e to edit)
- WCAG 2.1 AA accessibility

### 2. Deck Management

- CRUD + archive for decks
- **Subdeck hierarchy**: decks contain subdecks as a tree; rendered as collapsible lists
- Metadata: title, description, tags (hierarchical), author, difficulty level, category, custom cover color/image
- Organization: favorites, custom sort (alpha, date, frequency), deck templates
- Sharing: share via link, public deck marketplace (future), collaborative editing
- Import/export: JSON native format, CSV, Anki `.apkg` compatibility

### 3. Flashcard System

- **Rich text editor** (Tiptap or Lexical) for card content
- Card types: basic, cloze deletion, multiple choice, image occlusion, audio, reverse (auto-generated)
- Media: image upload + URL, audio files, LaTeX rendering (KaTeX), code highlighting (Shiki), embedded video
- Per-card metadata: tags, difficulty, source URL, custom key-value fields
- Bulk operations: multi-select, bulk tag, bulk move, bulk delete

### 4. Review Experience

- **Three-button rating**:
    - 🔴 Wrong — reschedule sooner (lapse)
    - 🟢 Correct — normal next interval
    - 🔵 Easy — extended interval / optional suspend
- Card flip animation, button press feedback, streak counter
- Session controls: undo last review, edit card inline, flag card, hint reveal (progressive), auto-play audio
- Real-time session stats: cards done, accuracy, time, remaining

### 5. Spaced Repetition Algorithms

- **SM-2** (default, classic Anki-style)
- **FSRS** (Free Spaced Repetition Scheduler — modern, evidence-based)
- **Custom algorithm support**: user-defined via config or plugin interface (future)
- Configurable per-deck: algorithm choice, parameters, new card rate, daily limits, graduating intervals, easy bonus, lapse handling
- Global defaults in user settings, overridable per deck

### 6. Dashboard & Analytics

- **Global dashboard page** (dedicated route, not a modal):
    - Heatmap calendar (GitHub contribution style)
    - Study streak tracker with daily/weekly/monthly goals
    - Total stats: cards, decks, time studied, retention rate
    - Due forecast: today, tomorrow, next 7 days
- **Per-deck analytics**:
    - Card state distribution (new / learning / review / mastered)
    - Retention rate over time (line chart)
    - Most difficult cards ranked
    - Time spent per deck
- **Filtering**: by deck, tags (AND/OR), card state, date range, difficulty, media type; saved filters
- Export: CSV data, PDF report (future)

### 7. Gamification

- Daily streak tracking with visual flame/counter
- XP system: points per review, bonuses for streaks and accuracy
- Achievement badges (e.g., "100 day streak", "1000 cards mastered", "Night Owl")
- Leaderboards (opt-in, for shared/public instances)
- Level progression with milestones
- Weekly challenges (e.g., "review 50 cards daily for 7 days")

### 8. User Settings & Configuration

- **Account**: profile, email, notification preferences (email/push reminders), timezone, locale, theme, data export, account deletion
- **Global study settings**: default SRS algorithm + parameters, session duration, break reminders, review priorities, card burial (hide siblings)
- **Per-deck overrides**: SRS params, daily limits, review order (random / due date / added date), time limits, auto-advance, custom theme/color

---

## API Design (REST)

### Auth

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
```

### Decks

```
GET    /api/v1/decks
POST   /api/v1/decks
GET    /api/v1/decks/:id
PUT    /api/v1/decks/:id
DELETE /api/v1/decks/:id
GET    /api/v1/decks/:id/subdecks
POST   /api/v1/decks/:id/export
POST   /api/v1/decks/import
```

### Cards

```
GET    /api/v1/decks/:deckId/cards
POST   /api/v1/decks/:deckId/cards
GET    /api/v1/cards/:id
PUT    /api/v1/cards/:id
DELETE /api/v1/cards/:id
POST   /api/v1/cards/bulk          # bulk operations
```

### Reviews

```
GET    /api/v1/decks/:deckId/review/next    # get next batch of due cards
POST   /api/v1/cards/:id/review             # submit a review rating
POST   /api/v1/cards/:id/review/undo
```

### Stats & Dashboard

```
GET    /api/v1/stats/global
GET    /api/v1/stats/decks/:deckId
GET    /api/v1/stats/heatmap?range=year
GET    /api/v1/stats/forecast
```

### Settings

```
GET    /api/v1/settings
PUT    /api/v1/settings
GET    /api/v1/decks/:deckId/settings
PUT    /api/v1/decks/:deckId/settings
```

---

## Infrastructure

### Target: Single VPS Deployment

```
┌─────────────────────────────────────┐
│             VPS (2GB RAM)           │
│  ┌───────────┐  ┌───────────────┐  │
│  │  Caddy    │→ │  Go API       │  │
│  │  (reverse │  │  :8080        │  │
│  │   proxy + │  └───────────────┘  │
│  │   TLS)    │  ┌───────────────┐  │
│  │           │→ │  Static files  │  │
│  └───────────┘  │  (web build)  │  │
│                 └───────────────┘  │
│  ┌───────────────┐                 │
│  │  PostgreSQL    │                 │
│  │  :5432         │                 │
│  └───────────────┘                 │
│  ┌───────────────┐                 │
│  │  S3-compatible │ (optional,     │
│  │  MinIO         │  for media)    │
│  └───────────────┘                 │
└─────────────────────────────────────┘
```

- **Caddy**: auto-TLS via Let's Encrypt, reverse proxy to Go API, serves static SPA build
- **Go binary**: single compiled binary, no runtime dependencies
- **PostgreSQL**: single instance, daily pg_dump backups to S3/local
- **MinIO** (optional): S3-compatible object storage for media uploads; can swap for local filesystem

### Terraform

- AWS S3 bucket + DynamoDB table for Terraform remote state
- VPS provisioning (Hetzner/DO provider)
- DNS records
- Backup bucket
- Modular: `infra/modules/` for reusable components

### CI/CD (GitHub Actions)

```
on push to main:
  1. Lint + test (Go + TS in parallel)
  2. Build Go binary
  3. Build web static assets
  4. Build Docker image → push to GHCR
  5. Deploy to VPS (SSH + docker compose pull + up)

on pull request:
  1. Lint + test
  2. Build check (no deploy)
```

---

## Non-Functional Requirements

|Concern|Target|
|---|---|
|Performance|< 200ms API response, < 3s initial load (LCP)|
|Offline|PWA service worker caches review queue; syncs when online|
|Hosting cost|< $10/mo for personal instance|
|Data ownership|All data stored in user's PostgreSQL; full export anytime|
|Security|JWT auth (access + refresh tokens), bcrypt passwords, HTTPS enforced, rate limiting|
|Scalability|Single user → small team. For larger scale: add Redis cache, read replicas (future)|
|Backup|Automated daily pg_dump; configurable retention|
|Observability|Structured logging (slog), Prometheus metrics endpoint (optional)|

---

## Implementation Phases

### Phase 1 — Foundation

- Project scaffolding (monorepo, Docker Compose, Makefile)
- Auth (register, login, JWT)
- Deck CRUD with subdeck hierarchy
- Basic card CRUD (front/back only)
- SM-2 algorithm implementation
- Review flow with three-button rating
- Dark theme UI shell

### Phase 2 — Core Experience

- Rich text editor for cards
- Additional card types (cloze, reverse, MCQ)
- Media upload (images, audio)
- Per-deck SRS configuration
- Dashboard with heatmap and basic stats
- Keyboard shortcuts
- Undo review

### Phase 3 — Polish & Features

- PWA + offline support
- FSRS algorithm
- Import/export (JSON, CSV, .apkg)
- Gamification (streaks, XP, badges)
- Advanced analytics and filtering
- Deck sharing via link
- LaTeX and code highlighting

### Phase 4 — Scale & Community

- Public deck marketplace
- Collaborative editing
- Custom algorithm plugin system
- Leaderboards
- Mobile-optimized gesture navigation
- Notification system (email/push reminders)
- Terraform modules for one-click deploy

