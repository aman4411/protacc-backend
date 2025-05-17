#!/bin/bash

# Check if command is provided
if [ -z "$1" ]; then
    echo "Usage: ./migrate.sh [up|down|create|force]"
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

# Check if DATABASE_URL is set
if [ -z "${DATABASE_URL}" ]; then
    echo "Error: DATABASE_URL environment variable is not set"
    echo "Please ensure your .env file contains DATABASE_URL or set it manually"
    echo "Example: DATABASE_URL=postgres://username:password@localhost:5432/dbname"
    exit 1
fi


# Migration directory
MIGRATION_DIR="db/migrations"

case "$1" in
    "up")
        migrate -database "${DATABASE_URL}" -path ${MIGRATION_DIR} up
        ;;
    "down")
        migrate -database "${DATABASE_URL}" -path ${MIGRATION_DIR} down
        ;;
    "create")
        if [ -z "$2" ]; then
            echo "Usage: ./migrate.sh create <migration_name>"
            exit 1
        fi
        migrate create -ext sql -dir ${MIGRATION_DIR} -seq "$2"
        ;;
    "force")
        if [ -z "$2" ]; then
            echo "Usage: ./migrate.sh force <version>"
            exit 1
        fi
        migrate -database "${DATABASE_URL}" -path ${MIGRATION_DIR} force "$2"
        ;;
    *)
        echo "Invalid command. Use 'up', 'down', 'create', or 'force'"
        exit 1
        ;;
esac 