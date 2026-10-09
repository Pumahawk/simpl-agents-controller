#!/bin/bash

function align_chart() {
  for remote in "$@"; do
    if ! git remote get-url "$remote" > /dev/null; then
      >&2 echo "Not found remote $remote"
      return 1
    fi
  done
  
  git fetch origin
  for remote in "$@"; do
    git fetch "$remote"
    if >/dev/null git diff --exit-code "$remote/feature/chart-participant-develop" origin/develop; then
      >&2 echo "No differences between $remote/feature/chart-participant-develop origin/develop."
      break
    fi
    local CM_PR_DEV="$(git commit-tree origin/develop^{tree} -p "$remote/feature/chart-participant-develop" -m "update $(git rev-list -n1 origin/develop)")"
    local TR_MR="$(git merge-tree "$remote/develop" "$CM_PR_DEV")"
    local CM_DEV="$(git commit-tree "$TR_MR" -p "$remote/develop" -p "$CM_PR_DEV" -m "update $(git rev-list -n1 origin/develop)")"

    git push "$remote" "$CM_PR_DEV":refs/heads/feature/chart-participant-develop "$CM_DEV":refs/heads/feature/update-develop
  done
}

align_chart "$@"
