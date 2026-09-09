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
  info)
    shift
    project_info "$@"
    ;;
  mr:create)
    shift
    merge_requests_create "$@"
    ;;
  pips:run)
    shift
    pipeline_run "$@"
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
  c_giteu "projects/$prid/packages" -G "$@" | apiout -r '.[] | "\(.id) \(.package_type) \(.name) \(.version) \(.pipeline.ref) https://code.europa.eu\(._links.web_path)"'
}

function merge_requests_create() {
  local source_branch="$(git branch --show-current)"
  local target_branch="develop"
  local title="$(git log --no-merges -1 --pretty=%s)"
  while true; do
    case "$1" in
    --title)
      title="$2"
      shift
      shift
      ;;
    --from)
      source_branch="$2"
      shift
      shift
      ;;
    --to)
      target_branch="$2"
      shift
      shift
      ;;
    *)
      break
      ;;
    esac
  done
  if [ -z "$source_branch" ]; then
    >&2 echo "Missing --from parameter"
    return 1
  fi
  if [ -z "$target_branch" ]; then
    >&2 echo "Missing --to parameter"
    return 1
  fi
  if [ -z "$title" ]; then
    >&2 echo "Missing --title parameter"
    return 1
  fi
  prid="$(git_get_project_id)"
  c_giteu /projects/"$prid"/merge_requests -X POST --data-raw "$(
    jq -n -c \
      --arg source "$source_branch" \
      --arg target "$target_branch" \
      --arg title "$title" \
      '{
         "source_branch":$source,
         "target_branch":$target,
         "title":$title,
       }'
  )" "$@" | apiout -r '"\(.project_id) \(.id) \(.web_url)"'
}

function project_info() {
  prid="$(git_get_project_id)"
  c_giteu "projects/$prid" -G "$@" | apiout -r '"\(.id) \(.path) \(.web_url)"'
}

function pipeline_run() {
  local ref="$(git branch --show-current)"
  local release="0"
  while true; do
    case "$1" in
    --ref)
      ref="$2"
      shift 2
      ;;
    --release)
      release="1"
      shift
      ;;
    *)
      break
      ;;
    esac
  done
  if [ -z "$ref" ]; then
    >&2 echo "Missing --ref parameter"
    return 1
  fi
  local variables="[]"
  if [ "$release" == "1" ]; then
    variables='[{"key":"RUN_RELEASE","value":"true"}]'
  fi
  prid="$(git_get_project_id)"
  c_giteu "projects/$prid/pipeline?ref=$ref" -X POST --data-raw "$(jq --argjson variables "$variables" -n '{"variables": $variables}')" | apiout -r '"\(.project_id) \(.iid) \(.web_url)"'
}

main "$@"
