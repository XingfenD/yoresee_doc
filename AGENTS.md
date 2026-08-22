# AGENTS.md

## Proto generation is external

Proto definitions live in `../yoresee_doc_proto`; this repo only consumes published packages (Go: `github.com/XingfenD/yoresee_doc_proto`; npm: `@yoresee/doc-proto-connect`, `@yoresee/doc-proto-grpc`).

Do NOT add `.proto` files or run `protoc` here. To change proto: edit in the proto repo, run its `scripts/gen.sh`, commit, bump `proto/VERSION`, `git tag v$(cat proto/VERSION)`, push the tag. For local testing: `bash deploy/script/link_proto.sh link ../yoresee_doc_proto` (restore with `unlink`).

## Two Go modules

`backend/` and `collab-go/` each have their own `go.mod`. Add dependencies to the correct module.

## Config rendering

`deploy/.env` → `*.tmpl` rendered by `bash deploy/script/prepare.sh`; rerun it after changing `.env`. Edit the `.tmpl` files, never the rendered outputs.

## Docker dev mode

```bash
bash deploy/script/start.sh dev up      # first time or after config change
bash deploy/script/start.sh dev rebuild # after code changes
```

## Coding style (Go)

`gofmt -s`; follow the existing layering (`transport/connectserver` → `service` → `repository` → `model`).

## Testing

Adjacent `*_test.go` with stdlib `testing`; stub Postgres/Redis/MQ, never depend on live services. Run `go test -race ./...` in the changed module before submitting.

## Branches

Create a `feat`/`fix`/`refactor`/`docs` branch for changes; do not commit directly to `master`.

## Commits

Conventional Commits (`fix: ...`, `refactor: ...`).

## Changelog

Record user-visible changes in `docs/CHANGELOG`, Keep-a-Changelog style. Newest version on top. The `[x.y.z]` header version must match `backend/pkg/constant/version.go` — bump the constant alongside the CHANGELOG entry. New-entry format:

```markdown
## [x.y.z] - YYYY-MM-DD

### Added / 新增
### Changed / 变更
### Fixed / 修复
### Removed / 移除

- One line per change, bilingual (zh + en) to match repo docs.
```

Pending work that is not yet shipped goes in `docs/TODO.md`, not here.

## Gotchas

- `frontend/package.json` has no `lint`/`typecheck` scripts.
- No CI pipeline (`.github/workflows`) exists.
