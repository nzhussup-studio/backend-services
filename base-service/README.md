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
- non-`GET` endpoints require `ROLE_ADMIN`
- JWT validation is performed locally against Keycloak JWKS
- small exception: `GET /v1/base/cv-generator-preference` also requires `ROLE_ADMIN`

## CV Generator Preferences

- Stored in `cv_generator_preferences` table as JSON.
- CRUD endpoint: `/v1/base/cv-generator-preference`.

## Runtime Inputs

The service expects database and Redis settings plus Keycloak JWT configuration:

- `SPRING_DATASOURCE_HOST`
- `SPRING_DATASOURCE_DB`
- `SPRING_DATASOURCE_PORT`
- `SPRING_DATASOURCE_USERNAME`
- `SPRING_DATASOURCE_PASSWORD`
- `SPRING_REDIS_HOST`
- `SPRING_REDIS_PORT`
- Trusted identity is supplied by the API gateway through `X-User-*` headers.
