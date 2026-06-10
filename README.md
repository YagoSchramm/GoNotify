# Notification Dispatch System

A production-grade asynchronous notification service built in Go. This project is designed to close real gaps in architecture, testing, and infrastructure — not as a portfolio toy, but as a system you maintain, operate, and debug under pressure.

---

## What This System Does

Users register and configure triggers. When a trigger fires, the system delivers a notification (email, webhook, etc.) with guaranteed delivery semantics: retry with exponential backoff, dead letter queue for unrecoverable failures, and an observable SLA.

---

## Tech Stack

- **Language:** Go
- **Router:** `gorilla/mux` + `net/http`
- **Database:** PostgreSQL
- **Queue/Cache:** Redis
- **Containerization:** Docker + Docker Compose
- **Testing:** `testing` stdlib + Testcontainers-Go
- **Logging:** `log/slog` (Go 1.21+ structured logging)
- **Deploy target:** VPS (Hetzner / Fly.io / Railway)

---

## Project Phases

### Phase 1 — Synchronous Core (2–3 weeks)

Get something working end-to-end before adding complexity.

- [ ] Initialize Go module and project structure
- [ ] Set up `gorilla/mux` router with versioned routes (`/api/v1/...`)
- [ ] Model the core domain:
  - `User` — who receives notifications
  - `Trigger` — conditions that fire a notification
  - `Notification` — the event to be delivered
  - `DeliveryAttempt` — a record of each delivery try
- [ ] Implement PostgreSQL schema with migrations (use `golang-migrate` or raw SQL files)
- [ ] Implement basic CRUD for users and triggers
- [ ] Deliver notifications **synchronously** first — HTTP call inline in the request handler
- [ ] Write a Dockerfile for the API service
- [ ] Write a `docker-compose.yml` with Postgres and Redis

**Checkpoint:** You can register a user, configure a trigger, fire it, and see a delivery attempt recorded in the database.

---

### Phase 2 — Make It Async (2–3 weeks)

Extract the delivery concern from the API. This is where bounded contexts become real, not theoretical.

- [ ] Design a job/event schema in Redis (or a Postgres-backed queue — justify your choice)
- [ ] Extract a **Worker** binary — separate `main.go`, separate process
- [ ] API enqueues a job; Worker picks it up and delivers
- [ ] Implement **exponential backoff retry**:
  - Attempt 1: immediate
  - Attempt 2: 30s delay
  - Attempt 3: 5min delay
  - Attempt 4: 30min delay
  - After max attempts: move to Dead Letter Queue (DLQ)
- [ ] Implement **Dead Letter Queue** — store failed notifications with full error context
- [ ] Implement **idempotency** — what happens if the worker processes the same job twice?
- [ ] Add a `/admin/dlq` endpoint to inspect and replay failed notifications
- [ ] Update `docker-compose.yml` to run both API and Worker services

**Checkpoint:** Kill the worker mid-run. Restart it. Verify no notifications are lost and no duplicates are created.

---

### Phase 3 — Tests That Matter (3–4 weeks)

No mocks for infrastructure. If your test doesn't use a real Postgres, it's not telling you the truth.

- [ ] Set up **Testcontainers-Go** to spin up real Postgres and Redis in tests
- [ ] Write **integration tests** for the repository layer:
  - Create/read/update delivery attempts
  - Assert retry count increments correctly
  - Assert DLQ behavior after max retries
- [ ] Write **API contract tests**:
  - All happy paths
  - Error cases (invalid payload, unknown trigger, duplicate event)
  - Assert correct HTTP status codes and response shapes
- [ ] Write **worker behavior tests**:
  - What happens when the downstream HTTP endpoint returns 500?
  - What happens when Redis is unavailable during enqueue?
  - What happens when Postgres is unavailable during job pickup?
- [ ] Write a **time-sensitive test**: verify retry delay logic without actually sleeping (use a clock interface — inject it)
- [ ] Set up a `Makefile` with `make test`, `make test-integration`, `make test-all`
- [ ] Measure and record test coverage — not to hit a number, but to find untested failure paths

**Checkpoint:** Run the full test suite against real containers. Every retry and failure scenario has a test. You can break any external dependency and the test tells you exactly what fails and why.

---

### Phase 4 — Deploy and Operate (2–3 weeks)

This is where everything abstract becomes concrete. You will discover things about your system you could not have discovered locally.

- [ ] Set up a VPS (Hetzner CX11 ~€4/mo, or Fly.io free tier)
- [ ] Write a production `docker-compose.yml` (separate from dev) with:
  - Resource limits
  - Restart policies (`restart: unless-stopped`)
  - Health checks for API and Worker
- [ ] Implement **structured logging** with `log/slog`:
  - Every request: method, path, status, duration
  - Every delivery attempt: notification ID, attempt number, outcome, latency
  - Every retry: reason for failure, next scheduled attempt
  - Every DLQ entry: full error context
- [ ] Implement a `/health` endpoint that checks Postgres and Redis connectivity
- [ ] Implement a `/metrics` endpoint (plain text is fine) exposing:
  - `notifications_enqueued_total`
  - `notifications_delivered_total`
  - `notifications_failed_total`
  - `notifications_in_dlq_total`
  - `delivery_latency_p99` (time from enqueue to delivery)
- [ ] Define your SLA: "X% of notifications delivered within 30 seconds"
- [ ] Set up basic alerting: if DLQ grows beyond N items, log a warning
- [ ] Write a `RUNBOOK.md` — how do you debug a notification that didn't deliver?

**Checkpoint:** Deploy to VPS. Fire 100 notifications. Read the logs and confirm you can trace every single one from enqueue to delivery (or DLQ). You should be able to answer "what happened to notification X?" in under 2 minutes.

---

## Key Design Decisions You Must Make

These are not answered here intentionally. Your job is to decide, implement, and live with the consequences.

| Question | Why It Matters |
|---|---|
| Postgres queue vs Redis queue? | Durability vs performance trade-off. Wrong choice = silent message loss. |
| How do you model idempotency keys? | Without this, retries cause duplicate deliveries. |
| Where do you store retry state? | In the job? In Postgres? What's the source of truth? |
| How do you inject the clock for testing? | Forces you to separate pure logic from side effects. |
| What's your migration strategy? | `golang-migrate`? Raw SQL? Embedded migrations? Each has a cost. |
| How do you handle partial failure? | Notification saved to DB but not enqueued — what now? |

---

## Definitions of Done

A phase is only done when:

1. The code compiles and all tests pass
2. You can explain every design decision and its trade-off
3. The `docker-compose up` works on a clean machine with no manual steps
4. You've broken it intentionally (kill a dependency) and verified the behavior

---
