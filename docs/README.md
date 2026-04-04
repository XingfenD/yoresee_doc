# Yoresee Doc

Yoresee Doc is a collaborative document platform built with a Go backend, Vue 3 frontend, and a Yjs-based real-time collaboration stack.

## Runtime Topology

Core app services:
- `frontend`: Vue 3 (Vite)
- `backend`: Go Connect RPC/gRPC service
- `collab-core`: Node.js Yjs collaboration core
- `collab`: `collab-go` WebSocket gateway
- `snapshot-worker`: document snapshot sync worker
- `notification-worker`: notification event consumer
- `search-sync-worker`: Elasticsearch sync worker

Infrastructure services:
- `postgres`
- `redis`
- `rabbitmq`
- `consul`
- `minio`
- `elasticsearch`
- `nginx` (single ingress)

## Request and Event Flow

Main request path:
1. Browser -> Nginx
2. Nginx `/` -> Frontend
3. Nginx `/grpc/` -> Backend gRPC-web endpoint
4. Backend -> Postgres/Redis/MinIO/Consul/Elasticsearch/MQ as needed

Realtime collaboration path:
1. Browser -> Nginx `/ws/doc/{docId}`
2. Nginx -> `collab-go`
3. `collab-go` validates JWT and proxies to `collab-core`
4. `collab-core` maintains Yjs doc state in memory and Redis

Async event path:
1. `collab-core` marks dirty docs in Redis and publishes dirty-doc events
2. `snapshot-worker` consumes `collab.dirty_docs` (RabbitMQ group mode) and also scans dirty set periodically
3. `snapshot-worker` fetches `/internal/yjs/doc-snapshot/{docId}` from `collab-core`, writes snapshot/content to DB, and emits search-sync event when content changed
4. `search-sync-worker` consumes `search.sync.document` and upserts Elasticsearch (when ES is enabled)
5. `notification-worker` consumes `notification.create` and writes notification records

## Health and Graceful Shutdown

Backend probes:
- `/health`
- `/readyz`
- `/livez`

Collab probes:
- `collab-core`: `/health`, `/readyz`, `/livez`
- `collab-go`: `/health`, `/readyz`, `/livez`

Shutdown behavior:
- Backend and collab services support draining and graceful shutdown.
- Readiness becomes `not_ready` during draining.

## Deployment Modes

`deploy/docker-compose.dev.yml`:
- Source-mounted development containers
- Backend debug port exposed (`2345`)
- Postgres/Redis/RabbitMQ/Elasticsearch ports exposed to host
- `start.sh dev ...` auto-runs `deploy/script/gen_proto.sh`

`deploy/docker-compose.yml`:
- Production-style images
- Fewer host-exposed ports
- Backend/worker binaries are prebuilt in image stages

## Dev Prerequisites

`start.sh dev ...` runs `deploy/script/gen_proto.sh` on the host machine before docker compose starts.

Required on host:
- `docker` + `docker compose`
- `go` + `protoc` (for Go stubs)
- frontend plugins in `frontend/node_modules` (`protoc-gen-es`, `protoc-gen-connect-es`)
- `grpc_tools_node_protoc_plugin` in `PATH` (for Node gRPC stubs used by `collab`)

## Configuration Workflow

Single source of truth:
- `deploy/.env`
- template: `deploy/.env.example`

Interactive init/update:
```bash
bash deploy/script/configure.sh
```

What `configure.sh` does:
- prompts for key ports/secrets/env values
- writes/updates `deploy/.env`
- derives `MINIO_BROWSER_REDIRECT_URL` from `VITE_API_BASE_URL`
- calls `prepare.sh` automatically

Render generated config files:
```bash
bash deploy/script/prepare.sh
```

Generated files:
- `backend/config.toml`
- `frontend/nginx.conf`
- `deploy/nginx/nginx.conf`
- `deploy/nginx/conf.d/default.conf`
- `deploy/redis/redis.conf`
- `deploy/rabbitmq/rabbitmq.conf`

Template files:
- `backend/config.toml.tmpl`
- `frontend/nginx.conf.tmpl`
- `deploy/nginx/nginx.conf.tmpl`
- `deploy/nginx/conf.d/default.conf.tmpl`
- `deploy/redis/redis.conf.tmpl`
- `deploy/rabbitmq/rabbitmq.conf.tmpl`

## Quick Start

Recommended (interactive):
```bash
bash deploy/script/configure.sh
bash deploy/script/start.sh dev up
```

Non-interactive:
```bash
cp deploy/.env.example deploy/.env
bash deploy/script/prepare.sh
bash deploy/script/start.sh dev up
```

Release mode:
```bash
bash deploy/script/start.sh release up
```

Backend image startup behavior:
- backend container runs `migrate`, `db_init`, `es_init`, then starts `cmd/main`.

Other actions:
```bash
bash deploy/script/start.sh dev rebuild
bash deploy/script/start.sh dev restart
bash deploy/script/start.sh dev clear
```

## Public Entry Points

Default dev ports:
- App UI: `http://localhost:8080`
- gRPC-web: `http://localhost:8080/grpc/`
- Collaboration WS: `ws://localhost:8080/ws/doc/{docId}?token={jwt}`
- MinIO console (via Nginx): `http://localhost:8080/minio/`
- Object public path (via Nginx): `http://localhost:8080/storage/...`
- RabbitMQ management (via Nginx): `http://localhost:8080/rabbitmq/`

Direct infra ports in dev:
- Postgres: `localhost:5432`
- Redis: `localhost:6379`
- RabbitMQ AMQP: `localhost:5672`
- Consul: `localhost:8500`
- Elasticsearch: `localhost:9200`

## Message Queue Notes

- Worker consumers are currently configured to use RabbitMQ backend.
- MQ interface supports Redis and RabbitMQ implementations.
- Redis Pub/Sub does not provide true consumer-group semantics.
- For dirty-doc pipeline, keep `DIRTY_DOC_MQ` as `rabbitmq` or `both` if snapshot-worker should consume events from RabbitMQ.

## Build and Protobuf

Manual protobuf generation:
```bash
bash deploy/script/gen_proto.sh
```

Generated targets:
- `backend/pkg/gen`
- `collab-go/pkg/gen`
- `frontend/src/gen`
- `collab/src/gen`

## Repository Structure

- `backend`: API service, workers, repositories, storage integration
- `frontend`: Vue application
- `collab`: Node.js collaboration core
- `collab-go`: Go WebSocket gateway for collaboration traffic
- `proto`: protobuf contracts
- `deploy`: compose files, scripts, infra templates
- `docs`: project docs
