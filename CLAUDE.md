# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

远楒文档 (Yoresee Doc) — a real-time collaborative document platform with Go backend, Vue 3 frontend, and Yjs-based collaboration.

## Commands

### Configuration & startup

```bash
# Interactive config (recommended first time)
bash deploy/script/configure.sh

# Non-interactive
cp deploy/.env.example deploy/.env
bash deploy/script/prepare.sh

# Dev mode
bash deploy/script/start.sh dev up

# Rebuild and restart
bash deploy/script/start.sh dev rebuild
bash deploy/script/start.sh dev restart

# Production mode
bash deploy/script/start.sh release up
```

`configure.sh` collects secrets/ports interactively and writes `deploy/.env`, then runs `prepare.sh`.
`prepare.sh` renders templates (`*.tmpl` → config files) using `deploy/.env` values.

### Proto generation

```bash
bash deploy/script/gen_proto.sh
```

Generates to `backend/pkg/gen`, `collab-go/pkg/gen`, `frontend/src/gen`, `collab/src/gen`.
Requires host `protoc`, `protoc-gen-es`, `protoc-gen-connect-es`, and `grpc_tools_node_protoc_plugin`.

### Frontend (standalone)

```bash
cd frontend && npm run dev     # Vite dev server
cd frontend && npm run build   # Production build
```

## Architecture

### Service topology

| Service | Dir | Role |
|---|---|---|
| `backend` | `backend/` | Go Connect/gRPC API server + 3 workers |
| `frontend` | `frontend/` | Vue 3 + Vite SPA |
| `collab-core` | `collab/` | Node.js Yjs document sync core |
| `collab-gateway` | `collab-go/` | Go WebSocket gateway (auth + proxy to collab-core) |
| Infrastructure | `deploy/` | Postgres, Redis, RabbitMQ, Consul, MinIO, Elasticsearch, Nginx |

### Request flow

1. Browser → Nginx (`:8080`)
2. Nginx `/` → Frontend; `/grpc/` → Backend gRPC-web
3. Backend accesses Postgres/Redis/MinIO/Consul/ES/RabbitMQ as needed

### Real-time collaboration flow

1. Browser → Nginx `/ws/doc/{docId}` → `collab-go`
2. `collab-go` validates JWT → proxies to `collab-core` (Node.js)
3. `collab-core` maintains Yjs doc state in memory + Redis

### Async event flow

1. `collab-core` marks dirty docs in Redis, emits events
2. `snapshot-worker` consumes `collab.dirty_docs` from RabbitMQ, persists snapshots to DB
3. `search-sync-worker` consumes `search.sync.document`, writes to Elasticsearch
4. `notification-worker` consumes `notification.create`, writes notification records

### Backend structure

- `backend/cmd/` — Entry points: `main.go` (API server), `migrate/`, `db_init/`, `es_init/`, plus 3 workers
- `backend/internal/bootstrap/` — Dependency wiring (config → DB → Redis → Consul → MinIO → ES → MQ → repos)
- `backend/internal/config/` — Viper-based config loading
- `backend/internal/service/` — Business logic layer
- `backend/internal/repository/` — Data access (GORM + PostgreSQL)
- `backend/internal/model/` — DB models
- `backend/internal/dto/` — Request/response DTOs
- `backend/internal/transport/connectserver/` — Connect RPC / gRPC-web server
- `backend/internal/domain_event/` — Domain event definitions and publishers
- `backend/internal/mapper/` — Model ↔ DTO mapping
- `backend/internal/auth/` — JWT authentication
- `backend/internal/cache/` — Redis caching layer
- `backend/internal/media/` — MinIO object storage integration
- `backend/internal/search/` — Elasticsearch integration
- `backend/internal/middleware/` — gRPC interceptors

MQ abstraction supports both Redis Pub/Sub and RabbitMQ backends. For dirty-doc events, use `rabbitmq` or `both` (Redis Pub/Sub lacks true group-consumption semantics).

### Frontend structure

- `frontend/src/views/` — Page-level components: `auth/`, `document/`, `workspace/`, `knowledge-base/`, `template/`, `manage/`, `user/`
- `frontend/src/components/` — Shared UI components
- `frontend/src/composables/` — Vue composables (reusable stateful logic)
- `frontend/src/services/` — API client layer (Connect RPC / gRPC-web calls)
- `frontend/src/gen/` — Generated protobuf/Connect ES code
- `frontend/src/router/` — Vue Router config
- `frontend/src/store/` — Pinia stores
- `frontend/src/i18n/` — vue-i18n localization

Key libraries: Element Plus (UI), TipTap (rich text editor), Yjs (CRDT collaboration), Markmap (mind maps), Vditor (Markdown editor).

### Health checks & graceful shutdown

All services expose `/health`, `/readyz`, `/livez`. During draining, `/readyz` returns not-ready. Backend and collab services support graceful shutdown with configurable timeouts.

### Config rendering

Template files (`*.tmpl`) are rendered by `prepare.sh` using envsubst-style substitution from `deploy/.env`. Key templates: `backend/config.toml.tmpl`, `frontend/nginx.conf.tmpl`, `deploy/nginx/nginx.conf.tmpl`, `deploy/nginx/conf.d/default.conf.tmpl`.
