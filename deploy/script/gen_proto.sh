#!/bin/sh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PROTO_REPO="${PROTO_REPO_PATH:-$ROOT_DIR/../yoresee_doc_proto}"

if [ ! -f "$PROTO_REPO/scripts/gen.sh" ]; then
  echo "error: proto repo not found at $PROTO_REPO" >&2
  echo "set PROTO_REPO_PATH to the proto repo location" >&2
  exit 1
fi

echo "generating proto code in $PROTO_REPO ..."
(cd "$PROTO_REPO" && bash scripts/gen.sh)
echo "done. consumers use 'go mod download' / 'npm install' to pick up changes."
