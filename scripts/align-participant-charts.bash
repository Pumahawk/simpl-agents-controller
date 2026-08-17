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
    git fetch --prune -m origin participant
    if >/dev/null git diff --exit-code origin/feature/chart-participant-develop participant/develop; then
      >&2 echo "No differences between origin/feature/chart-participant-develop participant/develop."
      exit 0
    fi

    CM_PR_DEV="$(git commit-tree participant/develop^{tree} -p origin/feature/chart-participant-develop -m "update $(git rev-list -n1 participant/develop)")"
    TR_MR="$(git merge-tree origin/develop "$CM_PR_DEV")"
    CM_DEV="$(git commit-tree "$TR_MR" -p origin/develop -p "$CM_PR_DEV" -m "update $(git rev-list -n1 participant/develop)")"

    git push origin "$CM_PR_DEV":refs/heads/feature/chart-participant-develop "$CM_DEV":refs/heads/feature/update-develop
EOF
  )"
}

align_chart "$@"
