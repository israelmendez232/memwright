# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Memwright is an open-source, web-first spaced repetition learning platform. Self-hostable alternative to Anki with a modern stack, clean dark UI, and extensible algorithm support.

## Build Commands

Commands will be available via Makefile once implemented:

```bash
# Full application (Docker Compose)
docker-compose up -d          # Start all services
docker-compose down           # Stop all services

# Frontend (web/)
cd web && npm install         # Install dependencies
cd web && npm run dev         # Development server
cd web && npm run build       # Production build
cd web && npm test            # Run tests

# Backend (api/)
cd api && go build ./cmd/server  # Build server binary
cd api && go run ./cmd/server    # Run server
cd api && go test ./...          # Run all tests
cd api && go test ./internal/srs # Run specific package tests

# Infrastructure (infra/)
cd infra && terraform init
cd infra && terraform plan
cd infra && terraform apply
```

## Architecture

### Stack
- **Frontend**: TypeScript, React, Tailwind CSS, shadcn/ui
- **Backend**: Go (net/http or Chi router), PostgreSQL
- **Infrastructure**: Terraform (AWS S3 + DynamoDB for state)
- **Deployment**: Docker Compose on single VPS

### Monorepo Structure
```
/web        - React frontend
/api        - Go backend
/infra      - Terraform infrastructure
/docs       - Documentation
```

### Backend Layout (api/)
- `cmd/server/` - Application entrypoint
- `internal/handler/` - HTTP handlers
- `internal/service/` - Business logic
- `internal/repository/` - Database queries
- `internal/model/` - Domain types
- `internal/srs/` - Spaced repetition algorithms (SM-2, FSRS)
- `migrations/` - SQL migrations

### Frontend Layout (web/src/)
- `components/` - React components
- `pages/` - Page components
- `hooks/` - Custom React hooks
- `lib/` - Utilities and SRS algorithm previews
- `styles/` - CSS/Tailwind styles

## Key Domain Concepts

### Core Entities
- **User** - account with global SRS settings
- **Deck** - collection of cards with tree hierarchy (subdecks)
- **Card** - flashcard with types: basic, cloze, MCQ, image occlusion, audio, reverse
- **CardSchedule** - per-user scheduling state (new → learning → review → relearning → mastered)
- **ReviewLog** - history of reviews for analytics
- **StudySession** - aggregated session metrics

### SRS Algorithms
- **SM-2** - default, classic Anki-style algorithm
- **FSRS** - modern evidence-based scheduler
- Algorithms are configurable per-deck with global user defaults

### Review Rating System
Three-button rating: Wrong (reschedule sooner), Correct (normal interval), Easy (extended interval)

## API Patterns

REST API at `/api/v1/`:
- Auth: `/auth/{register,login,refresh,logout}`
- Decks: `/decks`, `/decks/:id`, `/decks/:id/subdecks`
- Cards: `/decks/:deckId/cards`, `/cards/:id`
- Reviews: `/decks/:deckId/review/next`, `/cards/:id/review`
- Stats: `/stats/global`, `/stats/decks/:deckId`, `/stats/heatmap`

## Design Principles

- **Cheap to run**: target < $10/mo for personal instance
- **Self-host first**: Docker Compose as primary deployment
- **Single VPS**: Caddy (reverse proxy + TLS) → Go API → PostgreSQL
