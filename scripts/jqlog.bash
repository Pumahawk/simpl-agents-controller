#!/bin/bash
jq --unbuffered -r -R '. as $line | try (fromjson | "\(.timestamp) \(.micro // "") \(.level) \(.message) \(."error.stack_trace" // "")") catch $line'
