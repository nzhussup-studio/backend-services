#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
HELPERS_DIR="${SCRIPT_DIR}/openapi-generator-helpers"

# shellcheck source=/dev/null
source "${HELPERS_DIR}/common.sh"
# shellcheck source=/dev/null
source "${HELPERS_DIR}/java.sh"
# shellcheck source=/dev/null
source "${HELPERS_DIR}/go.sh"
# shellcheck source=/dev/null
source "${HELPERS_DIR}/workflows.sh"

usage() {
  echo "Usage: $(basename "$0") <path-to-folder>"
  echo "       $(basename "$0") -a|--all"
}

main() {
  [[ $# -ge 1 && $# -le 2 ]] || {
    usage
    exit 1
  }

  case "${1}" in
    -a|--all)
      [[ $# -eq 1 ]] || fail "The ${1} flag does not accept additional arguments."
      generate_all_openapis
      ;;
    *)
      [[ $# -eq 1 ]] || fail "Single-service mode expects exactly one path argument."
      local target_dir
      target_dir="$(resolve_path "$1")"
      [[ -d "$target_dir" ]] || fail "Directory not found: $target_dir"
      generate_service_openapi "$target_dir"
      ;;
  esac
}

main "$@"
