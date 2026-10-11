#!/usr/bin/env bash
set -euo pipefail

: "${BUCKET:?BUCKET is required}"
REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
KEY_FILE="${KEY_FILE:-/secrets/key.json}"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="redis_backup_${TIMESTAMP}.rdb"
BACKUP_DIR="/tmp/redis-backups"
TMP="${BACKUP_DIR}/${BACKUP_FILE}"
COMPRESSED="${TMP}.gz"

mkdir -p "$BACKUP_DIR"
trap 'rm -f "$TMP" "$COMPRESSED"' EXIT

echo "Dumping Redis at ${REDIS_HOST}:${REDIS_PORT}"
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" --rdb "$TMP"

# Check if the dump file is valid
if [ ! -s "$TMP" ] || [ "$(head -c 5 "$TMP")" != "REDIS" ]; then
  echo "Dump is empty or invalid, aborting" >&2
  exit 1
fi

echo "Compressing"
gzip -f "$TMP"          # replaces $TMP with $TMP.gz

echo "Uploading to gs://${BUCKET}/${BACKUP_FILE}.gz"
gcloud auth activate-service-account --key-file="$KEY_FILE" --quiet
gcloud storage cp "$COMPRESSED" "gs://${BUCKET}/${BACKUP_FILE}.gz"

echo "Backup done: ${BACKUP_FILE}.gz"