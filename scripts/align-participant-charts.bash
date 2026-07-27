#!/bin/bash

function align_chart() {
  local chname="${1?Missing chart name}"
  local repo="simpl-repo/$chname"
  if [ ! -d "$repo" ]; then
    >&2 echo "Not found repo \"$repo\""
    return 1
  fi

  env -C "$repo" bash -xec "$(
    cat <<'EOF'
    git add -A
    git stash
    git fetch --multiple --prune origin participant 
    git checkout -d
    git branch -D feature/chart-participant-develop || true
    git checkout participant/develop
    commit="$(git rev-list -n1 HEAD)"
    git reset --soft origin/feature/chart-participant-develop
    git commit -m "update chart participant $commit"
    git push HEAD:feature/chart-participant-develop
    git checkout origin/develop
    git branch -d chart-participant-update || true
    git checkout -b chart-participant-update
    git merge --no-edit origin/feature/chart-participant-develop
EOF
  )"
}

align_chart "$@"
