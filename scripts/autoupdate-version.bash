#!/bin/bash
set -e
function get_projects_info() {
  yq -oj -I0 '
  to_entries[]
    | .key as $path
    | .value
    | select(.projectID != null)
    | {"path": $path, "projectID": .projectID, "targetRevision": .targetRevision}
  '
}

function get_new_version() {
  read -r prj
  local projectID="$(<<<"$prj" jq -r .projectID)"
  local basetag="v$(<<<"$prj" jq -r .targetRevision | cut -d. -f1,2)"
  local latestTargetRevision="$(c_giteu "projects/$projectID/repository/tags?search=^$basetag&page_size=1" | jq -r '.[0]?.name' | sed 's/^v//')"
  jq -c -n --argjson prj "$prj"  --arg latestTargetRevision "$latestTargetRevision" '$prj + {"latestTargetRevision": $latestTargetRevision}'
}

function fetch_project_tags_path() {
  read -r prj
  path="$(<<<"$prj" jq -r .path)"
  prid="$(<<<"$prj" jq -r .projectID)"
  from="v$(<<<"$prj" jq -r .targetRevision)"
  to="v$(<<<"$prj" jq -r .latestTargetRevision)"
  repo="$(giteu api "projects/$prid" | jq -r .http_url_to_repo)"
  git fetch -n "$repo" +"refs/tags/$from:refs/from" +"refs/tags/$to:refs/to"
  msg="$path ($from -> $to):
$(git log --grep Changelog: --no-merges --oneline refs/from..refs/to | sed 's/^/  /')"
  <<<"$prj" jq -c --arg msg "$msg" '.+{"message":$msg}'
}

function collect_updates() {
  local msg=""
  source_content_yaml="${1?Missing source content}"
  new_version="${2?Missing new version}"
  source_content_yaml="$(echo "$source_content_yaml" | NEW_VER="$new_version" yq '.values.branch = strenv(NEW_VER)')"
  while read -r prj; do
    targetRevision="$(<<<"$prj" jq -r .targetRevision)"
    latestTargetRevision="$(<<<"$prj" jq -r .latestTargetRevision)"
    if [ "$targetRevision" == "$latestTargetRevision" ]; then
      continue
    fi
    path="$(<<<"$prj" jq -r .path)"
    message="$(<<<"$prj" jq -r .message)"
    source_content_yaml="$(
    echo "$source_content_yaml" |
    latestTargetRevision="$latestTargetRevision" \
    yq '('."$path"'.targetRevision) = strenv(latestTargetRevision)'
    )"
    if [ -n "$message" ]; then
      msg="$message
$msg"
    fi
  done
  if [ -n "$msg" ]; then
    msg="Bump versions

$msg

Changelog: fixed"
  fi
  source="$source_content_yaml" \
  msg="$msg" \
  yq -I0 -n -oj '{"message": strenv(msg), "source":strenv(source)}'
}

function get_version_from_pipeline_variables() {
 grep -w "PROJECT_VERSION_NUMBER" | grep -o '".*"' | sed 's/"//g'
}

function update_release() {
  prid="${1?Missing project id}"
  remote_branch="${2?Missing remote branch}"
  repo="$(giteu api "projects/$prid" | jq -r .http_url_to_repo)"
  git fetch --depth=1 -n "$repo" +"refs/$remote_branch":refs/source
  source="$(git show refs/source:charts/values.yaml)"
  pipeline_var_content="$(git show refs/source:pipeline.variables.sh)"
  ver_chart="$(<<<"$pipeline_var_content" get_version_from_pipeline_variables)"
  prefix="$(<<<"$ver_chart" cut -d. -f 1,2)"
  final="$(<<<"$ver_chart" cut -d. -f 3)"
  final="$(($final + 1))"
  ver_chart_new="$prefix.$final"
  out="$(
    echo "$source" | get_projects_info | while read -r line; do
    echo $line | get_new_version | fetch_project_tags_path
    done | collect_updates "$source" "v$ver_chart_new"
  )"
  pipeline_var_content="$(<<<"$pipeline_var_content" sed '/PROJECT_VERSION_NUMBER/s/".*"/"'"$ver_chart_new"'"/')"

  pipeline_sha="$(echo "$pipeline_var_content" | git hash-object -w --stdin)"
  chart_sha="$(echo "$out" | jq -r .source | git hash-object -w --stdin)"
  msg="$(echo "$out" | jq -r .message)"
  if [ -z "$msg" ]; then
    log "No message found. Exit"
    exit
  fi


  export GIT_INDEX_FILE="$(mktemp -u)"
  git read-tree refs/source
  git update-index --add --cacheinfo 100644,"$pipeline_sha",pipeline.variables.sh 
  git update-index --add --cacheinfo 100644,"$chart_sha",charts/values.yaml 
  tree="$(git write-tree)"
  commitID="$(git commit-tree "$tree" -m "$msg" -p "refs/source")"
  git update-ref refs/updated "$commitID"
}

function c_giteu() {
  path="${1?Missing path}"
  shift
  # if [ -n "$GITLAB_TOKEN" ]; then
  #   H_AUTH=(-H "Authorization: Bearer $GITLAB_TOKEN")
  # fi

  curl -s \
    -H "content-type: application/json" \
    "${H_AUTH[@]}" \
    "https://code.europa.eu/api/v4/$path" "$@"
}

function log() {
  >&2 echo "$@"
}

update_release "$@"
