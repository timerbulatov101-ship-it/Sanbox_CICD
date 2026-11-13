#!/bin/sh

set -e

IMAGE=${1:-sandbox-frontend:latest}

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

docker build -t "$IMAGE" "$SCRIPT_DIR"