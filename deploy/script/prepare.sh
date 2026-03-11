#!/bin/bash

# Script to prepare configuration files for different environments

# Function to copy example config files to actual config files
copy_config_file() {
    local src_path="$1"
    local dest_path="$2"

    if [ -f "$src_path" ]; then
        if [ ! -f "$dest_path" ]; then
            echo "Copied $(basename "$src_path") to $(basename "$dest_path")"
        else
            echo "$(basename "$dest_path") already exists, overwriting..."
        fi
        cp "$src_path" "$dest_path"
    else
        echo "Warning: $src_path not found"
    fi
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Preparing configuration files from examples..."
echo "Base path: $BASE_PATH"

# Prepare backend config
copy_config_file "$BASE_PATH/backend/config.example.toml" "$BASE_PATH/backend/config.toml"

# Prepare frontend config
copy_config_file "$BASE_PATH/frontend/nginx.example.conf" "$BASE_PATH/frontend/nginx.conf"

# Prepare nginx config
copy_config_file "$BASE_PATH/deploy/nginx/nginx.example.conf" "$BASE_PATH/deploy/nginx/nginx.conf"

# Prepare nginx default config
copy_config_file "$BASE_PATH/deploy/nginx/conf.d/default.example.conf" "$BASE_PATH/deploy/nginx/conf.d/default.conf"

# Prepare redis config
copy_config_file "$BASE_PATH/deploy/redis/redis.example.conf" "$BASE_PATH/deploy/redis/redis.conf"

# Prepare rabbitmq config
copy_config_file "$BASE_PATH/deploy/rabbitmq/rabbitmq.example.conf" "$BASE_PATH/deploy/rabbitmq/rabbitmq.conf"

echo "Configuration preparation completed!"