#!/bin/bash

set -e

NAMESPACE=""
SINCE="10m"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --namespace|-n)
      NAMESPACE="$2"
      shift 2
      ;;
    --since|-s)
      SINCE="$2"
      shift 2
      ;;
    --help|-h)
      cat <<EOF
Usage: $0 --namespace <namespace> [--since <duration>]

Options:
  -n, --namespace   Kubernetes namespace (required)
  -s, --since       Log time range (default: 10m)

Examples:
  $0 -n iaadsdevauthority
  $0 -n iaadsdevauthority -s 1h
EOF
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

: "${NAMESPACE:?Missing --namespace}"

PREFIX="$(openssl rand -hex 5)"

xargs -L1 -P4 bash -c '
kubectl logs -n "$1" --prefix --timestamps --since "$2" "$3" \
  > "/tmp/log-'"$PREFIX"'-${3#*/}"
' -- "$NAMESPACE" "$SINCE"

cat /tmp/log-"$PREFIX"-*
