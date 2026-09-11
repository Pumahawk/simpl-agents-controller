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
  help | -h | --help)
    shift
    help "$@"
    ;;
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
  rel:links)
    shift
    release_links "$@"
    ;;
  *)
    c_giteu "$@"
    ;;
  esac
}

function help() {
  cat <<'EOF'
giteu - Command-line helper for Code.europa.eu (GitLab)

USAGE
  giteu [GLOBAL OPTIONS] <command> [COMMAND OPTIONS]

GLOBAL OPTIONS

  --json
      Return the raw JSON response from the GitLab API.

  --remote <remote>
      Git remote to use when resolving the current project.
      Default: origin

COMMANDS

  info
      Display information about the current project.

      Example:
        giteu info

  mr | mergerequests
      List merge requests for the current project.

      Examples:
        giteu mr
        giteu mr -d state=opened

  mr:create
      Create a new merge request.

      Options:
        --from <branch>   Source branch.
                          Default: current branch.
        --to <branch>     Target branch.
                          Default: develop.
        --title <text>    Merge request title.
                          Default: latest commit subject.

      Examples:
        giteu mr:create
        giteu mr:create --to main
        giteu mr:create --title "Release 1.2.0"

  pips | pipelines
      List pipelines for the current project.

      Examples:
        giteu pips
        giteu pips -d ref=develop

  pips:run
      Trigger a new pipeline.

      Options:
        --ref <branch|tag>
                          Branch or tag to run.
                          Default: current branch.
        --release         Set RUN_RELEASE=true.

      Examples:
        giteu pips:run
        giteu pips:run --ref develop
        giteu pips:run --release

  pk | packages
      List packages published for the current project.

      Example:
        giteu pk

  rel:links [tag]
      Display useful links associated with a release:
        - GitLab Release
        - Fortify Report
        - SonarQube Report

      If a tag is provided, the corresponding release is used.
      Otherwise, the most recent release is selected.

      Examples:
        giteu rel:links
        giteu rel:links v1.2.3

  <api-path>
      Execute a direct GitLab API request.

      Examples:
        giteu projects
        giteu projects/123
        giteu projects/123/repository/tags

ENVIRONMENT VARIABLES

  GITLAB_TOKEN
      Personal access token used to authenticate against
      https://code.europa.eu/api/v4

EXAMPLES

  giteu info
  giteu mr
  giteu mr:create --to main
  giteu pips
  giteu pips:run --release
  giteu rel:links

EOF
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
  c_giteu "projects/$prid/pipelines" -G "$@" | apiout -r '.[] | "\(.project_id) \(.id) \(.status) \(.ref) \(.created_at) \(.updated_at) \(.web_url)"'
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
      shift
      shift
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

function release_links() {
  curlp_tag=()
  tag_p="$1"
  [ -n "$tag_p" ] && tag_p+=(-d "search=$tag_p")
  info="$(giteu --json info)"
  prid="$(echo "$info" | jq -r .id)"
  name="$(echo "$info" | jq -r .path)"
  tag="$(giteu "projects/$prid/repository/tags" -G "${tag_p[@]}" -d per_page=1 | jq -r .[].name)"
  releasej="$(giteu "projects/$prid/releases/$tag")"
  description="$(echo "$releasej" | jq -r .description)"
  releaseurl="$(echo "$releasej" | jq -r ._links.self)"
  fortify="$(echo "$description" | grep -e emea.fortify.com | grep -o "(.*)")"
  sonar="$(echo "$description" | grep -e sonarqube.tools.simpl-europe.eu | grep -o "(.*)")"
  [ -z "$fortify" ] && fortify='"not-found"'
  [ -z "$sonar" ] && sonar='"not-found"'
  echo "$prid $name $tag type:release $releaseurl"
  echo "$prid $name $tag type:fortify $fortify"
  echo "$prid $name $tag type:sonar $sonar"
}

main "$@"
