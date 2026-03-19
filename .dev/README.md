# Local Development Setup

This directory contains the local Docker Compose setup for the backend platform. It starts the main runtime dependencies and all backend services in a single local environment.

## What Starts

The compose stack includes:

- PostgreSQL
- Redis
- `auth-service`
- `base-service`
- `user-service`
- `image-service`
- `llm-service`
- `api-gateway`
- `db-seed` one-shot seeding container

Only these ports are exposed to the host:

- `8082` for `api-gateway`
- `5432` for PostgreSQL
- `6379` for Redis

The backend application services themselves are kept internal to the Compose network and are intended to be accessed through the API gateway.

## Start

Run from the repository root:

```bash
docker compose -f backend-services/.dev/compose.yml up --build
```

The stack works without a `.env` file because the compose file already defines defaults for all required variables.

## Optional Environment Overrides

You can override values at runtime if needed:

```bash
POSTGRES_PORT=5433 REDIS_PORT=6380 API_GATEWAY_PORT=8088 docker compose -f backend-services/.dev/compose.yml up --build
```

If you want the LLM integration to work locally, provide an API key:

```bash
OPENROUTER_API_KEY=your_key docker compose -f backend-services/.dev/compose.yml up --build
```

## Seed Data

The local setup includes a database seed step:

- `seed.sh` waits until PostgreSQL is ready and the application tables exist
- `seed.sql` clears existing seeded content and inserts a prepared dataset

The seed contains:

- two users
- work experience entries
- education entries
- projects
- skills
- certificates

The seed is intended for local development and is re-applied when the stack starts.

## Seeded Users

### Admin User

- username: `admin`
- password: `admin123!`
- role: `ROLE_ADMIN`

### Normal User

- username: `jane.user`
- password: `user123!`
- role: `ROLE_USER`

## Notes

- `image-service` file content is not seeded, because image data is stored on disk rather than in PostgreSQL.
- `llm-service` starts without an API key, but summarization requests will fail until `OPENROUTER_API_KEY` is set.
- Java services use local-friendly defaults for database and auth service access, while Compose overrides them with container hostnames.
