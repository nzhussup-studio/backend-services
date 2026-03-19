# Local Development Setup

This directory contains the local Docker Compose setup for the backend platform. It starts the main runtime dependencies and all backend services in a single local environment.

## What Starts

The compose stack includes:

- PostgreSQL
- Keycloak
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
- `8081` for Keycloak
- `5432` for PostgreSQL
- `6379` for Redis

The backend application services themselves are kept internal to the Compose network and are intended to be accessed through the API gateway.

## Start

Run from the repository root:

```bash
docker compose -f backend-services/.dev/compose.yml up --build
```

The stack works without a `.env` file because the compose file already defines defaults for all required variables.

## Keycloak

The local stack includes a Keycloak instance for IAM migration work.

- URL: `http://localhost:8081`
- admin username: `root`
- admin password: `root`
- realm: `backend-auth-dev`
- shared realm template: `keycloak-service/backend-auth-realm.template.json`
- dev user seed file: `.dev/.keycloak/dev-users.json` mounted outside Keycloak's import directory
- published on all host interfaces via `0.0.0.0:${KEYCLOAK_PORT:-8081}`

The admin credentials default to `root:root` and can be overridden with:

```bash
KEYCLOAK_ADMIN_USERNAME=your_user KEYCLOAK_ADMIN_PASSWORD=your_pass docker compose -f backend-services/.dev/compose.yml up --build
```

If you want Keycloak to behave like it is already behind a dedicated host or subdomain, set the hostname explicitly:

```bash
KEYCLOAK_HOSTNAME=auth.local.nzhussup.dev docker compose -f backend-services/.dev/compose.yml up --build
```

For free-form local access by IP or alternate hostnames, the defaults keep hostname checks relaxed:

- `KC_HOSTNAME_STRICT=false`
- `KC_HOSTNAME_STRICT_HTTPS=false`
- `KC_PROXY_HEADERS=xforwarded`

The local stack also patches the `master` realm to `sslRequired=NONE` after startup so the admin console works over plain HTTP on localhost.

The dev render currently configures:

- realm roles `ROLE_ADMIN` and `ROLE_USER`, with `ROLE_USER` as the default role for new users
- confidential client `backend-auth-client`
- public client `frontend-admin-auth-client`
- self-service user registration enabled
- Google and GitHub identity providers enabled only when credentials are provided
- local users `admin` and `jane.user` through the dev-only users JSON

Keycloak persists its state in a dedicated PostgreSQL database named `keycloak` by default. The local stack creates that database automatically through a small init container before Keycloak starts, so it also works with existing Postgres volumes.

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
