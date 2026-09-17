#!/bin/bash

set -e

until pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB"; do
    echo "Waiting for PostgreSQL..."
    sleep 2
done

psql \
    -v ON_ERROR_STOP=1 \
    --username "$POSTGRES_USER" \
    --dbname "$POSTGRES_DB" <<-EOSQL

CREATE TABLE IF NOT EXISTS counter (
    id INTEGER PRIMARY KEY,
    value INTEGER NOT NULL
);

INSERT INTO counter (id, value)
VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;

EOSQL

echo "Database initialized."