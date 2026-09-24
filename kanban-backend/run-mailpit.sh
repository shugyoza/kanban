#!/bin/bash

# DO NOT use in production!

# Exit instantly if any subcommand throws an unexpected fault code
set -e

DOCKER_CONTAINER_NAME="kanban-mailpit"
DOCKER_CONTAINER_IMAGE_NAME="axllent/mailpit"

# Spin / re-spin a docker container
echo "Running mailpit docker container for sending invitation email..."

docker_container_exists_as_id=$(docker ps -aqf "name=^/$DOCKER_CONTAINER_NAME") # short 12-character Container ID
docker_container_running_as_id=$(docker ps -qf "name=^/$DOCKER_CONTAINER_NAME" -f "status=running") # short 12-character Container ID
docker_container_is_dead_as_id=$(docker ps -aqf "name=^/$DOCKER_CONTAINER_NAME" -f "status=dead") # short 12-character Container ID

if [ -n "$docker_container_exists_as_id" ] && [ -n "$docker_container_running_as_id" ]; then # -n checks if a string is not empty
  echo "Docker container $DOCKER_CONTAINER_NAME exists and is already running..."
elif [ -n "$docker_container_exists_as_id" ] && [ -z "$docker_container_running_as_id" ] && [ -z "$docker_container_is_dead_as_id" ]; then # -z checks if a string is empty (zero length)
  echo "Docker container $DOCKER_CONTAINER_NAME exists, but not running. Restarting..."
  docker start "$docker_container_exists_as_id"
elif [ -n "$docker_container_exists_as_id" ] && [ -z "$docker_container_running_as_id" ] && [ -n "$docker_container_is_dead_as_id" ]; then
  echo "Purging the old container to free up the name..."

  if docker rm "$DOCKER_CONTAINER_NAME" >/dev/null 2>&1; then
    echo "Successfully purged the old container layer..."
    echo "Reinitializing a new docker container..."
    docker run -d --name "$DOCKER_CONTAINER_NAME" \
    -p 1025:1025 \
    -p 8025:8025 \
    "$DOCKER_CONTAINER_IMAGE_NAME"
  else
    echo "Failed to purge the container natively. Exiting..."
    exit 1
  fi
else
  echo "Docker container $DOCKER_CONTAINER_NAME does not exist. Spinning a new one..."
  docker run -d --name "$DOCKER_CONTAINER_NAME" \
  -p 1025:1025 \
  -p 8025:8025 \
  "$DOCKER_CONTAINER_IMAGE_NAME"
fi

# Clean up
# docker stop "$DOCKER_CONTAINER_NAME"
# docker rm "$DOCKER_CONTAINER_NAME"

echo "done."
