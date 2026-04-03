#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEPLOY_DIR="$BASE_PATH/deploy"
ENV_FILE="$DEPLOY_DIR/.env"
ENV_EXAMPLE_FILE="$DEPLOY_DIR/.env.example"

if [ ! -t 0 ]; then
    echo "Interactive terminal is required."
    echo "Fallback: copy deploy/.env.example to deploy/.env then run bash deploy/script/prepare.sh"
    exit 1
fi

if [ ! -f "$ENV_EXAMPLE_FILE" ]; then
    echo "Missing required file: $ENV_EXAMPLE_FILE"
    exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
    cp "$ENV_EXAMPLE_FILE" "$ENV_FILE"
fi

set -a
# shellcheck disable=SC1090
source "$ENV_EXAMPLE_FILE"
# shellcheck disable=SC1090
source "$ENV_FILE"
set +a

update_env_key() {
    local key="$1"
    local value="$2"
    local tmp_file=""

    tmp_file="$(mktemp)"
    awk -v key="$key" -v val="$value" '
        BEGIN { updated = 0 }
        $0 ~ "^" key "=" {
            print key "=" val
            updated = 1
            next
        }
        { print }
        END {
            if (!updated) {
                print key "=" val
            }
        }
    ' "$ENV_FILE" > "$tmp_file"
    mv "$tmp_file" "$ENV_FILE"
}

prompt_value() {
    local key="$1"
    local label="$2"
    local fallback="$3"
    local current_value=""
    local input=""

    current_value="${!key:-$fallback}"
    read -r -p "$label [$current_value]: " input
    if [ -z "$input" ]; then
        input="$current_value"
    fi

    printf -v "$key" '%s' "$input"
    update_env_key "$key" "$input"
}

echo "== Yoresee Deploy Config Init =="
echo "Press Enter to keep the default value in brackets."

# Host exposed ports
prompt_value NGINX_HTTP_PORT "Host port for Nginx HTTP" "8080"
prompt_value NGINX_HTTPS_PORT "Host port for Nginx HTTPS" "8443"
prompt_value BACKEND_GRPC_HOST_PORT "Host port for backend gRPC" "9090"
prompt_value BACKEND_DEBUG_HOST_PORT "Host port for backend debug (dev only)" "2345"
prompt_value POSTGRES_HOST_PORT "Host port for Postgres (dev only)" "5432"
prompt_value REDIS_HOST_PORT "Host port for Redis (dev only)" "6379"
prompt_value RABBITMQ_AMQP_HOST_PORT "Host port for RabbitMQ AMQP (dev only)" "5672"
prompt_value CONSUL_HOST_PORT "Host port for Consul UI/API" "8500"
prompt_value ELASTICSEARCH_HOST_PORT "Host port for Elasticsearch (dev only)" "9200"

# Credentials and app secrets
prompt_value POSTGRES_USER "Postgres user" "root"
prompt_value POSTGRES_PASSWORD "Postgres password" "your_password"
prompt_value POSTGRES_DB "Postgres database" "yoresee_doc_db"
prompt_value REDIS_PASSWORD "Redis password" "your_redis_password"
prompt_value RABBITMQ_DEFAULT_USER "RabbitMQ user" "guest"
prompt_value RABBITMQ_DEFAULT_PASS "RabbitMQ password" "guest"
prompt_value MINIO_ROOT_USER "MinIO root user" "minioadmin"
prompt_value MINIO_ROOT_PASSWORD "MinIO root password" "minioadmin"
prompt_value CONSUL_ROOT_TOKEN "Consul root token" "yoresee_doc_root_token"
prompt_value BACKEND_INTERNAL_RPC_KEY "Internal RPC key" "yoresee_doc_internal_key"
prompt_value JWT_SECRET "JWT secret" "yoresee_doc_jwt_secret_key"

# App behavior
prompt_value VITE_API_BASE_URL "Frontend API base URL" "http://localhost:${NGINX_HTTP_PORT}"
prompt_value VITE_GRPC_WEB_ENDPOINT "Frontend gRPC-web endpoint" "/grpc"
prompt_value DIRTY_DOC_NOTIFY_THRESHOLD "Dirty doc notify threshold" "5"

# Keep MinIO browser redirect bound to API base URL to avoid host drift in deployments.
MINIO_BROWSER_REDIRECT_URL="${VITE_API_BASE_URL%/}/minio"
update_env_key MINIO_BROWSER_REDIRECT_URL "$MINIO_BROWSER_REDIRECT_URL"
echo "MinIO browser redirect URL (derived): $MINIO_BROWSER_REDIRECT_URL"

prompt_value ELASTICSEARCH_ENABLED "Enable Elasticsearch" "true"
prompt_value ELASTICSEARCH_INDEX_PREFIX "Elasticsearch index prefix" "yoresee_doc"
prompt_value ELASTICSEARCH_USERNAME "Elasticsearch username (optional)" ""
prompt_value ELASTICSEARCH_PASSWORD "Elasticsearch password (optional)" ""

echo "Written: $ENV_FILE"

bash "$SCRIPT_DIR/prepare.sh"

echo "Initialization completed."
echo "You can start services with: bash deploy/script/start.sh dev up"
