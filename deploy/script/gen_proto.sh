#!/bin/sh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
PROTO_DIR="$ROOT_DIR/proto"
GO_OUT_DIR="$BACKEND_DIR/pkg/gen"
FRONTEND_OUT_DIR="$ROOT_DIR/frontend/src/gen"

export PATH="$(go env GOPATH)/bin:$PATH"

mkdir -p "$GO_OUT_DIR"

protoc -I "$PROTO_DIR" \
  --go_out="$GO_OUT_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$GO_OUT_DIR" --go-grpc_opt=paths=source_relative \
  "$PROTO_DIR/yoresee_doc/v1/yoresee_doc.proto"

ES_BIN="$ROOT_DIR/frontend/node_modules/.bin/protoc-gen-es"
CONNECT_ES_BIN="$ROOT_DIR/frontend/node_modules/.bin/protoc-gen-connect-es"

if [ -x "$ES_BIN" ] && [ -x "$CONNECT_ES_BIN" ]; then
  mkdir -p "$FRONTEND_OUT_DIR"
  protoc -I "$PROTO_DIR" \
    --plugin=protoc-gen-es="$ES_BIN" \
    --plugin=protoc-gen-connect-es="$CONNECT_ES_BIN" \
    --es_out="$FRONTEND_OUT_DIR" --es_opt=target=js,import_extension=.js \
    --connect-es_out="$FRONTEND_OUT_DIR" --connect-es_opt=target=js,import_extension=.js \
    "$PROTO_DIR/yoresee_doc/v1/yoresee_doc.proto"
else
  echo "protoc-gen-es or protoc-gen-connect-es not found; skipping frontend connect-es code generation" >&2
fi
