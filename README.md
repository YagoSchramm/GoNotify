# GoNotify

A notification service written in Go. This repository currently implements an HTTP API with PostgreSQL-backed user auth, trigger management, and notification firing.

## What This Project Contains

- HTTP API built with `net/http` and `gorilla/mux`
- JWT authentication for protected routes
- PostgreSQL connection via `jackc/pgx`
- Domain models for users, triggers, notifications, and delivery attempts
- Database migration file in `infrastructure/script/migrate/001-create-tables.up.sql`
- Postgres helper `Dockerfile.postgres` for local database initialization

## Current Features

- `POST /auth/register` — register a new user
- `POST /auth/login` — login and receive a bearer token
- `POST /triggers` — create a trigger for the authenticated user
- `GET /triggers` — list triggers for the authenticated user
- `GET /triggers/{id}` — get a trigger by ID
- `PATCH /triggers/{id}` — update a trigger
- `DELETE /triggers/{id}` — delete a trigger
- `POST /notifications` — fire a notification from a trigger
- `GET /notifications/{id}` — get notification status by ID

## Requirements

- Go 1.25+
- PostgreSQL database
- `DATABASE_URL` environment variable
- `JWT_SECRET` environment variable

## Quick Start

1. Copy the example env file:

   ```powershell
   copy .env-example .env
   ```

2. Edit `.env` and set `DATABASE_URL` and `JWT_SECRET`.

3. Start a PostgreSQL instance.
   - Use `Dockerfile.postgres` or any PostgreSQL server.
   - If you use `Dockerfile.postgres`, build and run it manually.

4. Run the service:

   ```powershell
   go run main.go
   ```

5. The API listens on `PORT` from `.env` or defaults to `8080`.

## `.env` Variables

- `PORT` — HTTP server port (default `8080`)
- `DATABASE_URL` — PostgreSQL connection string
- `JWT_SECRET` — secret used to sign JWT tokens

## Database Setup

The repository includes a migration SQL file:

- `infrastructure/script/migrate/001-create-tables.up.sql`

A local Postgres container can be created from `Dockerfile.postgres`, which initializes the database schema at startup.

## Notes

- The current implementation does not include Redis or a worker process.
- The API uses a simple JWT bearer scheme for auth.
- Protected routes require an `Authorization: Bearer <token>` header.

## Development

To add or modify behavior:

- `main.go` starts the HTTP server and loads env config
- `service/service.go` builds the router and database connection
- `infrastructure/router/module` defines HTTP modules and route handlers
- `domain/dto` contains request and response payloads
- `infrastructure/datastore/repository/impl` contains PostgreSQL repositories

## License

No license is specified in this repository.
