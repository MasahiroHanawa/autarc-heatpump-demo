# Heat Pump Planning Assistant

## Goal

Build a realistic fullstack demo inspired by heat pump installer workflows.

The purpose of this project is to demonstrate:

- Go backend development
- React frontend development
- asynchronous job processing
- AI-style structured outputs
- clean architecture
- product thinking

This project is intentionally designed to resemble a production SaaS feature rather than a coding exercise.

---

## Product Concept

Installers enter building information and site-survey notes.

The backend creates an analysis job.

A worker processes the job asynchronously.

The frontend polls for updates and renders structured recommendations.

The current implementation uses a mock analyzer but must be designed so that OpenAI, Anthropic, or other multimodal providers can replace it later.

---

## Feature Priorities

### Must Have (MVP)

1. Create analysis job via REST API (`POST /api/jobs`)
2. Background worker with async job processing
3. Job status polling endpoint (`GET /api/jobs/:id`)
4. Structured result rendering in the frontend
5. Loading, error, and empty states in the UI
6. Mock analyzer with deterministic output

### Should Have (Post-MVP)

7. Job history list — view past analyses (`GET /api/jobs`)
8. PostgreSQL persistence — replace in-memory store
9. Retry handling — automatic retry on transient failures
10. Input validation — server-side and client-side
11. Docker Compose — one-command local development
12. Basic unit tests — handler, worker, and analyzer coverage

### Nice To Have (Future)

13. OpenAI / Anthropic integration — real AI-powered analysis
14. Photo upload — site survey images as analysis input
15. PDF report export — generate customer-facing documents
16. Voice note transcription — convert field notes to text
17. Metrics and tracing — observability with OpenTelemetry
18. Multi-language support — German, English at minimum

---

## Core Domain Model

Input:

- buildingType
- livingAreaM2
- constructionYear
- insulationLevel
- currentHeatingSystem
- annualEnergyConsumptionKwh
- installerNotes

Output:

- estimatedHeatDemandKwh
- recommendedHeatPumpType
- suitability
- confidence
- riskFlags
- nextSteps
- summary

---

## Tech Stack

Backend:

- Go
- chi router
- context.Context
- in-memory repository first
- PostgreSQL later

Frontend:

- React
- TypeScript
- Vite
- Tailwind CSS

Development:

- Docker Compose
- README-first development

---

## Repository Structure

apps/api:

- cmd/server
- internal/server
- internal/jobs
- internal/analyzer
- internal/worker
- internal/platform

apps/web:

- src/api
- src/components
- src/pages
- src/types

---

## Engineering Principles

- Keep product workflows realistic.
- Prefer explicit Go code over abstractions.
- Use context.Context everywhere appropriate.
- Keep business logic outside handlers.
- Make analyzer implementations replaceable.
- Avoid overengineering.
- Write small, testable units.
- Optimize for clarity and maintainability.

---

## Future Extensions

- OpenAI Vision integration
- Voice note analysis
- PDF processing
- PostgreSQL persistence
- Retry policies
- Metrics and tracing