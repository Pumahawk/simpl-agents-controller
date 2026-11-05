#!/bin/bash

set -e

main() {
  jq \
    --unbuffered \
  -R \
  -c \
  '
  .
  | capture("(?<pod>^[^ ]*) (?<ktime>[^ ]*) (?<message>.*)")
  | {
      ktime: .ktime,
      micro: .pod,
      message: .message,
    }
  | . as $line
  | try ( .message | fromjson | $line * . ) catch $line
  '
}

main "$@"
