#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
PROTO_DIR="${1:-$ROOT_DIR/../yoresee_doc_proto}"

if [ ! -d "$PROTO_DIR" ]; then
    echo "Error: proto directory not found: $PROTO_DIR"
    echo "Usage: $0 [path-to-proto-repo]"
    exit 1
fi

VERSION=$(cat "$PROTO_DIR/proto/VERSION")
echo "Updating proto dependencies to v$VERSION..."

echo "  backend..."
(cd "$ROOT_DIR/backend" && GOPROXY=https://goproxy.cn,direct go get "github.com/XingfenD/yoresee_doc_proto@v$VERSION" && go mod tidy)

echo "  collab-go..."
(cd "$ROOT_DIR/collab-go" && GOPROXY=https://goproxy.cn,direct go get "github.com/XingfenD/yoresee_doc_proto@v$VERSION" && go mod tidy)

echo "  frontend..."
(cd "$ROOT_DIR/frontend" && npm install "@yoresee/doc-proto-connect@$VERSION")

echo "  collab..."
(cd "$ROOT_DIR/collab" && npm install "@yoresee/doc-proto-grpc@$VERSION")

echo "Done! Proto dependencies updated to v$VERSION"
