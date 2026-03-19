#!/usr/bin/env bash

set -euo pipefail

if [[ "${GITHUB_EVENT_NAME}" == "workflow_dispatch" ]]; then
  TAG="${INPUT_TAG:-}"
  if [[ -z "$TAG" ]]; then
    echo "'tag' input is required for manual runs" >&2
    exit 1
  fi
  TAG_STRIPPED="${TAG#v}"
else
  if [[ "${GITHUB_REF}" != refs/tags/* ]]; then
    echo "Not a tag push; changelog extraction requires a tag ref." >&2
    exit 1
  fi
  TAG="${GITHUB_REF#refs/tags/}"
  TAG_STRIPPED="${TAG#v}"
fi

ENTRY="$(awk -v tag="$TAG_STRIPPED" '
  BEGIN { in_section = 0 }
  $0 ~ "^#+ \\[" tag "\\]" { in_section = 1; next }
  in_section && $0 ~ /^#+ \[[0-9]+\.[0-9]+\.[0-9]+/ { exit }
  in_section { print }
' CHANGELOG.md)"

{
  echo "changelog<<EOF"
  echo "${ENTRY}"
  echo "EOF"
} >> "$GITHUB_OUTPUT"
