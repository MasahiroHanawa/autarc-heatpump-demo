# Heat Pump Planning Assistant

A fullstack demo inspired by heat pump installer workflows. Installers enter building data and site-survey notes; the system processes them asynchronously and returns structured recommendations.

## Why this project exists

- Learn Go backend development with idiomatic patterns
- Demonstrate React + TypeScript frontend development
- Explore asynchronous job processing
- Understand how AI-assisted workflows could support installers

## Architecture

```
┌────────────┐       ┌────────────┐       ┌────────────┐
│            │  HTTP  │            │  Job   │            │
│  React UI  │◄─────►│  Go API    │──────►│  Worker    │
│  (Vite)    │       │  (chi)     │       │  (async)   │
│            │       │            │       │            │
└────────────┘       └────────────┘       └────────────┘
                           │                    │
                           ▼                    ▼
                     ┌────────────┐       ┌────────────┐
                     │ Repository │       │  Analyzer  │
                     │ (in-mem /  │       │ (mock /    │
                     │  Postgres) │       │  AI later) │
                     └────────────┘       └────────────┘
```

1. Installer submits building information via the React frontend.
2. The Go API creates an analysis job and returns immediately.
3. A background worker picks up the job, runs the analyzer, and stores results.
4. The frontend polls for status and renders structured recommendations.

## Tech Stack

| Layer    | Technology                     |
|----------|--------------------------------|
| Backend  | Go, chi router, context.Context |
| Frontend | React, TypeScript, Vite, Tailwind CSS |
| Storage  | In-memory (MVP), PostgreSQL (later) |
| Infra    | Docker Compose                 |

## Repository Structure

```
├── AGENTS.md              # Single source of truth
├── apps/
│   ├── api/               # Go backend
│   │   ├── cmd/server/    # Entrypoint
│   │   └── internal/      # Business logic (server, jobs, analyzer, worker, platform)
│   └── web/               # React frontend
│       ├── public/
│       └── src/            # api, components, pages, types
├── docs/
│   ├── product.md         # Product requirements
│   ├── architecture.md    # Technical architecture
│   └── api.md             # REST API specification
└── docker-compose.yml
```

## Getting Started

> Prerequisites: Go 1.22+, Node.js 20+, npm

```bash
# Backend
cd apps/api
go run ./cmd/server

# Frontend
cd apps/web
npm install
npm run dev
```

## Documentation

- [AGENTS.md](AGENTS.md) — project source of truth
- [docs/product.md](docs/product.md) — product requirements and scope
- [docs/architecture.md](docs/architecture.md) — technical architecture
- [docs/api.md](docs/api.md) — REST API specification
