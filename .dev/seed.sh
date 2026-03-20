#!/bin/sh

set -eu

export PGPASSWORD="${POSTGRES_PASSWORD}"

echo "Waiting for PostgreSQL..."
until pg_isready -h postgres -p 5432 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; do
  sleep 2
done

echo "Waiting for application tables..."
until [ "$(psql -h postgres -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('work_experience', 'education', 'projects', 'skills', 'certificates');")" = "5" ]; do
  sleep 2
done

echo "Applying development seed..."
psql -h postgres -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -v ON_ERROR_STOP=1 -f /seed/seed.sql

echo "Development seed applied."
