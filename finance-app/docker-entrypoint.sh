#!/bin/sh

echo "Running migrations..."

./goose -dir /app/migrations postgres "$DB_DSN" up

echo "Starting app..."

./main