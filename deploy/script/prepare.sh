#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Preparing configuration files from examples..."

if [ -f "$BASE_PATH/backend/config.example.toml" ]; then
    cp "$BASE_PATH/backend/config.example.toml" "$BASE_PATH/backend/config.toml"
    echo "Copied backend/config.example.toml to backend/config.toml"
else
    echo "Warning: $BASE_PATH/backend/config.example.toml not found"
fi

if [ -f "$BASE_PATH/frontend/nginx.example.conf" ]; then
    cp "$BASE_PATH/frontend/nginx.example.conf" "$BASE_PATH/frontend/nginx.conf"
    echo "Copied frontend/nginx.example.conf to frontend/nginx.conf"
else
    echo "Warning: $BASE_PATH/frontend/nginx.example.conf not found"
fi

if [ -f "$BASE_PATH/deploy/nginx/nginx.example.conf" ]; then
    cp "$BASE_PATH/deploy/nginx/nginx.example.conf" "$BASE_PATH/deploy/nginx/nginx.conf"
    echo "Copied deploy/nginx/nginx.example.conf to deploy/nginx/nginx.conf"
else
    echo "Warning: $BASE_PATH/deploy/nginx/nginx.example.conf not found"
fi

if [ -f "$BASE_PATH/deploy/nginx/conf.d/default.example.conf" ]; then
    cp "$BASE_PATH/deploy/nginx/conf.d/default.example.conf" "$BASE_PATH/deploy/nginx/conf.d/default.conf"
    echo "Copied deploy/nginx/conf.d/default.example.conf to deploy/nginx/conf.d/default.conf"
else
    echo "Warning: $BASE_PATH/deploy/nginx/conf.d/default.example.conf not found"
fi

if [ -f "$BASE_PATH/deploy/redis/redis.example.conf" ]; then
    cp "$BASE_PATH/deploy/redis/redis.example.conf" "$BASE_PATH/deploy/redis/redis.conf"
    echo "Copied deploy/redis/redis.example.conf to deploy/redis/redis.conf"
else
    echo "Warning: $BASE_PATH/deploy/redis/redis.example.conf not found"
fi

echo "Configuration preparation completed!"