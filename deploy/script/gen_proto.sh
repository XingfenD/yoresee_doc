#!/bin/sh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
PROTO_DIR="$ROOT_DIR/proto"
GO_OUT_DIR="$BACKEND_DIR/pkg/gen"
FRONTEND_OUT_DIR="$ROOT_DIR/frontend/src/gen"
COLLAB_OUT_DIR="$ROOT_DIR/collab/gen"

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

GRPC_PLUGIN="$(which grpc_tools_node_protoc_plugin 2>/dev/null || echo "")"

if [ -n "$GRPC_PLUGIN" ] && [ -x "$GRPC_PLUGIN" ]; then
  mkdir -p "$COLLAB_OUT_DIR"

  protoc -I "$PROTO_DIR" \
    --plugin=protoc-gen-grpc="$GRPC_PLUGIN" \
    --grpc_out=grpc_js:"$COLLAB_OUT_DIR" \
    --js_out=import_style=commonjs,binary:"$COLLAB_OUT_DIR" \
    "$PROTO_DIR/yoresee_doc/v1/yoresee_doc.proto"

  echo "Generated gRPC code for collab at $COLLAB_OUT_DIR"
else
  echo "grpc_tools_node_protoc_plugin not found. Install with: npm install -g grpc-tools" >&2
  exit 1
fi
