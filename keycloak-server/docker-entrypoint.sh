#!/usr/bin/env bash

set -euo pipefail

TEMPLATE_PATH="${TEMPLATE_PATH:-/opt/keycloak/conf/backend-auth-realm.template.json}"
OUTPUT_PATH="${OUTPUT_PATH:-/opt/keycloak/data/import/backend-auth-realm.json}"

required_vars=(
  BACKEND_AUTH_CLIENT_SECRET
)

for var_name in "${required_vars[@]}"; do
  if [[ -z "${!var_name:-}" ]]; then
    echo "Missing required environment variable: ${var_name}" >&2
    exit 1
  fi
done

escape_sed_value() {
  printf '%s' "$1" | sed -e 's/[\/&]/\\&/g'
}

render_string_value() {
  local placeholder="$1"
  local value="$2"
  sed -i "s/${placeholder}/$(escape_sed_value "${value}")/g" "${OUTPUT_PATH}"
}

render_json_value() {
  local placeholder="$1"
  local value="$2"
  sed -i "s/\"${placeholder}\"/$(escape_sed_value "${value}")/g" "${OUTPUT_PATH}"
}

resolve_boolean() {
  local value="${1,,}"

  case "${value}" in
    true|false)
      printf '%s' "${value}"
      ;;
    *)
      echo "Invalid boolean value: ${1}" >&2
      exit 1
      ;;
  esac
}

cp "${TEMPLATE_PATH}" "${OUTPUT_PATH}"

realm_id="${KEYCLOAK_REALM_ID:-backend-auth}"
realm_name="${KEYCLOAK_REALM_NAME:-${realm_id}}"
realm_display_name="${KEYCLOAK_REALM_DISPLAY_NAME:-Backend Auth}"
ssl_required="${KEYCLOAK_SSL_REQUIRED:-external}"

registration_allowed="$(resolve_boolean "${KEYCLOAK_REGISTRATION_ALLOWED:-true}")"
login_with_email_allowed="$(resolve_boolean "${KEYCLOAK_LOGIN_WITH_EMAIL_ALLOWED:-true}")"
reset_password_allowed="$(resolve_boolean "${KEYCLOAK_RESET_PASSWORD_ALLOWED:-true}")"
remember_me="$(resolve_boolean "${KEYCLOAK_REMEMBER_ME:-true}")"
verify_email="$(resolve_boolean "${KEYCLOAK_VERIFY_EMAIL:-true}")"

backend_client_id="${KEYCLOAK_BACKEND_CLIENT_ID:-backend-auth-client}"
backend_client_name="${KEYCLOAK_BACKEND_CLIENT_NAME:-Backend Auth Client}"
backend_standard_flow_enabled="$(resolve_boolean "${KEYCLOAK_BACKEND_STANDARD_FLOW_ENABLED:-false}")"
backend_direct_access_grants_enabled="$(resolve_boolean "${KEYCLOAK_BACKEND_DIRECT_ACCESS_GRANTS_ENABLED:-false}")"

frontend_client_id="${KEYCLOAK_FRONTEND_CLIENT_ID:-frontend-admin-auth-client}"
frontend_client_name="${KEYCLOAK_FRONTEND_CLIENT_NAME:-Frontend Admin Auth Client}"
frontend_root_url="${KEYCLOAK_FRONTEND_ROOT_URL:-https://admin.example.com}"
frontend_base_url="${KEYCLOAK_FRONTEND_BASE_URL:-${frontend_root_url}}"
frontend_redirect_uris_json="${KEYCLOAK_FRONTEND_REDIRECT_URIS_JSON:-[\"https://admin.example.com/*\"]}"
frontend_web_origins_json="${KEYCLOAK_FRONTEND_WEB_ORIGINS_JSON:-[\"https://admin.example.com\"]}"

users_json="[]"
if [[ -n "${KEYCLOAK_USERS_FILE:-}" ]]; then
  if [[ ! -f "${KEYCLOAK_USERS_FILE}" ]]; then
    echo "Users file not found: ${KEYCLOAK_USERS_FILE}" >&2
    exit 1
  fi
  users_json="$(tr -d '\n' < "${KEYCLOAK_USERS_FILE}")"
fi

google_enabled="${KEYCLOAK_GOOGLE_ENABLED:-}"
if [[ -z "${google_enabled}" ]]; then
  if [[ -n "${GOOGLE_CLIENT_ID:-}" && -n "${GOOGLE_CLIENT_SECRET:-}" ]]; then
    google_enabled="true"
  else
    google_enabled="false"
  fi
fi
google_enabled="$(resolve_boolean "${google_enabled}")"

github_enabled="${KEYCLOAK_GITHUB_ENABLED:-}"
if [[ -z "${github_enabled}" ]]; then
  if [[ -n "${GITHUB_CLIENT_ID:-}" && -n "${GITHUB_CLIENT_SECRET:-}" ]]; then
    github_enabled="true"
  else
    github_enabled="false"
  fi
fi
github_enabled="$(resolve_boolean "${github_enabled}")"

render_string_value "__REALM_ID__" "${realm_id}"
render_string_value "__REALM_NAME__" "${realm_name}"
render_string_value "__REALM_DISPLAY_NAME__" "${realm_display_name}"
render_string_value "__SSL_REQUIRED__" "${ssl_required}"
render_json_value "__REGISTRATION_ALLOWED__" "${registration_allowed}"
render_json_value "__LOGIN_WITH_EMAIL_ALLOWED__" "${login_with_email_allowed}"
render_json_value "__RESET_PASSWORD_ALLOWED__" "${reset_password_allowed}"
render_json_value "__REMEMBER_ME__" "${remember_me}"
render_json_value "__VERIFY_EMAIL__" "${verify_email}"

render_string_value "__BACKEND_CLIENT_ID__" "${backend_client_id}"
render_string_value "__BACKEND_CLIENT_NAME__" "${backend_client_name}"
render_string_value "__BACKEND_AUTH_CLIENT_SECRET__" "${BACKEND_AUTH_CLIENT_SECRET}"
render_json_value "__BACKEND_STANDARD_FLOW_ENABLED__" "${backend_standard_flow_enabled}"
render_json_value "__BACKEND_DIRECT_ACCESS_GRANTS_ENABLED__" "${backend_direct_access_grants_enabled}"

render_string_value "__FRONTEND_CLIENT_ID__" "${frontend_client_id}"
render_string_value "__FRONTEND_CLIENT_NAME__" "${frontend_client_name}"
render_string_value "__FRONTEND_ROOT_URL__" "${frontend_root_url}"
render_string_value "__FRONTEND_BASE_URL__" "${frontend_base_url}"
render_json_value "__FRONTEND_REDIRECT_URIS_JSON__" "${frontend_redirect_uris_json}"
render_json_value "__FRONTEND_WEB_ORIGINS_JSON__" "${frontend_web_origins_json}"
render_json_value "__USERS_JSON__" "${users_json}"

render_json_value "__GOOGLE_ENABLED__" "${google_enabled}"
render_string_value "__GOOGLE_CLIENT_ID__" "${GOOGLE_CLIENT_ID:-}"
render_string_value "__GOOGLE_CLIENT_SECRET__" "${GOOGLE_CLIENT_SECRET:-}"

render_json_value "__GITHUB_ENABLED__" "${github_enabled}"
render_string_value "__GITHUB_CLIENT_ID__" "${GITHUB_CLIENT_ID:-}"
render_string_value "__GITHUB_CLIENT_SECRET__" "${GITHUB_CLIENT_SECRET:-}"

exec /opt/keycloak/bin/kc.sh "$@"
