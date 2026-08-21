# AGENTS.md

## Proto generation is external

Proto definitions live in a separate repo: `../yoresee_doc_proto`. This repo consumes published packages:
- Go: `github.com/XingfenD/yoresee_doc_proto` (backend, collab-go)
- Frontend: `@yoresee/doc-proto-connect` (npm)
- Collab: `@yoresee/doc-proto-grpc` (npm)

To modify proto: edit in `../yoresee_doc_proto`, run `bash scripts/gen.sh` there, bump `proto/VERSION`, commit and tag. The publish workflow auto-opens a PR in this repo to update dependencies.

Do NOT add proto files to this repo. Do NOT run protoc here.

## Two Go modules

- `backend/` — module `github.com/XingfenD/yoresee_doc`
- `collab-go/` — module `github.com/XingfenD/yoresee_doc/collab-go`

Each has its own `go.mod`. When adding dependencies, update the correct module.

## Config rendering

`deploy/.env` → `*.tmpl` files are rendered to config files by `bash deploy/script/prepare.sh`. After changing `deploy/.env`, run `prepare.sh` before restarting services.

Templates: `backend/config.toml.tmpl`, `frontend/nginx.conf.tmpl`, `deploy/nginx/nginx.conf.tmpl`, `deploy/nginx/conf.d/default.conf.tmpl`, `deploy/redis/redis.conf.tmpl`, `deploy/rabbitmq/rabbitmq.conf.tmpl`.

## Docker dev mode

```bash
bash deploy/script/start.sh dev up      # first time or after config change
bash deploy/script/start.sh dev rebuild # after code changes
```

All app services mount source via volumes. Code changes are picked up automatically (Go: `go run` with file watching, Node: Vite HMR / `node --watch`).

## Frontend has no lint or typecheck

`frontend/package.json` has no `lint` or `typecheck` scripts. If adding TypeScript or ESLint, configure scripts first.

## No CI pipeline

No `.github/workflows` exist. If adding CI, remember:
- Backend and collab-go are separate Go modules, need separate jobs
- Frontend needs `npm install` before build (no lockfile in repo root)
- Integration tests need Postgres/Redis/RabbitMQ/ES as service containers

## Backend entry points

- `backend/cmd/main.go` — API server
- `backend/cmd/migrate/` — DB migrations (run automatically on startup)
- `backend/cmd/db_init/` — DB initialization
- `backend/cmd/es_init/` — Elasticsearch index setup
- `backend/cmd/snapshot-worker/` — persists Yjs snapshots
- `backend/cmd/notification-worker/` — writes notification records
- `backend/cmd/search-sync-worker/` — syncs docs to Elasticsearch

## MQ abstraction

Backend supports both Redis Pub/Sub and RabbitMQ. For `collab.dirty_docs` events, use `rabbitmq` or `both` (Redis Pub/Sub lacks group-consumption semantics).
