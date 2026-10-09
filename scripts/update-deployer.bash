#!/bin/bash
set -e

function main() {
  local branch="main"
  local deployer="${1?Missing deployer arg}"
  local deployer_json="$(kubectl -n argocd get application -o json "$deployer")"
  local repourl="$(jq -r .spec.source.repoURL <<<"$deployer_json")"
  log "Repo url: $repourl"
  # Get projectId base chart
  projectId="$(sed 's|.*v4/projects/\([0-9]*\)/packages.*|\1|' <<<"$repourl")"
  log "ProjectID: $projectId"
  local value_file="$(c_giteu "projects/$projectId/repository/files/$(jq -nr '"charts/values.yaml" | @uri')/raw" -G -d "ref=$branch")"
  echo "$value_file" | yq -I0 -oj 'to_entries[] |
    .key as $name |
    select(.value.projectID != null and .value.targetRevision != null) |
    .value.projectID as $id |
    .value.targetRevision as $rev |
    { "projectId":$id, "property":"\($name).targetRevision", "value": $rev }
    ' | (
    while read -r line; do
      (
        id="$(jq -r .projectId <<<"$line")"
        newv="$(get_release_version "$id")"
        echo "$line" | jq -c --arg newv "$newv" '.value = $newv'
      ) &
    done
    wait
  )
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

function log() {
  >&2 echo "$@"
}

# TODO find the release given the ref
function get_release_version() {
  local prid="${1?get_release_version: Missing id}"
  c_giteu "projects/$prid/packages" | jq -r ".[0].version"
}

main "$@"
