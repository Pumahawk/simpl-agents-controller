#!/bin/bash
set -e

GIT_BRANCH_REMOTE="origin"
JSON_OUT="0"

function main() {
  while true; do
    case "$1" in
    --json)
      JSON_OUT="1"
      shift
      ;;
    --remote)
      GIT_BRANCH_REMOTE="${2?Missing remote parameter}"
      shift
      shift
      ;;
    *)
      break
      ;;
    esac
  done
  case $1 in
  mr | mergerequests)
    shift
    merge_requests "$@"
    ;;
  pips | pipelines)
    shift
    pipelines "$@"
    ;;
  pk | packages)
    shift
    packages "$@"
    ;;
  *)
    c_giteu "$@"
    ;;
  esac
}

function c_giteu() {
  path="${1?Missing path}"
  shift
  if [ -n "$GITLAB_TOKEN" ]; then
    H_AUTH=(-H "Authorization: Bearer $GITLAB_TOKEN")
  fi

  curl -s \
    -H "content-type: application/json" \
    "${H_AUTH[@]}" \
    "https://code.europa.eu/api/v4/$path" "$@"
}

function apiout() {
  if [ "$JSON_OUT" == "1" ]; then
    jq "."
  else
    jq "$@" | column --table
  fi
}

function git_get_project_id() {
  git remote get-url "$GIT_BRANCH_REMOTE" | sed 's|https://code.europa.eu/||;s|.git||;s|/$||' | jq -Rr @uri
}

function pipelines() {
  prid="$(git_get_project_id)"
  c_giteu "projects/$prid/pipelines" -G "$@" | apiout -r '.[] | "\(.project_id) \(.iid) \(.status) \(.ref) \(.created_at) \(.updated_at) \(.web_url)"'
}

function merge_requests() {
  prid="$(git_get_project_id)"
  c_giteu "projects/$prid/merge_requests" -G "$@" | apiout -r '.[] | "\(.project_id) \(.iid) \(.state) \(.detailed_merge_status) source:\(.source_branch) target:\(.target_branch) \(.created_at) \(.updated_at) \(.web_url)"'
}

function packages() {
  prid="$(git_get_project_id)"
  c_giteu "projects/$prid/packages" -G "$@" | apiout -r '.[] | "\(.id) \(.package_type) \(.name) \(.version) https://code.europa.eu\(._links.web_path)"'
}

main "$@"
