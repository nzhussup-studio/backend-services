# Backend Services

This repository contains the backend runtime for `nzhussup.com` and its admin ecosystem. It is organized as a set of focused services rather than a single application, which keeps deployment, scaling, and operational ownership explicit.

## Services

- `api-gateway`  
  Unified entrypoint for backend APIs and request routing.

- `auth-service`  
  Authentication and token-related flows for protected access.

- `base-service`  
  Core data APIs used to manage portfolio and profile content.

- `user-service`  
  User management and admin-facing user operations.

- `image-service`  
  Image and album handling, including media-related backend logic.

- `llm-service`  
  LLM-backed capabilities used by the platform.

- `redis-server`  
  Containerized Redis runtime used as infrastructure support for backend workloads.

## Stack

- `Go` for lightweight network services
- `Java / Spring Boot` for structured business APIs
- `Redis` for caching and supporting runtime concerns
- `Docker` for packaging
- `GitHub Actions` for CI/CD

## Delivery Model

The repository uses a single manifest-driven CI/CD pipeline defined under `.github/`.

- Changed services are detected automatically.
- Services marked as `full` run quality, test, build, and deploy stages.
- Services marked as `build-deploy` skip code validation and only build and deploy.
- Production image builds target `linux/amd64`, matching the current server runtime.

## Structure

Each service is isolated in its own top-level directory with its own code, build definition, and runtime concerns. Shared CI/CD configuration lives under `.github/`.

## OpenAPI Generation

This repository includes a repo-level OpenAPI generator at `scripts/generate-openapi.sh`.

### Usage

Generate OpenAPI for a single service:

```bash
./scripts/generate-openapi.sh <path-to-folder>
```

Examples:

```bash
./scripts/generate-openapi.sh auth-service
./scripts/generate-openapi.sh image-service
```

Generate OpenAPI for all configured services and build a unified spec:

```bash
./scripts/generate-openapi.sh --all
```

The `--all` flow reads the service list from `openapi-services.json`.

### Output

- Each supported service gets its own `openapi.yaml` in the service root.
- Running `--all` also creates a unified `openapi.yaml` in the repository root.

### How It Works

- Java services are generated through `springdoc`.
- Go services are generated through `swag`.
- The generator sanitizes the output so it is suitable for client generation:
  - removes generation-time local server values
  - normalizes Go Swagger 2 output into OpenAPI 3 format
  - keeps the root merged spec in OpenAPI 3 format

For Java services, the generator uses an `openapi` runtime profile so specs can be exported without requiring the full production infrastructure stack.

### Supported Services

The generator currently supports the services listed in `openapi-services.json`.

`api-gateway` and `redis-server` are intentionally excluded from OpenAPI generation.
