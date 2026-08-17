#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(cd "$SCRIPT_DIR/../.." && pwd)"
STATE_FILE="$BASE_PATH/deploy/.link_proto.state"

BACKEND_DIR="$BASE_PATH/backend"
GOPROXY_DIR="$BASE_PATH/collab-go"
FRONTEND_DIR="$BASE_PATH/frontend"
COLLAB_DIR="$BASE_PATH/collab"

GO_MODULE="github.com/XingfenD/yoresee_doc_proto"
CONNECT_PKG="@yoresee/doc-proto-connect"
GRPC_PKG="@yoresee/doc-proto-grpc"
CONNECT_PKG_DIR="packages/connect"
GRPC_PKG_DIR="packages/grpc"

DEFAULT_PROTO_PATH="$(dirname "$BASE_PATH")/yoresee_doc_proto"

usage() {
    cat <<EOF
Usage: $0 <command> [proto-repo-path]

Commands:
  link [path]    point all consumers at a local proto repo clone
  unlink         restore published versions
  gen [path]     regenerate code in the local proto repo clone
  status         show whether consumers are linked

The proto repo path defaults to PROTO_REPO_PATH, then \$1, then ../yoresee_doc_proto.
EOF
}

resolve_proto_path() {
    if [ -n "${PROTO_REPO_PATH:-}" ]; then
        echo "$PROTO_REPO_PATH"
    elif [ -n "${1:-}" ]; then
        echo "$1"
    else
        echo "$DEFAULT_PROTO_PATH"
    fi
}

validate_proto_path() {
    local proto="$1"
    if [ ! -f "$proto/go.mod" ]; then
        echo "error: proto repo go.mod not found at $proto" >&2
        exit 1
    fi
    if [ ! -f "$proto/$CONNECT_PKG_DIR/package.json" ]; then
        echo "error: $CONNECT_PKG_DIR not found at $proto" >&2
        exit 1
    fi
    if [ ! -f "$proto/$GRPC_PKG_DIR/package.json" ]; then
        echo "error: $GRPC_PKG_DIR not found at $proto" >&2
        exit 1
    fi
}

npm_version() {
    local dir="$1"
    local pkg="$2"
    node -p "try { require('$dir/package.json').dependencies['$pkg'] } catch (e) { '' }"
}

npm_has_pkg() {
    local dir="$1"
    local pkg="$2"
    [ -n "$(npm_version "$dir" "$pkg")" ]
}

cmd_link() {
    if [ -f "$STATE_FILE" ]; then
        echo "error: already linked ($(grep '^PROTO_REPO_PATH=' "$STATE_FILE"))" >&2
        echo "run '$0 unlink' first" >&2
        exit 1
    fi

    local proto
    proto="$(resolve_proto_path "${1:-}")"
    validate_proto_path "$proto"

    local connect_version=""
    local grpc_version=""
    if npm_has_pkg "$FRONTEND_DIR" "$CONNECT_PKG"; then
        connect_version="$(npm_version "$FRONTEND_DIR" "$CONNECT_PKG")"
    fi
    if npm_has_pkg "$COLLAB_DIR" "$GRPC_PKG"; then
        grpc_version="$(npm_version "$COLLAB_DIR" "$GRPC_PKG")"
    fi

    {
        echo "PROTO_REPO_PATH=$proto"
        echo "CONNECT_VERSION=$connect_version"
        echo "GRPC_VERSION=$grpc_version"
    } > "$STATE_FILE"

    (cd "$BACKEND_DIR" && go mod edit -replace "$GO_MODULE=$proto")
    (cd "$GOPROXY_DIR" && go mod edit -replace "$GO_MODULE=$proto")
    echo "linked Go module -> $proto (backend, collab-go)"

    if [ -n "$connect_version" ]; then
        (cd "$FRONTEND_DIR" && npm pkg set "dependencies.$CONNECT_PKG=file:$proto/$CONNECT_PKG_DIR" && npm install)
        echo "linked $CONNECT_PKG -> file:$proto/$CONNECT_PKG_DIR"
    else
        echo "skip: $CONNECT_PKG not a dependency of frontend"
    fi

    if [ -n "$grpc_version" ]; then
        (cd "$COLLAB_DIR" && npm pkg set "dependencies.$GRPC_PKG=file:$proto/$GRPC_PKG_DIR" && npm install)
        echo "linked $GRPC_PKG -> file:$proto/$GRPC_PKG_DIR"
    else
        echo "skip: $GRPC_PKG not a dependency of collab"
    fi
}

cmd_unlink() {
    if [ ! -f "$STATE_FILE" ]; then
        echo "error: not linked" >&2
        exit 1
    fi
    # shellcheck disable=SC1091
    source "$STATE_FILE"

    (cd "$BACKEND_DIR" && go mod edit -dropreplace "$GO_MODULE")
    (cd "$GOPROXY_DIR" && go mod edit -dropreplace "$GO_MODULE")
    echo "unlinked Go module (backend, collab-go)"

    if [ -n "${CONNECT_VERSION:-}" ]; then
        (cd "$FRONTEND_DIR" && npm pkg set "dependencies.$CONNECT_PKG=$CONNECT_VERSION" && npm install)
        echo "restored $CONNECT_PKG -> $CONNECT_VERSION"
    fi

    if [ -n "${GRPC_VERSION:-}" ]; then
        (cd "$COLLAB_DIR" && npm pkg set "dependencies.$GRPC_PKG=$GRPC_VERSION" && npm install)
        echo "restored $GRPC_PKG -> $GRPC_VERSION"
    fi

    rm -f "$STATE_FILE"
}

cmd_gen() {
    local proto
    proto="$(resolve_proto_path "${1:-}")"
    validate_proto_path "$proto"
    if [ ! -f "$proto/scripts/gen.sh" ]; then
        echo "error: scripts/gen.sh not found at $proto" >&2
        exit 1
    fi
    (cd "$proto" && bash scripts/gen.sh)
}

cmd_status() {
    if [ -f "$STATE_FILE" ]; then
        # shellcheck disable=SC1091
        source "$STATE_FILE"
        echo "linked to $PROTO_REPO_PATH"
    else
        echo "not linked (using published versions)"
    fi
}

case "${1:-}" in
    link)   shift; cmd_link "${1:-}" ;;
    unlink) shift; cmd_unlink ;;
    gen)    shift; cmd_gen "${1:-}" ;;
    status) cmd_status ;;
    -h|--help|help) usage ;;
    *)      usage; exit 1 ;;
esac
