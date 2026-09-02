main() {
file="$1"
if [ ! -f "$file" ]; then
  >&2 echo "not found $file"
  exit 1
fi
yqc="$(getyqpr "$file"  | sed 's/^/./;s/$/ |/')"
[ -n "$yqc" ] && yq -i "$yqc .=." $file
}

getyqpr() {
code="$(cat << 'EOF'
 line="$*"
 name="$(cut -d\  -f 1 <<<"$line")"
 prId="$(cut -d\  -f 2 <<<"$line")"
 v="$(cut -d\  -f 3 <<<"$line")"
 v1="$(cut -d.  -f 1 <<<"$v")"
 v2="$(cut -d.  -f 2 <<<"$v")"
 ver="$(curl -s https://code.europa.eu/api/v4/projects/$prId/repository/tags?search=^v$v1.$v2 | jq -r .[0].name)";
 echo "$name.targetRevision = \"${ver#v}\""
EOF
)"
yq '
to_entries[]   |
select((.value | kind) == "map" and .value.projectID  != null and .value.targetRevision != null) |
"\(.key) \(.value.projectID) \(.value.targetRevision)"
' "$file" |
xargs -P10 -L1 bash -c "$code" -- | sort
}

main "$@"
