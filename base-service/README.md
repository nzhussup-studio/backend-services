# Base Service

`base-service` is the Spring Boot backend that serves structured portfolio and CV content.

## Responsibilities

- projects
- work experience
- education
- skills
- certifications

## Auth Model

- `GET` endpoints are public
- Envoy Gateway validates JWTs and enforces admin access before requests reach this service.
- The service accepts only gateway-injected `X-User-*` identity headers; network policy blocks direct callers.
- `GET /v1/base/cv-generator-preference` requires an authenticated user; public GET endpoints remain public.

## CV Generator Preferences

- Stored in `cv_generator_preferences` table as JSON.
- CRUD endpoint: `/v1/base/cv-generator-preference`.

## Runtime Inputs

The service expects database and Redis settings:

- `SPRING_DATASOURCE_HOST`
- `SPRING_DATASOURCE_DB`
- `SPRING_DATASOURCE_PORT`
- `SPRING_DATASOURCE_USERNAME`
- `SPRING_DATASOURCE_PASSWORD`
- `SPRING_REDIS_HOST`
- `SPRING_REDIS_PORT`
- Trusted identity is supplied by the API gateway through `X-User-*` headers.
