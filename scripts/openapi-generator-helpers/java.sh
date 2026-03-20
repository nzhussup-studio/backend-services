#!/usr/bin/env bash

cleanup_pid() {
  local pid="${1:-}"

  if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
  fi
}

pick_free_port() {
  local port

  while true; do
    port=$((20000 + RANDOM % 20000))
    if ! lsof -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1; then
      printf '%s\n' "$port"
      return 0
    fi
  done
}

detect_java_docs_path() {
  local service_dir="$1"
  local props_file="${service_dir}/src/main/resources/application.properties"
  local docs_path=""
  local key
  local value

  if [[ -f "$props_file" ]]; then
    while IFS='=' read -r key value; do
      key="${key//[[:space:]]/}"
      if [[ "$key" == "springdoc.api-docs.path" ]]; then
        docs_path="${value//[[:space:]]/}"
        break
      fi
    done < "$props_file"
  fi

  if [[ -z "${docs_path:-}" ]]; then
    docs_path="/v3/api-docs"
  fi

  printf '%s.yaml\n' "$docs_path"
}

generate_java_openapi() {
  local service_dir="$1"
  local service_name="$2"
  local output_file="${service_dir}/openapi.yaml"
  local docs_path
  local port
  local log_file
  local maven_repo_dir="${service_dir}/.openapi-cache/m2"
  local app_pid=""

  docs_path="$(detect_java_docs_path "$service_dir")"
  port="$(pick_free_port)"
  log_file="$(mktemp -t "${service_name}.openapi")"
  mkdir -p "$maven_repo_dir"

  trap 'cleanup_pid "$app_pid"' EXIT

  (
    cd "$service_dir"
    SPRING_PROFILES_ACTIVE=openapi \
    SERVER_PORT="$port" \
    SPRING_MAIN_BANNER_MODE=off \
    SPRING_DATASOURCE_HOST=localhost \
    SPRING_DATASOURCE_PORT=5432 \
    SPRING_DATASOURCE_DB=openapi \
    SPRING_DATASOURCE_USERNAME=openapi \
    SPRING_DATASOURCE_PASSWORD=openapi \
    SPRING_REDIS_HOST=localhost \
    SPRING_REDIS_PORT=6379 \
    KEYCLOAK_JWK_SET_URI=http://127.0.0.1:8081/realms/backend-auth-dev/protocol/openid-connect/certs \
    KEYCLOAK_BACKEND_CLIENT_ID=backend-auth-client \
    ./mvnw -q -Dmaven.repo.local="$maven_repo_dir" -DskipTests spring-boot:run >"$log_file" 2>&1
  ) &
  app_pid=$!

  local url="http://127.0.0.1:${port}${docs_path}"
  local attempts=0
  local max_attempts=120

  until curl -fsS "$url" -o "$output_file" >/dev/null 2>&1; do
    if ! kill -0 "$app_pid" 2>/dev/null; then
      cat "$log_file" >&2
      fail "Java service '${service_name}' exited before OpenAPI was available."
    fi

    attempts=$((attempts + 1))
    if (( attempts >= max_attempts )); then
      cat "$log_file" >&2
      fail "Timed out waiting for ${url}."
    fi

    sleep 2
  done

  cleanup_pid "$app_pid"
  trap - EXIT
  rm -f "$log_file"
  sanitize_yaml_spec "$output_file"

  echo "Generated ${output_file}"
}
