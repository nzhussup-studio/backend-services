#!/usr/bin/env bash

go_swag_version() {
  local service_dir="$1"
  local version=""
  local line

  while IFS= read -r line; do
    case "$line" in
      *github.com/swaggo/swag\ *)
        line="${line#"${line%%[![:space:]]*}"}"
        line="${line#github.com/swaggo/swag }"
        version="${line%% *}"
        break
        ;;
    esac
  done < "${service_dir}/go.mod"

  if [[ -z "${version:-}" ]]; then
    version="v1.16.4"
  fi

  printf '%s\n' "$version"
}

generate_go_openapi() {
  local service_dir="$1"
  local service_name="$2"
  local output_file="${service_dir}/openapi.yaml"
  local docs_dir="${service_dir}/docs"
  local go_mod_cache="${service_dir}/.openapi-cache/gomod"
  local go_build_cache="${service_dir}/.openapi-cache/go-build"
  local main_file
  local main_dir_rel
  local main_name
  local swag_version

  mkdir -p "$go_mod_cache" "$go_build_cache"
  main_file="$(find "${service_dir}/cmd" -name main.go -print | head -n 1)"
  [[ -n "$main_file" ]] || fail "Could not locate a Go entrypoint under ${service_dir}/cmd."

  main_dir_rel="$(dirname "${main_file#${service_dir}/}")"
  main_name="$(basename "$main_file")"
  swag_version="$(go_swag_version "$service_dir")"

  (
    cd "$service_dir"
    export GOMODCACHE="$go_mod_cache"
    export GOCACHE="$go_build_cache"
    if command -v swag >/dev/null 2>&1; then
      swag init \
        --generalInfo "$main_name" \
        --dir "${main_dir_rel},internal" \
        --output docs \
        --outputTypes go,json,yaml
    else
      go run "github.com/swaggo/swag/cmd/swag@${swag_version}" init \
        --generalInfo "$main_name" \
        --dir "${main_dir_rel},internal" \
        --output docs \
        --outputTypes go,json,yaml
    fi
  )

  [[ -f "${docs_dir}/swagger.yaml" ]] || fail "Go OpenAPI generation did not produce ${docs_dir}/swagger.yaml."
  cp "${docs_dir}/swagger.yaml" "$output_file"
  sanitize_yaml_spec "$output_file"

  echo "Generated ${output_file}"
}
