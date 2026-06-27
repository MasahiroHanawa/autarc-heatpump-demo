# Technical Architecture

## Overview

The system follows a simple job-queue pattern: the API accepts requests, creates jobs, and a background worker processes them asynchronously. The frontend polls for updates.

```
Client → HTTP Handler → Repository (store job) → return job ID
                                    ↓
                              Worker (loop)
                                    ↓
                              Analyzer (process)
                                    ↓
                              Repository (store result)
                                    ↓
Client ← HTTP Handler ← Repository (read result)
```

## Backend Package Structure (Go)

```
apps/api/
├── cmd/
│   └── server/
│       └── main.go              # Entrypoint: wires dependencies, starts HTTP server
├── internal/
│   ├── server/
│   │   ├── handler.go           # HTTP handlers (create job, get job, list jobs)
│   │   ├── router.go            # chi router setup and middleware
│   │   └── response.go          # JSON response helpers
│   ├── jobs/
│   │   ├── model.go             # Domain types: Job, JobInput, JobResult, JobStatus
│   │   ├── service.go           # Business logic: depends on Repository interface
│   │   ├── repository.go        # Repository interface
│   │   └── memory.go            # In-memory repository implementation
│   ├── analyzer/
│   │   ├── analyzer.go          # Analyzer interface
│   │   └── mock.go              # Mock analyzer with deterministic logic
│   ├── worker/
│   │   └── worker.go            # In-process background worker using goroutines and periodic polling.
│   │                            # Can later be replaced with a dedicated queue consumer.
│   └── platform/
│       └── config.go            # Configuration loading (env vars, defaults)
├── go.mod
└── go.sum
```

### Package Responsibilities

| Package    | Role | Depends On |
|------------|------|------------|
| `cmd/server` | Wires all dependencies and starts the server | All internal packages |
| `internal/server` | HTTP transport layer — routes requests to domain logic | `jobs` |
| `internal/jobs` | Domain model, service logic, and persistence interface | Repository abstraction only |
| `internal/analyzer` | Analysis logic behind an interface | `jobs` (for types) |
| `internal/worker` | Background processing loop | `jobs`, `analyzer` |
| `internal/platform` | Cross-cutting concerns (config, logging) | None |

### Design Decisions

- **`internal/` prefix** — Go convention; prevents external imports of implementation details.
- **Interface-based analyzer** — `analyzer.Analyzer` is an interface so mock, rule-based, and AI implementations are interchangeable.
- **Repository pattern** — `jobs.Repository` is an interface; `memory.go` implements it. Postgres can be added later without changing handlers or workers.
- **`server` over `http`** — avoids ambiguity with Go's `net/http` standard library package.
- **No framework** — chi is a lightweight router, not a framework. Handlers are plain `http.HandlerFunc`.
- **`context.Context` everywhere** — passed through handlers → repository → analyzer for cancellation and timeouts.

## Frontend Folder Structure (React + TypeScript)

```
apps/web/
├── public/
│   └── favicon.ico
├── src/
│   ├── api/
│   │   └── client.ts            # HTTP client: createJob(), getJob(), listJobs()
│   ├── components/
│   │   ├── JobForm.tsx           # Input form for building data
│   │   ├── JobStatus.tsx         # Status badge (pending/processing/completed/failed)
│   │   ├── JobResult.tsx         # Structured result display
│   │   └── Layout.tsx            # Page layout wrapper
│   ├── pages/
│   │   ├── NewJob.tsx            # Create new analysis page
│   │   └── JobDetail.tsx         # View job status and results
│   ├── types/
│   │   └── index.ts             # TypeScript types mirroring backend models
│   ├── App.tsx                  # Root component with routing
│   └── main.tsx                 # Vite entry point
├── index.html
├── package.json
├── tsconfig.json
├── tailwind.config.js
└── vite.config.ts
```

### Frontend Design Decisions

- **Polling over WebSockets** — simpler to implement; sufficient for job status updates that take seconds.
- **No state management library** — React state + polling is enough for MVP. Add Zustand or similar only if complexity warrants it.
- **Collocated types** — `src/types/` mirrors backend models to keep the contract explicit.
- **API client module** — centralizes all HTTP calls; easy to swap base URL or add auth headers later.

## Data Flow

### Create Job

1. User fills `JobForm` and submits
2. `client.createJob(input)` → `POST /api/jobs`
3. Handler validates input, calls `repository.Create(job)`
4. Returns `202 Accepted` with job ID
5. Frontend navigates to `JobDetail` page

### Process Job

1. Worker polls `repository.ListPending()` on interval
2. Picks up job, sets status to `processing`
3. Calls `analyzer.Analyze(input)` → returns result
4. Sets status to `completed` with result (or `failed` with error)

### View Result

1. `JobDetail` page polls `GET /api/jobs/:id` every 2 seconds
2. Renders `JobStatus` badge
3. When completed, renders `JobResult` with structured data
4. Stops polling when terminal state reached

## Future Architecture Changes

| Change | Impact |
|--------|--------|
| Add PostgreSQL | New `internal/jobs/postgres.go` implementing `Repository` interface |
| Add AI analyzer | New `internal/analyzer/openai.go` implementing `Analyzer` interface |
| Add authentication | Middleware in `internal/server/router.go` |
| Add file uploads | New `internal/storage/` package + multipart handler |
