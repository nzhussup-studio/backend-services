#!/usr/bin/env bash

fail() {
  echo "Error: $*" >&2
  exit 1
}

sanitize_yaml_spec() {
  ruby "${HELPERS_DIR}/sanitize_yaml_spec.rb" "$1"
}

require_manifest() {
  local manifest_file="${REPO_ROOT}/openapi-services.json"
  [[ -f "$manifest_file" ]] || fail "Manifest not found: ${manifest_file}"
  printf '%s\n' "$manifest_file"
}

read_manifest_services() {
  ruby "${HELPERS_DIR}/read_manifest_services.rb" "$1"
}

build_unified_openapi() {
  local output_file="${REPO_ROOT}/openapi.yaml"
  local manifest_file="$1"

  ruby "${HELPERS_DIR}/build_unified_openapi.rb" "$manifest_file" "$REPO_ROOT"
  sanitize_yaml_spec "$output_file"
  echo "Generated ${output_file}"
}

resolve_path() {
  local input="$1"

  if [[ "$input" = /* ]]; then
    printf '%s\n' "$input"
  else
    printf '%s\n' "$(cd "$PWD" && cd "$(dirname "$input")" && pwd)/$(basename "$input")"
  fi
}

validate_service_dir() {
  local service_dir="$1"
  local service_name

  service_name="$(basename "$service_dir")"

  case "$service_name" in
    api-gateway)
      fail "OpenAPI generation is not supported for ${service_name}."
      ;;
  esac

  case "$service_dir" in
    "${REPO_ROOT}"/*) ;;
    *)
      fail "Target must be inside ${REPO_ROOT}."
      ;;
  esac
}
