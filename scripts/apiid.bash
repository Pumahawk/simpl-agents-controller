#!/bin/bash
set -e
file="$1"

if [ -z "$file" ]; then
  file="-"
fi

FILE="$file" yq '
.paths |
  to_entries[] |
  .key as $path |
  .value |
  to_entries[] |
  .key as $method |
  .value.operationId as $id |
  .value.operationId | line as $line |
  "\(strenv(FILE)):\($line) \($id) \($method) \($path)"' "$@" | column --table -l4
