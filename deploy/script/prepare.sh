#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEPLOY_DIR="$BASE_PATH/deploy"
ENV_EXAMPLE_FILE="$DEPLOY_DIR/.env.example"
ENV_FILE="$DEPLOY_DIR/.env"

if [ ! -f "$ENV_EXAMPLE_FILE" ]; then
    echo "Missing required file: $ENV_EXAMPLE_FILE"
    exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
    echo "deploy/.env not found, creating from .env.example"
    cp "$ENV_EXAMPLE_FILE" "$ENV_FILE"
fi

set -a
# shellcheck disable=SC1090
source "$ENV_EXAMPLE_FILE"
# shellcheck disable=SC1090
source "$ENV_FILE"
set +a

prepare_target_file() {
    local target="$1"
    if [ -d "$target" ]; then
        local backup="${target}.backup.$(date +%s)"
        mv "$target" "$backup"
        echo "Detected directory at $target, moved to $backup"
    fi
}

render_template() {
    local template="$1"
    local output="$2"
    local tmp_file=""

    if [ ! -f "$template" ]; then
        echo "Template not found: $template"
        return 1
    fi

    mkdir -p "$(dirname "$output")"
    prepare_target_file "$output"

    tmp_file="$(mktemp)"
    perl -pe 's/\$\{([A-Z0-9_]+)\}/defined $ENV{$1} ? $ENV{$1} : ""/ge' "$template" > "$tmp_file"
    mv "$tmp_file" "$output"
}

TEMPLATE_MAPPINGS=(
    "$BASE_PATH/backend/config.toml.tmpl:$BASE_PATH/backend/config.toml"
    "$BASE_PATH/frontend/nginx.conf.tmpl:$BASE_PATH/frontend/nginx.conf"
    "$DEPLOY_DIR/nginx/nginx.conf.tmpl:$DEPLOY_DIR/nginx/nginx.conf"
    "$DEPLOY_DIR/nginx/conf.d/default.conf.tmpl:$DEPLOY_DIR/nginx/conf.d/default.conf"
    "$DEPLOY_DIR/redis/redis.conf.tmpl:$DEPLOY_DIR/redis/redis.conf"
    "$DEPLOY_DIR/rabbitmq/rabbitmq.conf.tmpl:$DEPLOY_DIR/rabbitmq/rabbitmq.conf"
)

echo "Rendering configuration files from templates..."
for mapping in "${TEMPLATE_MAPPINGS[@]}"; do
    template="${mapping%%:*}"
    output="${mapping#*:}"
    render_template "$template" "$output"
    rel_output="${output#$BASE_PATH/}"
    if [ "$rel_output" = "$output" ]; then
        rel_output="$output"
    fi
    echo "  - $rel_output"
done

echo "Configuration preparation completed!"
