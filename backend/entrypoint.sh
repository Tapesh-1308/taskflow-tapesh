#!/bin/sh

echo "⏳ Waiting for database..."

until nc -z db 5432; do
  sleep 2
done

echo "✅ Database is up"

echo "🚀 Running migrations..."

migrate -path /app/migrations \
  -database "$DATABASE_URL" \
  up || true

echo "🔥 Starting backend..."

./main