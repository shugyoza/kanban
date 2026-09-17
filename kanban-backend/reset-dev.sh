#!/bin/bash

# DO NOT use in production!

# Exit instantly if any subcommand throws an unexpected fault code
set -e

# Initialize a flag variable
FRESH_RESET=false

# Parse command line arguments
for arg in "$@"; do
  case $arg in
    -f|--fresh)
      FRESH_RESET=true
      shift
      ;;
  esac
done

# Execute file purge ONLY if the fresh flag was explicitly provided
if [ "FRESH_RESET" = true ]; then
  echo "⚠️ Warning: Fresh reset requested. Clearing out old system states..."# 1. Safely remove the old database file if it exists

  if [ -f "./kanban.db" ]; then
    rm "./kanban.db"
    echo "Old local kanban.db file successfully purged."
  else
    echo "No existing kanban.db file found. Skipping file purge."
  fi
else
  echo "🚀 Normal Boot: Preserving existing local database file..."
fi

echo "Booting up Kanban Backend monolith..."
go run cmd/api/main.go