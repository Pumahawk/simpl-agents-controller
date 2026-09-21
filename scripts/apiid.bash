#!/bin/bash
set -e
yq '
.paths |
  to_entries[] |
  .key as $path |
  .value |
  to_entries[] |
  .key as $method |
  .value.operationId as $id |
  .value.operationId | line as $line |
  "\($line) \($id) \($method) \($path)"' | column --table
