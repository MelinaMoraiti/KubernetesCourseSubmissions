#!/bin/sh

set -e

REDIS_HOST="${REDIS_HOST}"
REDIS_PORT="${REDIS_PORT}"
TASK_NAME="${TASK_NAME}"

# Get the next available ID
ID=$(redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" INCR todo:next_id)

# Create the todo
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" \
    SET "todo:$ID" \
    "{\"id\":$ID,\"task\":\"$TASK\",\"done\":false}"

# Add the ID to the todos index
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" \
    SADD todos "$ID"
