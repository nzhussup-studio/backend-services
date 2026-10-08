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
  local gateway_spec="${REPO_ROOT}/nginx-gateway/docs/specs/openapi.yaml"

  ruby "${HELPERS_DIR}/build_unified_openapi.rb" "$manifest_file" "$REPO_ROOT"
  sanitize_yaml_spec "$output_file"
  echo "Generated ${output_file}"

  mkdir -p "$(dirname "$gateway_spec")"
  cp "$output_file" "$gateway_spec"
  echo "Synced ${gateway_spec}"
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
    nginx-gateway)
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

configure_java_runtime() {
  local service_dir="$1"
  local expected_version
  local current_version
  local java_home

  expected_version="$(sed -n 's:.*<java.version>[[:space:]]*\([0-9][0-9]*\)[[:space:]]*</java.version>.*:\1:p' "${service_dir}/pom.xml" | head -n 1)"
  [[ -n "$expected_version" ]] || return 0

  current_version="$(java -version 2>&1 | sed -n 's/.*version "\([0-9][0-9]*\).*/\1/p' | head -n 1)"
  if [[ "$current_version" == "$expected_version" ]]; then
    return 0
  fi

  if [[ "$(uname -s)" == "Darwin" ]] && command -v /usr/libexec/java_home >/dev/null 2>&1; then
    java_home="$(/usr/libexec/java_home -v "$expected_version" 2>/dev/null || true)"
    if [[ -n "$java_home" ]]; then
      current_version="$("${java_home}/bin/java" -version 2>&1 | sed -n 's/.*version "\([0-9][0-9]*\).*/\1/p' | head -n 1)"
      if [[ "$current_version" == "$expected_version" ]]; then
        export JAVA_HOME="$java_home"
        export PATH="${JAVA_HOME}/bin:${PATH}"
        return 0
      fi
    fi
  fi

  fail "Java service '${service_dir##*/}' requires JDK ${expected_version}; found JDK ${current_version:-unknown}. Install JDK ${expected_version} or set JAVA_HOME before generating OpenAPI."
}
