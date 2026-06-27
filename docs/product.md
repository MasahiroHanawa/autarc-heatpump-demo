# Product Requirements

## Target Users

**Primary: Heat pump installers and energy consultants** who perform on-site surveys, gather building data, and need quick preliminary assessments before committing to a full engineering analysis.

**Secondary: Small installation companies** looking for a digital tool to standardize their quoting workflow and reduce the time between site visit and customer proposal.

## Business Problem

Heat pump sizing and suitability assessment is currently a manual, experience-dependent process. Installers gather data on-site (building type, insulation, existing heating system) and then spend hours cross-referencing tables and guidelines to determine the right heat pump type and capacity.

Key pain points:
- Assessments are slow — days between site visit and recommendation
- Quality varies — depends on individual installer experience
- No structured output — customers receive inconsistent formats
- Difficult to scale — each assessment requires senior expertise

## MVP Scope

The MVP delivers a single workflow: **submit building data, receive a structured recommendation.**

### Included in MVP

- **Job submission** — installer enters building parameters (type, area, construction year, insulation level, current heating, energy consumption, notes)
- **Async processing** — backend queues the analysis and processes it in a background worker
- **Status polling** — frontend polls for job progress (pending → processing → completed/failed)
- **Structured results** — display estimated heat demand, recommended pump type, suitability score, confidence level, risk flags, and next steps
- **Loading and error states** — clear feedback during processing and on failure
- **Mock analyzer** — deterministic analysis logic for demo purposes; designed to be replaceable with AI providers

### Excluded from MVP

- User authentication
- Multi-tenant support
- Real AI/ML integration (prepared-for but not implemented)
- PDF report generation
- Photo/document upload

## Feature Priorities

### Must Have (MVP)

1. Create analysis job via REST API
2. Background worker with async job processing
3. Job status polling endpoint
4. Structured result rendering in the frontend
5. Loading, error, and empty states
6. Mock analyzer with deterministic output

### Should Have (Post-MVP)

7. Job history list — view past analyses
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

## Success Metrics

For this demo, success is measured by:

- A working end-to-end flow from submission to result display
- Clean separation of concerns allowing analyzer replacement
- Code that reads as a realistic production service, not a tutorial
