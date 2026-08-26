#!/bin/bash

# DO NOT use in production!

# Exit instantly if any subcommand throws an unexpected fault code
set -e

echo "Stopping and clearing out old system states..."

# 1. Safely remove the old database file if it exists
if [ -f "./kanban.db" ]; then
  rm "./kanban.db"
  echo "Old local kanban.db file successfully purged."
else
  echo "No existing kanban.db file found. Skipping file purge."
fi

echo "Re-initializing database tables and seeding fresh models..."

# 2. Boot up the Go environment, which automatically processes schema.sql
echo "Booting up Kanban Backend monolith..."
go run cmd/api/main.go