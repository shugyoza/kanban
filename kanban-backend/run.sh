#!/bin/bash

# DO NOT use in production!

# Exit instantly if any subcommand throws an unexpected fault code
set -e

# Prompt the user for a Yes/No response
read -p "Do you want to fresh reset the database? [y/N]: " response

# Execute file purge ONLY if the fresh flag was explicitly provided
if [[ "$response" = "Y" || "$response" = "y" || "$response" = "yes" || "$response" = "YES" ]]; then
  echo "⚠️ Warning: Fresh reset requested. Clearing out old system states..."
  
  # Safely remove the old database file if it exists
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