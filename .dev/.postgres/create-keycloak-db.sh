#!/bin/sh

set -eu

export PGPASSWORD="${POSTGRES_PASSWORD}"

echo "Waiting for PostgreSQL..."
until pg_isready -h postgres -p 5432 -U "${POSTGRES_USER}" >/dev/null 2>&1; do
  sleep 2
done

DB_EXISTS="$(psql -h postgres -U "${POSTGRES_USER}" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = '${KEYCLOAK_DB}'")"

if [ "${DB_EXISTS}" = "1" ]; then
  echo "Database '${KEYCLOAK_DB}' already exists."
  exit 0
fi

echo "Creating database '${KEYCLOAK_DB}'..."
psql -h postgres -U "${POSTGRES_USER}" -d postgres -v ON_ERROR_STOP=1 -c "CREATE DATABASE \"${KEYCLOAK_DB}\";"
echo "Database '${KEYCLOAK_DB}' created."
