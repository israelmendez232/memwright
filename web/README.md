# Memwright Web

React frontend for the Memwright spaced repetition learning platform.

## Prerequisites

- Docker and Docker Compose
- Make

## Development

All commands are run from the project root directory using Make.

Start all services (API, web, database) in development mode:

```bash
make dev
```

Build and start all services:

```bash
make dev-build
```

Stop all services:

```bash
make dev-down
```

View logs:

```bash
make dev-logs
```

Open a shell in the web container:

```bash
make web-shell
```

## Building for Production

Build web assets:

```bash
make build-web
```

The output will be in the `dist/` directory.

## Testing

Run web tests:

```bash
make test-web
```

Run all tests (API + web):

```bash
make test
```

## Local Development (without Docker)

If you prefer to run the web app outside Docker:

```bash
cd web
npm install
npm run dev          # Development server at http://localhost:5173
npm test             # Run tests once
npm run test:watch   # Run tests in watch mode
npm run lint         # Type checking
```

## Project Structure

```
web/
├── src/
│   ├── components/    # Reusable React components
│   ├── hooks/         # Custom React hooks
│   ├── lib/           # Utilities and helpers
│   ├── pages/         # Page components
│   ├── styles/        # CSS and Tailwind styles
│   ├── App.tsx        # Root application component
│   └── main.tsx       # Application entry point
├── tests/
│   ├── integration/   # Integration tests
│   ├── unit/          # Unit tests
│   └── setup.ts       # Test setup configuration
├── index.html         # HTML entry point
├── vite.config.ts     # Vite configuration
├── vitest.config.ts   # Vitest test configuration
└── tsconfig.json      # TypeScript configuration
```

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Vitest** - Test runner
- **Testing Library** - React testing utilities
