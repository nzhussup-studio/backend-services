#!/usr/bin/env bash

generate_service_openapi() {
  local target_dir="$1"
  local service_dir
  local service_name

  service_dir="$(cd "$target_dir" && pwd)"
  service_name="$(basename "$service_dir")"

  validate_service_dir "$service_dir"

  if [[ -f "${service_dir}/pom.xml" ]]; then
    generate_java_openapi "$service_dir" "$service_name"
  elif [[ -f "${service_dir}/go.mod" ]]; then
    generate_go_openapi "$service_dir" "$service_name"
  else
    fail "Unsupported service type. Expected pom.xml or go.mod in ${service_dir}."
  fi
}

generate_all_openapis() {
  local manifest_file
  local service

  manifest_file="$(require_manifest)"

  while IFS= read -r service; do
    [[ -n "$service" ]] || continue
    generate_service_openapi "${REPO_ROOT}/${service}"
  done < <(read_manifest_services "$manifest_file")

  build_unified_openapi "$manifest_file"
}
