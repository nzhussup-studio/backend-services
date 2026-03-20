# Keycloak Service

This directory contains the source for the platform Keycloak image.

## Files

- `Dockerfile`
- `backend-auth-realm.template.json`
- `docker-entrypoint.sh`
- `dev.env.example`
- `prod.env.example`

## Runtime Model

The image bakes in a generic realm template.
At startup, `docker-entrypoint.sh` renders a concrete import file from environment variables into Keycloak's import directory and then starts Keycloak with `--import-realm`.
The template itself is stored outside the import directory so Keycloak does not try to import placeholder values directly.

Live secrets are not stored in git.

## Required Runtime Variables

- `BACKEND_AUTH_CLIENT_SECRET`

## Common Realm Variables

- `KEYCLOAK_REALM_ID`
- `KEYCLOAK_REALM_NAME`
- `KEYCLOAK_REALM_DISPLAY_NAME`
- `KEYCLOAK_SSL_REQUIRED`
- `KEYCLOAK_REGISTRATION_ALLOWED`
- `KEYCLOAK_LOGIN_WITH_EMAIL_ALLOWED`
- `KEYCLOAK_RESET_PASSWORD_ALLOWED`
- `KEYCLOAK_REMEMBER_ME`
- `KEYCLOAK_VERIFY_EMAIL`

## Client Variables

- `KEYCLOAK_BACKEND_CLIENT_ID`
- `KEYCLOAK_BACKEND_CLIENT_NAME`
- `KEYCLOAK_BACKEND_STANDARD_FLOW_ENABLED`
- `KEYCLOAK_BACKEND_DIRECT_ACCESS_GRANTS_ENABLED`
- `KEYCLOAK_FRONTEND_CLIENT_ID`
- `KEYCLOAK_FRONTEND_CLIENT_NAME`
- `KEYCLOAK_FRONTEND_ROOT_URL`
- `KEYCLOAK_FRONTEND_BASE_URL`
- `KEYCLOAK_FRONTEND_REDIRECT_URIS_JSON`
- `KEYCLOAK_FRONTEND_WEB_ORIGINS_JSON`
- `KEYCLOAK_USERS_FILE`

`KEYCLOAK_FRONTEND_REDIRECT_URIS_JSON` and `KEYCLOAK_FRONTEND_WEB_ORIGINS_JSON` must be valid JSON arrays.
`KEYCLOAK_USERS_FILE` is optional and should point to a JSON array of Keycloak users for environment-specific seed users.

For CORS and browser redirects, use exact origins and hosts rather than parent domains.
Example production list:

- `https://admin.nzhussup.dev`
- `https://admin.nzhussup.com`

Only keep the hosts you actually need. Subdomains are separate origins for CORS.

## Social Login Variables

- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GITHUB_CLIENT_ID`
- `GITHUB_CLIENT_SECRET`
- `KEYCLOAK_GOOGLE_ENABLED`
- `KEYCLOAK_GITHUB_ENABLED`

If `KEYCLOAK_GOOGLE_ENABLED` or `KEYCLOAK_GITHUB_ENABLED` are not set, the entrypoint enables the provider automatically only when both client ID and secret are present.

Use `dev.env.example` and `prod.env.example` as the concrete starting points for local and production runtime values.

The intended split is:

- base realm structure in `backend-auth-realm.template.json`
- no seeded users in production
- optional seeded users in development through `KEYCLOAK_USERS_FILE`

The backend client is configured as a confidential API client, but service accounts are intentionally disabled by default.
That keeps the default realm focused on browser-issued user tokens for the backend API.
