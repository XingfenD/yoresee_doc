#!/bin/bash

# Check parameter count
if [ $# -lt 2 ]; then
    echo "Usage: $0 <dev|release> <rebuild|clear|restart|up>"
    echo "  dev|release: Use docker-compose.dev.yml or docker-compose.yml"
    echo "  rebuild: Clean containers and volumes, then rebuild and start"
    echo "  clear: Clean containers and volumes"
    echo "  restart: Stop and restart services"
    echo "  up: Start services"
    exit 1
fi

ENVIRONMENT=$1
ACTION=$2

# Parameter validation
if [[ "$ENVIRONMENT" != "dev" && "$ENVIRONMENT" != "release" ]]; then
    echo "Error: First argument must be 'dev' or 'release'"
    exit 1
fi

if [[ "$ACTION" != "rebuild" && "$ACTION" != "clear" && "$ACTION" != "restart" && "$ACTION" != "up" ]]; then
    echo "Error: Second argument must be 'rebuild', 'clear', 'restart', or 'up'"
    exit 1
fi

ORIGINAL_DIR=$(pwd)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Environment: $ENVIRONMENT"
echo "Action: $ACTION"
echo "Script directory: $SCRIPT_DIR"
echo "Base path: $BASE_PATH"

cd "$BASE_PATH"/deploy/

# Select compose file based on environment
if [ "$ENVIRONMENT" = "dev" ]; then
    COMPOSE_FILE="docker-compose.dev.yml"
else
    COMPOSE_FILE="docker-compose.yml"
fi

echo "Using compose file: $COMPOSE_FILE"

case "$ACTION" in
    "rebuild")
        echo "Rebuilding and starting services..."
        docker compose -f "$COMPOSE_FILE" down -v
        docker compose -f "$COMPOSE_FILE" up -d --build
        ;;
    "clear")
        echo "Clearing containers and volumes..."
        docker compose -f "$COMPOSE_FILE" down -v
        ;;
    "restart")
        echo "Restarting services..."
        docker compose -f "$COMPOSE_FILE" down
        docker compose -f "$COMPOSE_FILE" up -d
        ;;
    "up")
        echo "Starting services..."
        docker compose -f "$COMPOSE_FILE" up -d
        ;;
    *)
        echo "Unknown action: $ACTION"
        cd "$ORIGINAL_DIR"
        exit 1
        ;;
esac

cd "$ORIGINAL_DIR"

if [ $? -eq 0 ]; then
    echo "Operation completed successfully."
else
    echo "Operation failed."
    exit 1
fi